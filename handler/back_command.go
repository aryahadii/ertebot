package handler

import (
	"strconv"

	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleBackCommand(message *botAPI.Message) error {
	userState.Delete(strconv.Itoa(message.From.ID))
	return nil
}
