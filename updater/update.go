package updater

import (
	log "github.com/Sirupsen/logrus"
	"gitlab.com/arha/Ertebot/handler"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

const (
	botToken = "bot-token"
)

var (
	bot *botAPI.BotAPI
)

func init() {
	var err error
	bot, err = botAPI.NewBotAPI(botToken)
	if err != nil {
		log.WithError(err).Fatalln("Can't initialize bot")
	}
	bot.Debug = true

	log.Infof("Authorized on account %s", bot.Self.UserName)
}

func Update() {
	updateConfig := botAPI.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.WithError(err).Warnln("Updater isn't working")
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := handler.Handle(update.Message)
		bot.Send(msg)
	}
}
