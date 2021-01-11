package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbPort "company.com/seaports/proto/src/api/port"

	"company.com/seaports/services/ports-service/repository"
)

var _ pbPort.PortsServiceServer = (*Port)(nil) // check if service implements generated proto interface

type Port struct {
	pbPort.UnimplementedPortsServiceServer

	repo repository.Repository
}

func NewPort(r repository.Repository) *Port {
	return &Port{
		repo: r,
	}
}

func (s Port) GetPort(ctx context.Context, req *pbPort.GetPortRequest) (*pbPort.Port, error) {
	val, rErr := s.repo.Get(ctx, req.Id)
	if rErr != nil {
		return nil, convertRepoErrorToGrpcError(rErr)
	}

	port := &pbPort.Port{}
	err := json.Unmarshal([]byte(val), port)
	if err != nil {
		log.Printf("Failed to unmarshal bytes to port. PortID: %s Error: %s", req.Id, err)
		return nil, status.Errorf(codes.Internal, "Failed to get port %s", req.Id)
	}
	return port, nil
}

func (s Port) StreamImportedPorts(stream pbPort.PortsService_StreamImportedPortsServer) error {
	response := &pbPort.ImportPortsResponse{}
	for {
		port, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(response)
			}
			return err
		}

		s.process(stream.Context(), port, response)
	}
}

func (s Port) process(ctx context.Context, port *pbPort.Port, response *pbPort.ImportPortsResponse) {
	if port == nil {
		return
	}

	bytes, err := json.Marshal(port)
	if err != nil {
		log.Printf("Failed to marshal port[%s] to bytes. Error: %s", port.Id, err)
		response.Errors = append(response.Errors,
			fmt.Sprintf("Failed to save port %s. Error while marshaling: %s", port.Id, err))
		return
	}

	val, rErr := s.repo.Get(ctx, port.Id)
	if rErr != nil {
		if rErr.Type == repository.NOT_FOUND {
			if rErr := s.repo.Set(ctx, port.Id, string(bytes)); rErr != nil {
				log.Printf("Failed to save port %s. Repository error: %s", port.Id, rErr.Err.Error())
				response.Errors = append(response.Errors,
					fmt.Sprintf("Failed to save port %s. Repository error: %s", port.Id, rErr.Err.Error()))
			} else {
				response.Created++
			}
		}
	} else {
		if val != string(bytes) { // updated only if something changed
			if err := s.repo.Set(ctx, port.Id, string(bytes)); err != nil {
				log.Printf("Failed to save port %s. Repository error: %s", port.Id, rErr.Err.Error())
				response.Errors = append(response.Errors,
					fmt.Sprintf("Failed to save port %s. Repository error: %s", port.Id, rErr.Err.Error()))
			} else {
				response.Updated++
			}
		}
	}
	response.Total++
}

func convertRepoErrorToGrpcError(re *repository.RepoError) error {
	if re == nil {
		return nil
	}

	switch re.Type {
	case repository.NONE:
		return status.Error(codes.OK, re.Err.Error())
	case repository.NOT_FOUND:
		return status.Error(codes.NotFound, re.Err.Error())
	case repository.ALREADY_EXISTS:
		return status.Error(codes.AlreadyExists, re.Err.Error())
	default:
		return status.Error(codes.Internal, re.Err.Error())
	}
}
