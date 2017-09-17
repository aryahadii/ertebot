package handler

import (
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/model"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

var (
	userState = cache.New(30*time.Minute, 10*time.Minute)
)

func Handle(message *botAPI.Message) botAPI.MessageConfig {
	msg := botAPI.NewMessage(message.Chat.ID, "خطایی رخ داد! دوباره دستور بده")

	if message.IsCommand() {
		msg.Text = handleCommand(message)
	} else {
		state, found := userState.Get(strconv.Itoa(message.From.ID))
		if !found {
			return msg
		}
		if (state.(model.UserState)).Command == "newmessage" {
			msg.Text = handleNewMessageArgs(message, state.(model.UserState))
		}
	}

	return msg
}
