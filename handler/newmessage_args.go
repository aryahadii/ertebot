package handler

import (
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/db"
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
		secretMessage := &model.SecretMessage{
			Message:          state.Args[0].(string),
			SenderID:         strconv.Itoa(message.From.ID),
			SenderUsername:   message.From.UserName,
			ReceiverUsername: strings.TrimLeft(state.Args[1].(string), "@"),
			SendEpoch:        time.Now().Unix(),
			SeenEpoch:        0,
		}
		err := db.MessagesCollection.Insert(secretMessage)
		if err != nil {
			log.WithError(err).Errorln("Can't send message")
			return "خطایی پیش آمد! دوباره تلاش کنید"
		}

		return "پیام ارسال شد"
	}
}
