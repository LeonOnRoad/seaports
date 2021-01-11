package controller

import (
	"net/http"

	"company.com/seaports/services/ports-api/controller/writer"
)

type Ports struct {
}

func NewPorts() *Ports {
	return &Ports{}
}

func (c *Ports) Get(w http.ResponseWriter, r *http.Request) {
	writer.Write(w, http.StatusOK, "GET /api/port called")
}

func (c *Ports) Import(w http.ResponseWriter, r *http.Request) {
	writer.Write(w, http.StatusCreated, "POST /api/ports:import called")
}
