package sqs

import (
	"encoding/json"
	"github.com/greenac/fifoqueue"
	"github.com/greenac/sqsmock/logger"
	"github.com/greenac/sqsmock/models"
	"github.com/greenac/sqsmock/response"
	"github.com/greenac/sqsmock/worker"
	"net/http"
	"reflect"
	"time"
)

type RequestHandler struct {
	WorkerUrls   *[]string
	Delay       int64
	q           *fifoqueue.FifoQueue
	maxMessages int
	count       int
	workerCounter int
}

func (rh *RequestHandler) setUp() {
	if rh.q == nil {
		q := fifoqueue.FifoQueue{}
		rh.q = &q
	}

	rh.maxMessages = 10
}

func (rh *RequestHandler) Add(w http.ResponseWriter, req *http.Request) {
	logger.Log("Request handler adding message")

	var m models.Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to add message to queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	m.SetIdentifiers()
	n := rh.q.Insert(m)

	logger.Log("Added new message to queue. Node:", n)

	if m.ToWorker {
		logger.Log("Sending message to worker")
		if rh.Delay > 0 {
			rh.count += 1
			c := rh.count
			logger.Log("Delaying message by:", rh.Delay, "milliseconds", "for message count:", c)
			timer := time.NewTimer(time.Millisecond * time.Duration(rh.Delay))
			go func() {
				<-timer.C
				logger.Log("Sending delayed message to worker for message count:", c)
				err = rh.sendToWorker(&m)
				if err != nil {
					logger.Warn("`RequestHandler::Add` sending message to worker with delay:", rh.Delay, "error:", err)
				}

				rh.q.Delete(n)
			}()
		} else {
			err = rh.sendToWorker(&m)
			if err != nil {
				logger.Warn("`RequestHandler::Add` sending message to worker with delay:", rh.Delay, "error:", err)
			}

			rh.q.Delete(n)
		}
	}

	logger.Log("request handler queue has:", rh.q.Length(), "nodes after add")

	pl := map[string]interface{}{"success": true, "MessageId": m.GetIdentifier()}
	rh.sendOk(w, pl)
}

func (rh *RequestHandler) Delete(w http.ResponseWriter, req *http.Request) {
	var m models.Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to retrieve messages from queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	nodes := rh.q.AsSlice()
	logger.Log("Got # of nodes in handler delete:", len(*nodes))
	var target *fifoqueue.QueueNode
	for _, n := range *nodes {
		logger.Log("node:", n)
		mm, ok := n.Payload.(models.Message)
		if !ok {
			logger.Warn("Could not get message from node. Payload was not a Message")
			continue
		}

		if mm.ReceiptHandle == m.ReceiptHandle {
			target = n
			break
		}
	}

	logger.Log("Found target:", target)

	var pl map[string]interface{}
	if target == nil {
		pl = map[string]interface{}{"success": false}
	} else {
		logger.Warn("Before delete # of queue entries:", rh.q.Length())
		suc := rh.q.Delete(target)
		logger.Warn("After delete # of queue entries:", rh.q.Length())
		pl = map[string]interface{}{"success": suc}
	}

	logger.Log("request handler queue has:", rh.q.Length(), "nodes after delete")
	rh.sendOk(w, pl)
}

func (rh *RequestHandler) RetrieveSingle(w http.ResponseWriter, req *http.Request) {
	var m models.Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to retrieve messages from queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	nodes := rh.q.AsSlice()
	var target *fifoqueue.QueueNode = nil
	for _, n := range *nodes {
		mm, ok := n.Payload.(models.Message)
		if !ok {
			logger.Warn("Could not get message from node. Payload was not a Message")
			continue
		}

		if mm.MessageId == m.MessageId {
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

func (rh *RequestHandler) Retrieve(w http.ResponseWriter, req *http.Request) {
	var m models.Message
	rh.setUp()
	err := json.NewDecoder(req.Body).Decode(&m)
	if err != nil {
		logger.Error("Failed to retrieve messages from queue with error:", err)
		rh.error(w, ResponseInternalServerError)
		return
	}

	logger.Log("Got message in handler retrieve", m.Info())

	c := rh.maxMessages
	if m.MaxNumberOfMessages != 0 && m.MaxNumberOfMessages < rh.maxMessages {
		c = m.MaxNumberOfMessages
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

func (rh *RequestHandler) sendToWorker(m *models.Message) error {
	if rh.WorkerUrls == nil {
		return nil
	}

	var workerPl interface{}
	if m.ToWorker {
		var mb map[string]interface{}
		err := json.Unmarshal([]byte(m.MessageBody), &mb)
		if err != nil {
			logger.Error("Could not unmarshal message body.", err)
			return err
		}

		workerPl = mb
	} else {
		workerPl = &m
	}

	urls := *(rh.WorkerUrls)
	url := urls[rh.workerCounter % len(urls)]
	rh.updateWorkerCounter()
	wi := worker.Interface{BaseUrl: url}
	err := wi.SendNewMessage(workerPl)
	if err != nil {
		logger.Warn("`RequestHandler::sendToWorker error sending message to:", url, err)
	}

	return nil
}

func (rh *RequestHandler) Print(w http.ResponseWriter, req *http.Request) {
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
	rr := response.Response{Error: &e, ResponseMetadata: nil}
	rr.Respond(w)
}

func (rh *RequestHandler) sendOk(w http.ResponseWriter, payload interface{}) {
	rr := response.Response{Error: nil, ResponseMetadata: &payload}
	rr.Respond(w)
}

func (rh *RequestHandler) updateWorkerCounter() {
	rh.workerCounter += 1
	if rh.workerCounter % 10000 == 0 {
		rh.workerCounter = 0
	}
}
