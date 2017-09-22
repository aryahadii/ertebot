package handler

import (
	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleCommand(message *botAPI.Message) (string, interface{}) {
	// Handle commands
	if message.Command() == model.StartCommand {
		return model.WelcomeMessage, nil
	}
	if message.Command() == model.NewMessageRawCommand {
		return handleNewMessage(message)
	}
	//	if message.Command() == model.InboxRawCommand {
	//		return handleInboxCommand(message, []string{})
	//	}
	if message.Command() == model.HelpRawCommand {
		return handleHelpCommand(message)
	}

	return model.WrongCommandMessage, keyboard.NewMainKeyboard()
}
