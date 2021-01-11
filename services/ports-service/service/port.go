package service

import (
	"context"

	pbPort "company.com/seaports/proto/src/api/port"
)

var _ pbPort.PortsServiceServer = (*Port)(nil) // check if service implements generated proto interface

type Port struct {
	pbPort.UnimplementedPortsServiceServer
}

func NewPort() *Port {
	return &Port{}
}

func (s Port) GetPort(ctx context.Context, req *pbPort.GetPortRequest) (*pbPort.Port, error) {
	return nil, nil
}

func (s Port) StreamImportedPorts(stream pbPort.PortsService_StreamImportedPortsServer) error {
	return nil
}
