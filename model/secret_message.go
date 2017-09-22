package model

type SecretMessage struct {
	Message          string
	SenderID         string
	SenderUsername   string
	ReceiverID       string
	ReceiverUsername string
	SendEpoch        int64
	SeenEpoch        int64
}

type Message struct {
	Message          string
	SenderID         string
	SenderUsername   string
	ReceiverID       string
	ReceiverUsername string
	SendEpoch        int64
	SeenEpoch        int64
}
