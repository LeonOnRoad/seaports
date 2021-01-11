package model

import (
	pbPort "company.com/seaports/proto/src/api/port"
)

type Port pbPort.Port // map modeling model to proto

type ImportPortsResponse pbPort.ImportPortsResponse
