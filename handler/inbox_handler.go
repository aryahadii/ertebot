package handler

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gopkg.in/mgo.v2/bson"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func HandleInboxCommand(message *botAPI.Message) string {
	if len(message.From.UserName) == 0 {
		log.WithField("User", message.From).Infoln("User doesn't have username")
		return model.NoUsernameError
	}

	var inboxMessages []model.SecretMessage
	err := db.MessagesCollection.Find(bson.M{"receiverusername": strings.ToLower(message.From.UserName)}).All(&inboxMessages)
	if err != nil {
		return model.NoSecretMessageFoundMessage
	}
	if len(inboxMessages) == 0 {
		return model.NoSecretMessageFoundMessage
	}

	resultMessage := ""
	for _, msg := range inboxMessages {
		resultMessage += fmt.Sprintf(model.InboxMessagesTemplate, msg.SenderID, msg.Message)
	}
	return resultMessage
}
