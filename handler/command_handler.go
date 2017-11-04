package handler

import (
	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleCommand(message *botAPI.Message) (string, interface{}) {
	// Handle commands
	if message.Command() == model.StartCommand {
		if len(message.CommandArguments()) > 0 {
			return handleNewMessageByLink(message)
		}
		return model.WelcomeMessage, nil
	}
	if message.Command() == model.HelpRawCommand {
		return handleHelpCommand(message)
	}

	return model.WrongCommandMessage, keyboard.NewMainKeyboard()
}
