package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"company.com/seaports/services/ports-api/controller/writer"
	"company.com/seaports/services/ports-api/model"
	"company.com/seaports/services/ports-api/service"
)

type Port struct {
	portService service.PortInterface
}

func NewPort(ps service.PortInterface) *Port {
	return &Port{
		portService: ps,
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
	defer r.Body.Close()

	parseChan, err := readPortsFromStream(r.Body)
	if err != nil {
		writer.WriteError(w, writer.NewResponseError(http.StatusBadRequest, "Invalid json provided", err))
		return
	}

	portsChan := make(chan *model.Port)
	parseErrors := make([]string, 0)
	go func() {
		for {
			parseResult := <-parseChan
			if parseResult == nil {
				close(portsChan)
				return
			}
			if parseResult.err != nil {
				log.Printf("Failed to parse port from json. Error: %s", err)
				parseErrors = append(parseErrors, parseResult.err.Error())
				continue
			}
			portsChan <- parseResult.port
		}
	}()

	response, err := c.portService.ImportPorts(portsChan)
	if err != nil {
		writer.WriteError(w, writer.NewResponseError(http.StatusInternalServerError, "Failed to process ports", err))
		return
	}

	response.Errors = append(response.Errors, parseErrors...) // append parse errors to response errors
	writer.Write(w, http.StatusCreated, response)
}

type parseResult struct {
	port *model.Port
	err  error
}

func readPortsFromStream(inStream io.Reader) (<-chan *parseResult, error) {
	decoder := json.NewDecoder(inStream)
	_, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	parseChan := make(chan *parseResult, 100)
	go func() {
		for decoder.More() {
			token, err := decoder.Token()
			if err != nil {
				log.Printf("Failed to decode json token. Error: %s\n", err)
				parseChan <- &parseResult{err: err}
				continue
			}

			s, ok := token.(string)
			if !ok {
				log.Printf("Json token not a string: %v\n", token)
				parseChan <- &parseResult{err: fmt.Errorf("Json token not a string: %v\n", token)}
				continue
			}

			port := &model.Port{}
			err = decoder.Decode(port)
			if err != nil {
				log.Printf("Failed to decode json object. Error: %s\n", err)
				parseChan <- &parseResult{err: fmt.Errorf("Failed to decode json object. Error: %s\n", err)}
				continue
			}
			port.Id = s
			parseChan <- &parseResult{port: port}
		}
		close(parseChan)
	}()

	return parseChan, nil
}
