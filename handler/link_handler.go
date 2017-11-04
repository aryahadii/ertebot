package handler

import (
	"fmt"
	"strconv"

	"gitlab.com/arha/Ertebot/model"
	"gitlab.com/arha/Ertebot/util"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

func handleMyLinkCommand(message *botAPI.Message) (string, error) {
	hashID, err := util.GetHashID(strconv.Itoa(message.From.ID))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(model.LinkTemplate, hashID), nil
}
