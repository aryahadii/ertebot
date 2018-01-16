package handler

import (
	"strconv"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleReplyCommand(message *botAPI.Message, callbackQuery *botAPI.CallbackQuery, user *botAPI.User, chat *botAPI.Chat) []botAPI.Chattable {
	callback := strings.Split(callbackQuery.Data, model.CallbackSeparator)
	threadOwnerID, userID, username := callback[1], callback[2], ""
	if len(callback) > 3 {
		username = callback[3]
	}

	state := &model.UserState{
		Command: model.ReplyCommand,
		Args:    []interface{}{threadOwnerID, userID, username},
	}
	userState.Set(strconv.Itoa(user.ID), *state, cache.DefaultExpiration)

	msg := botAPI.NewMessage(chat.ID, model.ReplyCommandMessageInputMessage)
	msg.ReplyMarkup = keyboard.NewBackKeyboard()
	return []botAPI.Chattable{msg}
}

func handleReplyCommandArgs(message *botAPI.Message, state model.UserState) (model.NewMessageState, *model.SecretMessage) {
	state.Args = append(state.Args, message.Text)
	userState.Set(strconv.Itoa(message.From.ID), state, cache.DefaultExpiration)

	argsLen := len(state.Args)
	if argsLen == 4 {
		secretMessage := &model.SecretMessage{
			Message:          state.Args[3].(string),
			SenderID:         strconv.Itoa(message.From.ID),
			SenderUsername:   strings.ToLower(message.From.UserName),
			ReceiverUsername: state.Args[2].(string),
			ReceiverID:       state.Args[1].(string),
			ThreadOwnerID:    state.Args[0].(string),
			SendEpoch:        time.Now().Unix(),
			SeenEpoch:        0,
		}
		log.Debugln(secretMessage)

		err := db.MessagesCollection.Insert(secretMessage)
		if err != nil {
			log.WithError(err).Errorln("Can't send message")
			return model.NewMessageStateError, nil
		}

		return model.NewMessageStateSent, secretMessage
	}

	return model.NewMessageStateError, nil
}
