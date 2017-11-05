package handler

import (
	"strconv"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"gitlab.com/arha/Ertebot/db"
	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	"gitlab.com/arha/Ertebot/util"
	"gopkg.in/mgo.v2/bson"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

var (
	userState = cache.New(30*time.Minute, 10*time.Minute)
)

func HandleMessage(message *botAPI.Message) []botAPI.Chattable {
	updateUseEpoch(message)

	answerMessage := botAPI.NewMessage(message.Chat.ID, model.SomeErrorOccured)
	answerMessage.ReplyMarkup = keyboard.NewMainKeyboard()
	var receiverNotifyMessage *botAPI.MessageConfig

	if message.IsCommand() {
		if message.Command() == model.StartCommand {
			if len(message.CommandArguments()) > 0 {
				handleNewMessageByLink(message)
				answerMessage.Text = model.NewMessageCommandMessageInputMessage
				answerMessage.ReplyMarkup = keyboard.NewBackKeyboard()
			} else {
				answerMessage.Text = model.WelcomeMessage
			}
		} else if message.Command() == model.HelpRawCommand {
			answerMessage.Text = model.HelpCommandMessage
		} else {
			answerMessage.Text = model.WrongCommandMessage
		}
	} else if message.Text == model.BackCommand {
		handleBackCommand(message)
		answerMessage.Text = model.BackCommandMessage
		userState.Delete(strconv.Itoa(message.From.ID))
	} else if message.Text == model.InboxCommand {
		return handleInboxCommand(message, nil, message.From, message.Chat)
	} else if message.Text == model.LinkCommand {
		link, err := handleMyLinkCommand(message)
		if err != nil {
			answerMessage.Text = model.SomeErrorOccured
		} else {
			answerMessage.Text = link
		}
	} else if message.Text == model.HelpCommand {
		answerMessage.Text = model.HelpCommandMessage
	} else if message.Text == model.NewMessageCommand {
		handleNewMessage(message)
		answerMessage.Text = model.NewMessageCommandMessageInputMessage
		answerMessage.ReplyMarkup = keyboard.NewBackKeyboard()
	} else {
		state, found := userState.Get(strconv.Itoa(message.From.ID))
		if found {
			var messageState model.NewMessageState
			var secretMessage *model.SecretMessage
			if (state.(model.UserState)).Command == model.NewMessageCommand {
				messageState, secretMessage = handleNewMessageArgs(message, state.(model.UserState))
			} else if (state.(model.UserState)).Command == model.ReplyCommand {
				messageState, secretMessage = handleReplyCommandArgs(message, state.(model.UserState))
			} else if (state.(model.UserState)).Command == model.NewMessageByLinkCommand {
				messageState, secretMessage = handleNewMessageByLinkArgs(message, state.(model.UserState))
			}

			if messageState == model.NewMessageStateInputUsername {
				answerMessage.Text = model.NewMessageCommandUsernameMessage
				answerMessage.ReplyMarkup = keyboard.NewBackKeyboard()
			} else if messageState == model.NewMessageStateInputText {
				answerMessage.Text = model.NewMessageCommandMessageInputMessage
				answerMessage.ReplyMarkup = keyboard.NewBackKeyboard()
			} else if messageState == model.NewMessageStateSent {
				answerMessage.Text = model.NewMessageCommandSentMessage
				if secretMessage != nil && secretMessage.ReceiverID != "" {
					id, _ := strconv.ParseInt(secretMessage.ReceiverID, 10, 64)
					receiverNotifyMessage = &botAPI.MessageConfig{}
					*receiverNotifyMessage = botAPI.NewMessage(id, model.NewMessageReceived)
				}
			} else if messageState == model.NewMessageStateError {
				answerMessage.Text = model.NewMessageCommandSentMessage
			}

			userState.Delete(strconv.Itoa(message.From.ID))
		}
	}

	chattables := []botAPI.Chattable{
		answerMessage,
	}
	if receiverNotifyMessage != nil {
		chattables = append(chattables, receiverNotifyMessage)
	}
	return chattables
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
			HashID:       util.GetMD5(strconv.Itoa(message.From.ID))[:model.HashIDLength],
			LastUseEpoch: time.Now().Unix(),
		}

		db.PeopleCollection.Insert(person)
	} else {
		person.LastUseEpoch = time.Now().Unix()
		//db.PeopleCollection.Update(&model.Person{Username: person.Username}, person)
	}
}
