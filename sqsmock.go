package main

import (
	"github.com/gorilla/mux"
	"github.com/greenac/sqsmock/logger"
	"github.com/greenac/sqsmock/sqs"
	"net/http"
)

const ServerUrl = "127.0.0.1:4242"

func main() {
	logger.Log("Running sqs mock server on:", ServerUrl)
	handler := sqs.RequestHandler{}
	router := mux.NewRouter()
	router.HandleFunc(sqs.AddMessageEndpoint, handler.Add).Methods("POST")
	router.HandleFunc(sqs.RetrieveMessageEndpoint, handler.Retrieve).Methods("POST")
	router.HandleFunc(sqs.DeleteMessageEndpoint, handler.Delete).Methods("POST")
	router.HandleFunc(sqs.PrintQueueEndpoint, handler.Print).Methods("GET")
	logger.Error(http.ListenAndServe(ServerUrl, router))
}
