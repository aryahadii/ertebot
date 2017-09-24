package handler

import (
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	"gitlab.com/arha/Ertebot/util"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleNewMessage(message *botAPI.Message) (string, interface{}) {
	state := &model.UserState{
		Command: model.NewMessageCommand,
	}
	userState.Set(strconv.Itoa(message.From.ID), *state, cache.DefaultExpiration)
	return model.NewMessageCommandMessageInputMessage, keyboard.NewBackKeyboard()
}

func handleNewMessageArgs(message *botAPI.Message, state model.UserState) (string, interface{}) {
	state.Args = append(state.Args, message.Text)
	userState.Set(strconv.Itoa(message.From.ID), state, cache.DefaultExpiration)

	argsLen := len(state.Args)
	if argsLen == 1 {
		return model.NewMessageCommandUsernameMessage, keyboard.NewBackKeyboard()
	} else {
		receiverUsername := strings.ToLower(strings.TrimLeft(state.Args[1].(string), "@"))
		id, err := util.GetUserID(receiverUsername)
		if err != nil {
			log.WithField("Receiver", receiverUsername).Debugln("Doesn't have userID")
		}

		secretMessage := &model.SecretMessage{
			Message:          state.Args[0].(string),
			SenderID:         strconv.Itoa(message.From.ID),
			SenderUsername:   strings.ToLower(message.From.UserName),
			ReceiverUsername: receiverUsername,
			ReceiverID:       id,
			ThreadOwnerID:    strconv.Itoa(message.From.ID),
			SendEpoch:        time.Now().Unix(),
			SeenEpoch:        0,
		}
		log.Debugln(secretMessage)

		err = db.MessagesCollection.Insert(secretMessage)
		if err != nil {
			log.WithError(err).Errorln("Can't send message")
			return model.NewMessageCommandSendErrorMessage, keyboard.NewMainKeyboard()
		}

		return model.NewMessageCommandSentMessage, keyboard.NewMainKeyboard()
	}
}

func handleNewMessageByLink(message *botAPI.Message) (string, interface{}) {
	state := &model.UserState{
		Command: model.NewMessageByLinkCommand,
		Args:    []interface{}{message.CommandArguments()},
	}
	userState.Set(strconv.Itoa(message.From.ID), *state, cache.DefaultExpiration)
	return model.NewMessageCommandMessageInputMessage, keyboard.NewBackKeyboard()
}

func handleNewMessageByLinkArgs(message *botAPI.Message, state model.UserState) (string, interface{}) {
	state.Args = append(state.Args, message.Text)
	userState.Set(strconv.Itoa(message.From.ID), state, cache.DefaultExpiration)

	argsLen := len(state.Args)
	if argsLen == 2 {
		receiverHashID := state.Args[0].(string)
		person, err := util.GetPersonByHashID(receiverHashID)
		if err != nil {
			log.WithField("Receiver HashID", receiverHashID).Debugln("Doesn't have userID")
		}

		secretMessage := &model.SecretMessage{
			Message:          state.Args[1].(string),
			SenderID:         strconv.Itoa(message.From.ID),
			SenderUsername:   strings.ToLower(message.From.UserName),
			ReceiverUsername: person.Username,
			ReceiverID:       person.UserID,
			ThreadOwnerID:    strconv.Itoa(message.From.ID),
			SendEpoch:        time.Now().Unix(),
			SeenEpoch:        0,
		}
		log.Debugln(secretMessage)

		err = db.MessagesCollection.Insert(secretMessage)
		if err != nil {
			log.WithError(err).Errorln("Can't send message")
			return model.NewMessageCommandSendErrorMessage, keyboard.NewMainKeyboard()
		}

		return model.NewMessageCommandSentMessage, keyboard.NewMainKeyboard()
	}
	return model.SomeErrorOccured, keyboard.NewMainKeyboard()
}
