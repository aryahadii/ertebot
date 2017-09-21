package handler

import (
	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleHelpCommand(message *botAPI.Message) (string, interface{}) {
	return model.HelpCommandMessage, keyboard.NewMainKeyboard()
}
