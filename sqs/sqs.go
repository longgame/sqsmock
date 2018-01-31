package sqs

type Endpoint string

const (
	AddMessage	Endpoint = "add-message"
	DeleteMessage Endpoint = "delete-message"
	GetMessages Endpoint = "get-messages"
)

