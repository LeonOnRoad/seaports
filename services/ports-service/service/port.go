package service

import (
	"context"
	"encoding/json"
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
		log.Printf("Failed to unmarshal bytes to port. Error: %s", err)
		return nil, status.Error(codes.Internal, "Failed to get port")
	}
	return port, nil
}

func (s Port) StreamImportedPorts(stream pbPort.PortsService_StreamImportedPortsServer) error {
	return nil
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
