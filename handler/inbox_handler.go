package handler

import (
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	"gitlab.com/arha/Ertebot/util"
	"gopkg.in/mgo.v2/bson"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleInboxCommand(message *botAPI.Message, callbackQuery *botAPI.CallbackQuery, user *botAPI.User, chat *botAPI.Chat) []botAPI.Chattable {
	var callback []string
	if callbackQuery != nil {
		callback = strings.Split(callbackQuery.Data, model.CallbackSeparator)
	}

	var allMessages []model.SecretMessage
	db.MessagesCollection.Find(bson.M{
		"$or": []interface{}{
			bson.M{"receiverusername": strings.ToLower(user.UserName)},
			bson.M{"receiverid": strconv.Itoa(user.ID)},
		},
	}).All(&allMessages)

	if len(allMessages) == 0 {
		msg := botAPI.NewMessage(chat.ID, model.NoSecretMessageFoundMessage)
		msg.ReplyMarkup = keyboard.NewMainKeyboard()
		return []botAPI.Chattable{msg}
	}

	// Sort threads
	receivedMessagesMap := make(map[string]([]model.SecretMessage))
	for _, msg := range allMessages {
		key := util.GetMD5(msg.ThreadOwnerID + msg.SenderID + msg.SenderUsername)
		receivedMessagesMap[key] = append(receivedMessagesMap[key], msg)

		// TODO: Make batch update
		seenMsg := msg
		seenMsg.SeenEpoch = time.Now().Unix()
		err := db.MessagesCollection.Update(msg, seenMsg)
		if err != nil {
			log.WithError(err).Errorln("Can't update seen message in DB")
		}
	}
	sortedAllMessages := util.SortInboxMessagesByTime(receivedMessagesMap)

	currentMessage := 0
	// Check if message update is needed
	if len(callback) > 0 {
		var err error
		currentMessage, err = strconv.Atoi(callback[1])
		if err != nil {
			log.WithError(err).Errorln("Can't extract currentMessage_index")
			currentMessage = 0
		}
	}

	// Create inbox's inline keyboard
	fwrdless, backless := false, false
	if len(sortedAllMessages) <= currentMessage+1 {
		fwrdless = true
	}
	if currentMessage == 0 {
		backless = true
	}
	back := model.InboxUpdateCallback + model.CallbackSeparator + strconv.Itoa(currentMessage-1)
	fwrd := model.InboxUpdateCallback + model.CallbackSeparator + strconv.Itoa(currentMessage+1)
	reply := model.InboxReplyCallback +
		model.CallbackSeparator +
		sortedAllMessages[currentMessage][0].ThreadOwnerID +
		model.CallbackSeparator +
		sortedAllMessages[currentMessage][0].SenderID +
		model.CallbackSeparator +
		sortedAllMessages[currentMessage][0].SenderUsername
	inboxKeyboard := keyboard.NewInboxInlineKeyboard(back, fwrd, reply, backless, fwrdless)

	// Add my messages
	currentMessages := sortedAllMessages[currentMessage]
	var myMessages []model.SecretMessage
	db.MessagesCollection.Find(bson.M{
		"senderid":      strconv.Itoa(user.ID),
		"threadownerid": currentMessages[0].ThreadOwnerID,
		"$or": []interface{}{
			bson.M{"receiverid": "", "receiverusername": currentMessages[0].SenderUsername},
			bson.M{"receiverid": currentMessages[0].SenderID},
		},
	}).All(&myMessages)

	currentMessages = append(currentMessages, myMessages...)
	currentMessages = util.SortMessagesByTime(currentMessages)

	if len(callback) > 0 {
		editMsgText := botAPI.NewEditMessageText(chat.ID, callbackQuery.Message.MessageID, util.ThreadToStringSlice(currentMessages, strconv.Itoa(user.ID)))
		editReplyMarkup := botAPI.NewEditMessageReplyMarkup(chat.ID, callbackQuery.Message.MessageID, inboxKeyboard)
		return []botAPI.Chattable{editMsgText, editReplyMarkup}
	}
	msg := botAPI.NewMessage(chat.ID, util.ThreadToStringSlice(currentMessages, strconv.Itoa(user.ID)))
	msg.ReplyMarkup = inboxKeyboard
	return []botAPI.Chattable{msg}
}
