package server

import (
	"fmt"
	"log"
	"net/http"

	"company.com/seaports/services/ports-api/controller"
	"github.com/gorilla/mux"
)

func StartAsync(port int) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: configure(),
	}

	go func() {
		log.Println("Starting http server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return server
}

func configure() *mux.Router {
	router := mux.NewRouter()

	portsController := controller.NewPorts()
	router.HandleFunc("/api/ports/{id}", portsController.Get).Methods("GET")
	router.HandleFunc("/api/ports:import", portsController.Import).Methods("POST")

	return router
}
