package service

import "company.com/seaports/services/ports-api/model"

type PortInterface interface {
	Get(id string) (*model.Port, error)
	ImportPorts(portsChan <-chan *model.Port) (*model.ImportPortsResponse, error)
}
