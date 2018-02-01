package sqs

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/fatih/structs"
)

type Message struct {
	MessageId           *string `json:"MessageId"`
	MessageBody         string  `json:"MessageBody"`
	QueueUrl            string  `json:"QueueUrl"`
	MessageGroupId      *string `json:"MessageGroupId"`
	MaxNumberOfMessages *int    `json:"MaxNumberOfMessages"`
}

func (m *Message) makeIdentifier() string {
	h := md5.New()
	id := h.Sum([]byte(m.MessageBody))
	hid := hex.EncodeToString(id)
	return hid
}

func (m *Message) SetIdentifier() {
	mId := m.makeIdentifier()
	m.MessageId = &mId
}

func (m *Message) GetIdentifier() *string {
	if m.MessageId != nil {
		return m.MessageId
	}

	m.SetIdentifier()
	return m.MessageId
}

func (m *Message) Info() map[string]interface{} {
	return structs.Map(m)
}
