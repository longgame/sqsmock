package sqs

import (
	"github.com/greenac/fifoqueue"
	"net/http"
	"encoding/json"
	"github.com/greenac/sqsmock/logger"
	"github.com/greenac/sqsmock/response"
	"github.com/greenac/sqsmock/models"
	"reflect"
	"github.com/greenac/sqsmock/worker"
)


type RequestHandler struct {
	WorkerUrl string
	q *fifoqueue.FifoQueue
	maxMessages int
}

func (rh *RequestHandler) setUp() {
	if rh.q == nil {
		q := fifoqueue.FifoQueue{}
		rh.q = &q
	}

	rh.maxMessages = 10
}

func (rh *RequestHandler) Add(w http.ResponseWriter, req *http.Request){
	var m models.Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to add message to queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	m.SetIdentifier()
	rh.q.Insert(m)

	logger.Log("Added new message to queue:", m.Info())

	wi := worker.Interface{BaseUrl: rh.WorkerUrl}
	wi.SendNewMessage(&m)

	pl := map[string]interface{}{"success": true, "MessageId": m.GetIdentifier()}
	rh.sendOk(w, pl)
}

func (rh *RequestHandler) Delete(w http.ResponseWriter, req *http.Request){
	var m models.Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to retrieve messages from queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	logger.Log("Got body to retrieving message", m.Info())
	nodes := rh.q.AsSlice()
	var target *fifoqueue.QueueNode = nil
	for _, n := range *nodes {
		mm, ok := n.Payload.(models.Message)
		if !ok {
			logger.Warn("Could not get message from node. Payload was not a Message")
			continue
		}

		if *(mm.MessageId) == *(m.MessageId) {
			target = n
			break
		}
	}

	var pl map[string]interface{}
	if target == nil {
		pl = map[string]interface{}{"success": false}
	} else {
		suc := rh.q.Delete(target)
		pl = map[string]interface{}{"success": suc}
	}

	rh.sendOk(w, pl)
}

func (rh *RequestHandler) RetrieveSingle(w http.ResponseWriter, req *http.Request){
	var m models.Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to retrieve messages from queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	logger.Log("Got body to retrieving message", m.Info())
	nodes := rh.q.AsSlice()
	var target *fifoqueue.QueueNode = nil
	for _, n := range *nodes {
		mm, ok := n.Payload.(models.Message)
		if !ok {
			logger.Warn("Could not get message from node. Payload was not a Message")
			continue
		}

		if *(mm.MessageId) == *(m.MessageId) {
			target = n
			break
		}
	}

	var pl map[string]interface{}
	if target == nil {
		pl = map[string]interface{}{"message": nil}
	} else {
		pl = map[string]interface{}{"message": target.Payload}
	}

	rh.sendOk(w, pl)
}

func (rh *RequestHandler) Retrieve(w http.ResponseWriter, req *http.Request){
	var m models.Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to retrieve messages from queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	logger.Log("Got body to retrieving message", m.Info())

	c := rh.maxMessages
	if m.MaxNumberOfMessages != nil && *(m.MaxNumberOfMessages) < rh.maxMessages {
		c = *(m.MaxNumberOfMessages)
	}

	payloads := rh.q.GetPayloads(c)
	messages := make([]models.Message, len(*payloads))
	for i, pl := range *payloads {
		mm, ok := pl.(models.Message)
		if !ok {
			logger.Warn("Could not get message from payload. Payload not a Message. It is of type:", reflect.TypeOf(pl))
			continue
		}
		messages[i] = mm
	}

	pl := map[string]interface{}{"Messages": messages}
	rh.sendOk(w, pl)
}

func (rh *RequestHandler) Print(w http.ResponseWriter, req *http.Request){
	rh.setUp()
	nodes := rh.q.AsSlice()
	for i, n := range *nodes {
		logger.Log("node:", i, "=", *n, "\n")
	}

	logger.Log("Finished printing mock sqs queue for:", len(*nodes), "nodes")

	pl := map[string]interface{}{"success": true}
	rh.sendOk(w, pl)
}

func (rh *RequestHandler) error(w http.ResponseWriter, rc ResponseCode) {
	e := map[string]interface{}{"code": rc, "message": "Sqs mock failed with a generic error"}
	rr := response.Response{Error: &e, Data: nil}
	rr.Respond(w)
}

func (rh *RequestHandler) sendOk(w http.ResponseWriter, payload interface{}) {
	rr := response.Response{Error: nil, Data: &payload}
	rr.Respond(w)
}

