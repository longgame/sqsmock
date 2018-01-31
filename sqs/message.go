package sqs

import (
	"crypto/md5"
)

type Message struct {
	MessageId *string `json:"MessageId"`
	MessageBody *[]byte `json:"MessageBody"`
	QueueUrl string `json:"QueueUrl"`
	MessageGroupId string `json:"MessageGroupId"`
	MaxNumberOfMessages *int `json:"MaxNumberOfMessages"`
}

func (m *Message)makeIdentifier() string {
	h := md5.New()
	if m.MessageBody == nil {
		return "NO_BODY"
	}

	return string(h.Sum(m.*MessageBody))
}

func (m *Message)SetIdentifier() {
	mId := m.makeIdentifier()
	m.MessageId = &mId
}
