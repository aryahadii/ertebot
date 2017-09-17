package handler

import (
	"strconv"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/model"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleNewMessageArgs(message *botAPI.Message, state model.UserState) string {
	state.Args = append(state.Args, message.Text)
	userState.Set(strconv.Itoa(message.From.ID), state, cache.DefaultExpiration)

	argsLen := len(state.Args)
	if argsLen == 1 {
		return "نام کاربری فرد موردنظر را وارد کن"
	} else {
		return "در حال ارسال پیام..."
	}
}
