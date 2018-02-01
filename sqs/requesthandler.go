package sqs

import (
	"github.com/greenac/fifoqueue"
	"net/http"
	"encoding/json"
	"github.com/greenac/sqsmock/logger"
	"github.com/greenac/restresponse"
)

type RequestHandler struct {
	q *fifoqueue.FifoQueue
}

func (rh *RequestHandler) setUp() {
	if rh.q == nil {
		q := fifoqueue.FifoQueue{}
		rh.q = &q
	}
}

func (rh *RequestHandler) Add(w http.ResponseWriter, req *http.Request){
	var m Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to add message to queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	logger.Log("Adding new message to queue:", m)
	rh.q.Insert(m)
}

func (rh *RequestHandler) Retrieve(w http.ResponseWriter, req *http.Request){
	var m Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to retrieve messages to queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	logger.Log("Adding new message to queue:", m)
	rh.q.Insert(m)
}

func (rh *RequestHandler) error(w http.ResponseWriter, rc ResponseCode) {
	rr := restresponse.Response{Code: int(rc), Payload: nil}
	rr.Respond(w)
}

