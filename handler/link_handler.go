package handler

import (
	"fmt"
	"strconv"

	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/ui/keyboard"
	"gitlab.com/arha/Ertebot/util"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleMyLinkCommand(message *botAPI.Message) (string, interface{}) {
	hashID, err := util.GetHashID(strconv.Itoa(message.From.ID))
	if err != nil {
		return model.SomeErrorOccured, keyboard.NewMainKeyboard()
	}

	return fmt.Sprintf(model.LinkTemplate, hashID), keyboard.NewMainKeyboard()
}
