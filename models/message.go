package models

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/fatih/structs"
)

type Message struct {
	MessageId           string `json:"MessageId"`
	MessageBody         string  `json:"MessageBody"`
	QueueUrl            string  `json:"QueueUrl"`
	MessageGroupId      string `json:"MessageGroupId"`
	MaxNumberOfMessages int    `json:"MaxNumberOfMessages"`
	ReceiptHandle 			string `json:"ReceiptHandle"`
}

func (m *Message) makeIdentifier() string {
	h := md5.New()
	id := h.Sum([]byte(m.MessageBody))
	hid := hex.EncodeToString(id)
	return hid
}

func (m *Message) SetIdentifier() {
	m.MessageId = m.makeIdentifier()
}

func (m *Message) GetIdentifier() string {
	if m.MessageId != "" {
		return m.MessageId
	}

	m.SetIdentifier()
	return m.MessageId
}

func (m *Message) SetReceiptHandle() {
	m.ReceiptHandle = m.makeIdentifier()
}

func (m *Message) GetReceiptHandle() string {
	if m.ReceiptHandle != "" {
		return m.ReceiptHandle
	}

	m.SetReceiptHandle()
	return m.ReceiptHandle
}

func (m *Message) SetIdentifiers() {
	m.SetIdentifier()
	m.SetReceiptHandle()
}

func (m *Message) Info() map[string]interface{} {
	return structs.Map(m)
}
