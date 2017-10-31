package webserver

import (
	"../model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Wrapper for the data we send to the UI
// For now only contains the docker stacks
type WebserverData struct {
	Stacks []model.Stack
}

// Replaces the stacks currently held by the webserverdata
func (wd *WebserverData) UpdateStacks(stacks []model.Stack) {
	wd.Stacks = stacks
}

// Handle function for http router
func (wd *WebserverData) HandleGetStacks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(wd.Stacks)
}

// Wrapper object for managing the http router and logger
type Server struct {
	logger *log.Logger
	mux    *http.ServeMux
}

// Starts the webserver on the port given
// The channel bool is for telling the server to shutdown
func StartServer(port string, data *WebserverData, c chan bool) {
	router := mux.NewRouter()
	router.HandleFunc("/stacks", data.HandleGetStacks).Methods("GET")
	listenAddress := fmt.Sprintf(":%s", port)
	server := &http.Server{Addr: listenAddress, Handler: router}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	if b := <-c; b {
		fmt.Printf("We got told to quit\n")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		server.Shutdown(ctx)
		cancel()
	}
	c <- true
}
