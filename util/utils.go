package util

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gopkg.in/mgo.v2/bson"
)

func GetUserID(username string) (string, error) {
	var id string
	err := db.PeopleCollection.Find(bson.M{"username": strings.ToLower(username)}).One(&id)
	if err != nil {
		return id, errors.New("Not found")
	}
	return id, nil
}

func SortInboxMessagesByTime(messages map[string]([]model.SecretMessage)) []([]model.SecretMessage) {
	messagesSlice := make([]([]model.SecretMessage), 0)
	for _, v := range messages {
		SortMessagesByTime(v)
		messagesSlice = append(messagesSlice, v)
	}

	sort.Sort(model.ThreadNewFirst(messagesSlice))

	return messagesSlice
}

func SortMessagesByTime(messages []model.SecretMessage) []model.SecretMessage {
	sort.Sort(model.SecretMessageNewFirst(messages))
	return messages
}

func ThreadToStringSlice(thread []model.SecretMessage, myID string) string {
	threadMessages := ""
	for i, _ := range thread {
		message := thread[len(thread)-i-1]
		if myID == message.SenderID {
			threadMessages += fmt.Sprintf(model.InboxMessagesTemplate, "من", message.Message)
		} else {
			if myID == message.ThreadOwnerID {
				babeName := message.SenderUsername
				if len(babeName) == 0 {
					babeName = "او"
				}
				threadMessages += fmt.Sprintf(model.InboxMessagesTemplate, babeName, message.Message)
			} else {
				threadMessages += fmt.Sprintf(model.InboxMessagesTemplate, message.SenderID, message.Message)
			}
		}
	}
	return threadMessages
}
