package service

import (
	"company.com/seaports/services/ports-api/model"
)

type Port struct {
}

func NewPort() *Port {
	return &Port{}
}

func (s Port) Get(id string) (*model.Port, error) {
	return nil, nil
}

func (s Port) ImportPorts(portsChan <-chan *model.Port) (*model.ImportPortsResponse, error) {
	return nil, nil
}
