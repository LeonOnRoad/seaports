package controller

import (
	"net/http"

	"company.com/seaports/services/ports-api/controller/writer"
)

type Root struct{}

func NewRoot() *Root {
	return &Root{}
}

func (c *Root) Get(w http.ResponseWriter, r *http.Request) {
	writer.Write(w, http.StatusOK, "Welcome!")
}
