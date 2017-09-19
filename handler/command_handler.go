package handler

import (
	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gopkg.in/mgo.v2/bson"
	botAPI "gopkg.in/telegram-bot-api.v4"
	"strconv"
	"strings"
	"time"
)

func handleCommand(message *botAPI.Message) string {
	// Update LastUseEpoch or Create new user if it's needed
	person := &model.Person{}
	err := db.PeopleCollection.Find(bson.M{"UserID": message.From.ID}).One(person)
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

	// Handle commands
	if message.Command() == model.StartCommand {
		return model.WelcomeMessage
	}
	if message.Command() == model.NewMessageCommand {
		state := &model.UserState{
			Command: model.NewMessageCommand,
		}
		userState.Set(strconv.Itoa(message.From.ID), *state, cache.DefaultExpiration)
		return "متن پیام را وارد کنید"
	}
	if message.Command() == model.InboxCommand {
		return HandleInboxCommand(message)
	}
	if message.Command() == model.HelpCommand {
		return HandleHelpCommand(message)
	}

	return "دستور به درستی وارد نشده"
}
