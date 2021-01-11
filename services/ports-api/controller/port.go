package controller

import (
	"net/http"

	"github.com/gorilla/mux"

	"company.com/seaports/services/ports-api/controller/writer"
	"company.com/seaports/services/ports-api/service"
)

type Port struct {
	portService service.PortInterface
}

func NewPort(portService service.PortInterface) *Port {
	return &Port{
		portService: portService,
	}
}

func (c *Port) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	port, err := c.portService.Get(id)
	if err != nil {
		writer.WriteError(w, writer.ConvertGrpcError("port", err))
		return
	}
	writer.Write(w, http.StatusOK, port)
}

func (c *Port) Import(w http.ResponseWriter, r *http.Request) {
	writer.Write(w, http.StatusCreated, "POST /api/ports:import called")
}
