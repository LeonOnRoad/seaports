package server

import (
	"fmt"
	"net/http"

	"company.com/seaports/services/ports-api/controller"
	"company.com/seaports/services/ports-api/controller/writer"
	"company.com/seaports/services/ports-api/service"
	"github.com/gorilla/mux"
)

type Resources struct {
	PortService service.PortInterface
}

func StartAsync(port int, res *Resources) *http.Server {
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: configure(res)}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return server
}

func configure(res *Resources) *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = &notFoundHandler{}

	rootController := controller.NewRoot()
	router.HandleFunc("/", rootController.Get)

	portsController := controller.NewPort(res.PortService)
	router.HandleFunc("/api/ports/{id}", portsController.Get).Methods(http.MethodGet)
	router.HandleFunc("/api/ports:import", portsController.Import).Methods(http.MethodPost)

	return router
}

type notFoundHandler struct{}

// called on an unexisting route
func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writer.Write(w, http.StatusNotFound, "Page not found")
}
