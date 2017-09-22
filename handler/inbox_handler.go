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

	if len(user.UserName) == 0 {
		log.WithField("User", user).Infoln("User doesn't have username")
		msg := botAPI.NewMessage(chat.ID, model.NoUsernameError)
		msg.ReplyMarkup = keyboard.NewMainKeyboard()
		return []botAPI.Chattable{msg}
	}
	var inboxMessages []model.SecretMessage
	err := db.MessagesCollection.Find(bson.M{"receiverusername": strings.ToLower(user.UserName)}).All(&inboxMessages)
	if err != nil {
		msg := botAPI.NewMessage(chat.ID, model.NoSecretMessageFoundMessage)
		msg.ReplyMarkup = keyboard.NewMainKeyboard()
		return []botAPI.Chattable{msg}
	}

	// Find messages
	resultMessagesMap := make(map[string]([]model.SecretMessage))
	for _, msg := range inboxMessages {
		resultMessagesMap[msg.SenderID] = append(resultMessagesMap[msg.SenderID], msg)

		seenMsg := msg
		seenMsg.SeenEpoch = time.Now().Unix()
		err := db.MessagesCollection.Update(msg, seenMsg)
		if err != nil {
			log.WithError(err).Errorln("Can't update seen message in DB")
		}
	}
	resultMessages := util.SortInboxMessagesByTime(resultMessagesMap)

	// No message
	if len(resultMessages) == 0 {
		msg := botAPI.NewMessage(chat.ID, model.NoSecretMessageFoundMessage)
		msg.ReplyMarkup = keyboard.NewMainKeyboard()
		return []botAPI.Chattable{msg}
	}

	current_message := 0
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
	if len(resultMessages) <= current_message+1 {
		fwrdless = true
	}
	if current_message == 0 {
		backless = true
	}
	back := model.InboxUpdateCallback + model.CallbackSeparator + strconv.Itoa(current_message-1)
	fwrd := model.InboxUpdateCallback + model.CallbackSeparator + strconv.Itoa(current_message+1)
	inboxKeyboard := keyboard.NewInboxInlineKeyboard(back, fwrd, backless, fwrdless)

	if len(callback) > 0 {
		editMsgText := botAPI.NewEditMessageText(chat.ID, callbackQuery.Message.MessageID, util.ThreadToStringSlice(resultMessages[current_message]))
		editReplyMarkup := botAPI.NewEditMessageReplyMarkup(chat.ID, callbackQuery.Message.MessageID, inboxKeyboard)
		return []botAPI.Chattable{editMsgText, editReplyMarkup}
	} else {
		msg := botAPI.NewMessage(chat.ID, util.ThreadToStringSlice(resultMessages[current_message]))
		msg.ReplyMarkup = inboxKeyboard
		return []botAPI.Chattable{msg}
	}
}
