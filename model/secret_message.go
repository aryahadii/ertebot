package model

type SecretMessage struct {
	Message          string
	SenderID         string
	SenderUsername   string
	ReceiverUsername string
	SendEpoch        int
	SeenEpoch        int
}
