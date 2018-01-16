package updater

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/arha/Ertebot/configuration"
	"gitlab.com/arha/Ertebot/handler"
	botAPI "gopkg.in/telegram-bot-api.v4"
)

const (
	botToken = "419987007:AAFkHJPaj2bwW-mvY8g219ZIGFAXV8jzcws"
)

var (
	bot *botAPI.BotAPI
)

func InitBot() {
	var err error
	bot, err = botAPI.NewBotAPI(configuration.ErtebotConfig.GetString("bot-token"))
	if err != nil {
		log.WithError(err).Fatalln("Can't initialize bot")
	}
	bot.Debug = configuration.ErtebotConfig.GetBool("debug")

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
		go func(update botAPI.Update) {
			var msgs []botAPI.Chattable
			if update.Message != nil {
				msgs = handler.HandleMessage(update.Message)
			} else {
				msgs = handler.HandleCallback(update.CallbackQuery)
			}

			for _, msg := range msgs {
				bot.Send(msg)
			}
		}(update)
	}
}
