package handler

import (
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
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
		receivedMessagesMap[msg.ThreadOwnerID] = append(receivedMessagesMap[msg.ThreadOwnerID], msg)

		// TODO: Make batch update
		seenMsg := msg
		seenMsg.SeenEpoch = time.Now().Unix()
		err := db.MessagesCollection.Update(msg, seenMsg)
		if err != nil {
			log.WithError(err).Errorln("Can't update seen message in DB")
		}
	}
	sortedAllMessages := util.SortInboxMessagesByTime(receivedMessagesMap)

	current_message := 0
	// Check if message update is needed
	if len(callback) > 0 {
		var err error
		current_message, err = strconv.Atoi(callback[1])
		if err != nil {
			log.WithError(err).Errorln("Can't extract current_message_index")
			current_message = 0
		}
	}

	// Create inbox's inline keyboard
	fwrdless, backless := false, false
	if len(sortedAllMessages) <= current_message+1 {
		fwrdless = true
	}
	if current_message == 0 {
		backless = true
	}
	back := model.InboxUpdateCallback + model.CallbackSeparator + strconv.Itoa(current_message-1)
	fwrd := model.InboxUpdateCallback + model.CallbackSeparator + strconv.Itoa(current_message+1)
	reply := model.InboxReplyCallback +
		model.CallbackSeparator +
		sortedAllMessages[current_message][0].ThreadOwnerID +
		model.CallbackSeparator +
		sortedAllMessages[current_message][0].SenderID +
		model.CallbackSeparator +
		sortedAllMessages[current_message][0].SenderUsername
	inboxKeyboard := keyboard.NewInboxInlineKeyboard(back, fwrd, reply, backless, fwrdless)

	// Add my messages
	current_messages := sortedAllMessages[current_message]
	var myMessages []model.SecretMessage
	db.MessagesCollection.Find(bson.M{"senderid": strconv.Itoa(user.ID), "threadownerid": current_messages[0].ThreadOwnerID, "receiverid": current_messages[0].SenderID}).All(&myMessages)
	current_messages = append(current_messages, myMessages...)
	current_messages = util.SortMessagesByTime(current_messages)

	if len(callback) > 0 {
		editMsgText := botAPI.NewEditMessageText(chat.ID, callbackQuery.Message.MessageID, util.ThreadToStringSlice(current_messages, strconv.Itoa(user.ID)))
		editReplyMarkup := botAPI.NewEditMessageReplyMarkup(chat.ID, callbackQuery.Message.MessageID, inboxKeyboard)
		return []botAPI.Chattable{editMsgText, editReplyMarkup}
	} else {
		msg := botAPI.NewMessage(chat.ID, util.ThreadToStringSlice(current_messages, strconv.Itoa(user.ID)))
		msg.ReplyMarkup = inboxKeyboard
		return []botAPI.Chattable{msg}
	}
}
