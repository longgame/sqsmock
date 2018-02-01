package main

import (
	"github.com/gorilla/mux"
	"github.com/greenac/sqsmock/logger"
	"github.com/greenac/sqsmock/sqs"
	"net/http"
	"github.com/joho/godotenv"
	"os"
)

const ServerUrl = "127.0.0.1:4242"

func main() {
	err := godotenv.Load()
	workerUrl := ""
	if err == nil {
		workerUrl = os.Getenv("LG_NOTIFICATIONS_WORKER_URL")
		logger.Log("Loaded .env file...worker url:", workerUrl)
	} else {
		logger.Error("Error loading .env file")
	}



	logger.Log("Running sqs mock server on:", ServerUrl)

	handler := sqs.RequestHandler{WorkerUrl: workerUrl}
	router := mux.NewRouter()
	router.HandleFunc(sqs.AddMessageEndpoint, handler.Add).Methods("POST")
	router.HandleFunc(sqs.RetrieveMessageEndpoint, handler.Retrieve).Methods("POST")
	router.HandleFunc(sqs.DeleteMessageEndpoint, handler.Delete).Methods("POST")
	router.HandleFunc(sqs.PrintQueueEndpoint, handler.Print).Methods("GET")
	logger.Error(http.ListenAndServe(ServerUrl, router))
}
