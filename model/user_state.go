package model

type NewMessageState int

const (
	NewMessageStateInputText NewMessageState = iota
	NewMessageStateInputUsername
	NewMessageStateSent
	NewMessageStateError
)

type UserState struct {
	Command string
	Args    []interface{}
}
