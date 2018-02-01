package main

import (
	"github.com/greenac/sqsmock/logger"
	"github.com/greenac/sqsmock/sqs"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/greenac/fifoqueue"
)

const ServerUrl = ":4242"

func main() {
	q := fifoqueue.FifoQueue{}
	for i := 0; i < 10; i++ {
		q.Insert(i)
	}

	qs := q.AsSlice()
	for _, n := range *qs {
		if n.Payload == 5 {
			logger.Log("found payload", n.Payload, "deleting")
			q.Delete(n)
			nqs := q.AsSlice()
			for _, nn := range *nqs {
				logger.Log(*nn)
			}
		}
	}

	logger.Log("Running sqs mock server on:", ServerUrl)
	handler := sqs.RequestHandler{}
	router := mux.NewRouter()
	router.HandleFunc(sqs.AddMessageEndpoint, handler.Add).Methods("POST")
	router.HandleFunc(sqs.RetrieveMessagesEndpoint, handler.Retrieve).Methods("POST")
	logger.Error(http.ListenAndServe(ServerUrl, router))
}
