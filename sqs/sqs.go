package sqs

type Endpoint string

const (
	SendMessage     Endpoint = "send-message"
	DeleteMessage   Endpoint = "delete-message"
	RetrieveMessage Endpoint = "receive-message"
)
