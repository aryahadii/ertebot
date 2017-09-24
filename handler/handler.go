package handler

import (
	"strconv"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gopkg.in/mgo.v2/bson"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

var (
	userState = cache.New(30*time.Minute, 10*time.Minute)
)

func HandleMessage(message *botAPI.Message) []botAPI.Chattable {
	updateUseEpoch(message)

	msg := botAPI.NewMessage(message.Chat.ID, model.SomeErrorOccured)

	if message.IsCommand() {
		msg.Text, msg.ReplyMarkup = handleCommand(message)
	} else if message.Text == model.BackCommand {
		msg.Text, msg.ReplyMarkup = handleBackCommand(message)
	} else if message.Text == model.InboxCommand {
		return handleInboxCommand(message, nil, message.From, message.Chat)
	} else if message.Text == model.LinkCommand {
		msg.Text, msg.ReplyMarkup = handleMyLinkCommand(message)
	} else if message.Text == model.HelpCommand {
		msg.Text, msg.ReplyMarkup = handleHelpCommand(message)
	} else if message.Text == model.NewMessageCommand {
		msg.Text, msg.ReplyMarkup = handleNewMessage(message)
	} else {
		state, found := userState.Get(strconv.Itoa(message.From.ID))
		if found {
			if (state.(model.UserState)).Command == model.NewMessageCommand {
				msg.Text, msg.ReplyMarkup = handleNewMessageArgs(message, state.(model.UserState))
			} else if (state.(model.UserState)).Command == model.ReplyCommand {
				msg.Text, msg.ReplyMarkup = handleReplyCommandArgs(message, state.(model.UserState))
			} else if (state.(model.UserState)).Command == model.NewMessageByLinkCommand {
				msg.Text, msg.ReplyMarkup = handleNewMessageByLinkArgs(message, state.(model.UserState))
			}
		}
	}

	return []botAPI.Chattable{msg}
}

func HandleCallback(callbackQuery *botAPI.CallbackQuery) []botAPI.Chattable {
	var callback []string
	if callbackQuery != nil {
		callback = strings.Split(callbackQuery.Data, model.CallbackSeparator)
	}

	if callback[0] == model.InboxUpdateCallback {
		return handleInboxCommand(nil, callbackQuery, callbackQuery.From, callbackQuery.Message.Chat)
	} else if callback[0] == model.InboxReplyCallback {
		return handleReplyCommand(nil, callbackQuery, callbackQuery.From, callbackQuery.Message.Chat)
	}

	return []botAPI.Chattable{}
}

func updateUseEpoch(message *botAPI.Message) {
	// Update LastUseEpoch or Create new user if it's needed
	person := &model.Person{}
	err := db.PeopleCollection.Find(bson.M{"userid": strconv.Itoa(message.From.ID)}).One(person)
	if err != nil {
		person = &model.Person{
			UserID:       strconv.Itoa(message.From.ID),
			FirstName:    message.From.FirstName,
			LastName:     message.From.LastName,
			Username:     strings.ToLower(message.From.UserName),
			HashID:       util.GetHashID(stconv.Itoa(message.From.ID)),
			LastUseEpoch: time.Now().Unix(),
		}

		db.PeopleCollection.Insert(person)
	} else {
		person.LastUseEpoch = time.Now().Unix()
		db.PeopleCollection.Update(&model.Person{Username: person.Username}, person)
	}
}
