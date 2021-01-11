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

	rootController := controller.NewRoot()
	router.HandleFunc("/", rootController.Get).Methods(http.MethodGet)

	portsController := controller.NewPorts()
	router.HandleFunc("/api/ports/{id}", portsController.Get).Methods(http.MethodGet)
	router.HandleFunc("/api/ports:import", portsController.Import).Methods(http.MethodPost)

	return router
}
