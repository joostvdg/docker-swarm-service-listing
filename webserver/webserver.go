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

type WebserverData struct {
	Stacks []model.Stack
}

func (wd *WebserverData) UpdateStacks(stacks []model.Stack) {
	wd.Stacks = stacks
}

func (wd *WebserverData) HandleGetStacks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(wd.Stacks)
}

type Server struct {
	logger *log.Logger
	mux    *http.ServeMux
}

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
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		server.Shutdown(ctx)
	}
	c <- true
}
