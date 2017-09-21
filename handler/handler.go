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

func Handle(message *botAPI.Message) botAPI.MessageConfig {
	updateUseEpoch(message)

	msg := botAPI.NewMessage(message.Chat.ID, model.SomeErrorOccured)

	if message.IsCommand() {
		msg.Text, msg.ReplyMarkup = handleCommand(message)
	} else if message.Text == model.BackCommand {
		msg.Text, msg.ReplyMarkup = handleBackCommand(message)
	} else if message.Text == model.InboxCommand {
		msg.Text, msg.ReplyMarkup = handleInboxCommand(message)
	} else if message.Text == model.HelpCommand {
		msg.Text, msg.ReplyMarkup = handleHelpCommand(message)
	} else if message.Text == model.NewMessageCommand {
		msg.Text, msg.ReplyMarkup = handleNewMessage(message)
	} else {
		state, found := userState.Get(strconv.Itoa(message.From.ID))
		if found {
			if (state.(model.UserState)).Command == model.NewMessageCommand {
				msg.Text, msg.ReplyMarkup = handleNewMessageArgs(message, state.(model.UserState))
			}
		}
	}

	return msg
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
			LastUseEpoch: time.Now().Unix(),
		}

		db.PeopleCollection.Insert(person)
	} else {
		person.LastUseEpoch = time.Now().Unix()
		db.PeopleCollection.Update(&model.Person{Username: person.Username}, person)
	}
}
