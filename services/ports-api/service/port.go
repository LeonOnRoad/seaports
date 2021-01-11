package service

import (
	"context"
	"log"

	pbPort "company.com/seaports/proto/src/api/port"
	"company.com/seaports/services/ports-api/model"
)

type Port struct {
	grpcPortsServiceClient pbPort.PortsServiceClient
}

func NewPort(pc pbPort.PortsServiceClient) *Port {
	return &Port{
		grpcPortsServiceClient: pc,
	}
}

func (s Port) Get(id string) (*model.Port, error) {
	port, err := s.grpcPortsServiceClient.GetPort(context.Background(), &pbPort.GetPortRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return (*model.Port)(port), nil
}

func (s Port) ImportPorts(portsChan <-chan *model.Port) (*model.ImportPortsResponse, error) {
	stream, err := s.grpcPortsServiceClient.StreamImportedPorts(context.Background())
	if err != nil {
		return nil, err
	}

	for port := range portsChan {
		protoPort := model.ConvertModelPortToProtoPort(port)
		err = stream.Send(protoPort)
		if err != nil {
			log.Printf("Failed to send port through stream. Error: %s", err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	return model.ConvertProtoResponseToModelResponse(resp), nil
}
