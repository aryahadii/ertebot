package handler

import (
	"gitlab.com/arha/Ertebot/model"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func HandleHelpCommand(message *botAPI.Message) string {
	return model.HelpCommandMessage
}
