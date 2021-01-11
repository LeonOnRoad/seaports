package model

import (
	pbPort "company.com/seaports/proto/src/api/port"
)

type Port pbPort.Port // map modeling model to proto

type ImportPortsResponse pbPort.ImportPortsResponse

func ConvertModelPortToProtoPort(mp *Port) *pbPort.Port {
	return (*pbPort.Port)(mp)
}

func ConvertProtoResponseToModelResponse(pr *pbPort.ImportPortsResponse) *ImportPortsResponse {
	return (*ImportPortsResponse)(pr)
}
