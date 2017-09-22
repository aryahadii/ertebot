package keyboard

import (
	"gitlab.com/arha/Ertebot/model"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func NewMainKeyboard() botAPI.ReplyKeyboardMarkup {
	newMessageKey := botAPI.NewKeyboardButton(model.NewMessageCommand)
	inboxKey := botAPI.NewKeyboardButton(model.InboxCommand)
	row1 := botAPI.NewKeyboardButtonRow(newMessageKey, inboxKey)

	helpKey := botAPI.NewKeyboardButton(model.HelpCommand)
	linkKey := botAPI.NewKeyboardButton(model.LinkCommand)
	row2 := botAPI.NewKeyboardButtonRow(helpKey, linkKey)

	return botAPI.NewReplyKeyboard(row1, row2)
}

func NewBackKeyboard() botAPI.ReplyKeyboardMarkup {
	backKey := botAPI.NewKeyboardButton(model.BackCommand)
	row := botAPI.NewKeyboardButtonRow(backKey)
	return botAPI.NewReplyKeyboard(row)
}

func NewInboxInlineKeyboard(back, fwrd, reply string, backless, fwrdless bool) botAPI.InlineKeyboardMarkup {
	var row []botAPI.InlineKeyboardButton
	if !backless {
		row = append(row, botAPI.NewInlineKeyboardButtonData("<", back))
	}
	row = append(row, botAPI.NewInlineKeyboardButtonData("پاسخ", reply))
	if !fwrdless {
		row = append(row, botAPI.NewInlineKeyboardButtonData(">", fwrd))
	}

	return botAPI.NewInlineKeyboardMarkup(row)
}
