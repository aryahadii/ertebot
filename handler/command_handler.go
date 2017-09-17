package handler

import (
	"strconv"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/model"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleCommand(message *botAPI.Message) string {
	if message.Command() == "newmessage" {
		state := &model.UserState{
			Command: "newmessage",
		}
		userState.Set(strconv.Itoa(message.From.ID), *state, cache.DefaultExpiration)
		return "متن پیام را وارد کنید"
	}

	return "دستور به درستی وارد نشده"
}
