package main

import (
	"github.com/greenac/gologger"
	"github.com/greenac/sqsmock/sqs"
	"github.com/gorilla/mux"
	"net/http"
)

var logger = gologger.GoLogger{}

func main() {
	handler := sqs.RequestHandler{}
	router := mux.NewRouter()
	router.HandleFunc(sqs.AddMessageEndpoint, handler.Add).Methods("POST")
	router.HandleFunc(sqs.RetrieveMessagesEndpoint, handler.Retrieve).Methods("POST")
	logger.Error(http.ListenAndServe(":4242", router))
}

