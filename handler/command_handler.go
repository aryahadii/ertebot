package handler

import (
	"strconv"
	"time"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gopkg.in/mgo.v2/bson"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleCommand(message *botAPI.Message) string {
	person := &model.Person{}
	err := db.PeopleCollection.Find(bson.M{"UserID": message.From.ID}).One(person)
	if err != nil {
		person = &model.Person{
			UserID:       strconv.Itoa(message.From.ID),
			FirstName:    message.From.FirstName,
			LastName:     message.From.LastName,
			Username:     message.From.UserName,
			LastUseEpoch: time.Now().Unix(),
		}

		db.PeopleCollection.Insert(person)
	} else {
		person.LastUseEpoch = time.Now().Unix()
		db.PeopleCollection.Update(&model.Person{Username: person.Username}, person)
	}

	if message.Command() == "newmessage" {
		state := &model.UserState{
			Command: "newmessage",
		}
		userState.Set(strconv.Itoa(message.From.ID), *state, cache.DefaultExpiration)
		return "متن پیام را وارد کنید"
	}

	return "دستور به درستی وارد نشده"
}
