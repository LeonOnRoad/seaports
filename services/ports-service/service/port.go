package service

import (
	"context"

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
	return nil, nil
}

func (s Port) StreamImportedPorts(stream pbPort.PortsService_StreamImportedPortsServer) error {
	return nil
}
