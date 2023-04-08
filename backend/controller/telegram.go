package controller

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"primeprice.com/dal"
)

// SendTgNotification sends notification for telegram user
func SendTgNotification(message string) error {
	// check app configuration for scraping
	cfg, err := dal.GetAppConfig()
	if err != nil {
		return err
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return err
	}

	bot.Debug = false

	chatID, err := strconv.ParseInt(cfg.ChatID, 10, 64)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(int64(chatID), "")
	msg.Text = message
	msg.ParseMode = tgbotapi.ModeHTML

	_, err = bot.Send(msg)
	return err
}
