package httptrigger

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/edgexfoundry/edgex-go/pkg/models"

	"github.com/edgexfoundry-holdings/app-functions-sdk-go/pkg/context"

	"github.com/edgexfoundry-holdings/app-functions-sdk-go/pkg/configuration"
	"github.com/edgexfoundry-holdings/app-functions-sdk-go/pkg/runtime"
)

// HTTPTrigger implements ITrigger to support HTTPTriggers
type HTTPTrigger struct {
	Configuration configuration.Configuration
	Runtime       runtime.GolangRuntime
	outputData    interface{}
}

// Initialize ...
func (h *HTTPTrigger) Initialize() error {
	http.HandleFunc("/", h.requestHandler)   // set router - just a GET for now
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	return nil
}
func (h *HTTPTrigger) requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("BOOM EXECUTED")
	decoder := json.NewDecoder(r.Body)

	// event := event.Event{Data: "DATA FROM HTTP"}
	edgexContext := context.Context{Configuration: h.Configuration,
		Trigger: h,
	}
	var event models.Event
	decoder.Decode(&event)

	h.Runtime.ProcessEvent(edgexContext, event)
	bytes, _ := getBytes(h.outputData)
	w.Write(bytes)
}
func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

// GetConfiguration gets the config
func (h *HTTPTrigger) GetConfiguration() configuration.Configuration {
	//
	return h.Configuration
}

// GetData This function might return data
func (h *HTTPTrigger) GetData() interface{} {
	return "data"
}

// Complete ...
func (h *HTTPTrigger) Complete(outputData interface{}) {
	//
	h.outputData = outputData

}
