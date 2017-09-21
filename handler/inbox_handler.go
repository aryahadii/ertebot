package handler

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	"gopkg.in/mgo.v2/bson"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleInboxCommand(message *botAPI.Message) (string, interface{}) {
	if len(message.From.UserName) == 0 {
		log.WithField("User", message.From).Infoln("User doesn't have username")
		return model.NoUsernameError, keyboard.NewMainKeyboard()
	}

	var inboxMessages []model.SecretMessage
	err := db.MessagesCollection.Find(bson.M{"receiverusername": strings.ToLower(message.From.UserName)}).All(&inboxMessages)
	if err != nil {
		return model.NoSecretMessageFoundMessage, keyboard.NewMainKeyboard()
	}

	resultMessage := ""
	for _, msg := range inboxMessages {
		if msg.SeenEpoch != 0 {
			continue
		}
		resultMessage += fmt.Sprintf(model.InboxMessagesTemplate, msg.SenderID, msg.Message)

		seenMsg := msg
		seenMsg.SeenEpoch = time.Now().Unix()
		err := db.MessagesCollection.Update(msg, seenMsg)
		if err != nil {
			log.WithError(err).Errorln("Can't update seen message in DB")
		}
	}

	if len(resultMessage) == 0 {
		return model.NoSecretMessageFoundMessage, keyboard.NewMainKeyboard()
	}
	return resultMessage, keyboard.NewMainKeyboard()
}
