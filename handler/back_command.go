package handler

import (
	"strconv"

	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleBackCommand(message *botAPI.Message) (string, interface{}) {
	userState.Delete(strconv.Itoa(message.From.ID))

	return model.BackCommandMessage, keyboard.NewMainKeyboard()
}
