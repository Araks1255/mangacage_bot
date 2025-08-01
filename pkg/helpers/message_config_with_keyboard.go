package helpers

import (
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MessageConfigWithKeyboard(text string, keyboard tgbotapi.ReplyKeyboardMarkup) fsm.MessageConfig {
	return fsm.MessageConfig{
		MessageConfig: tgbotapi.MessageConfig{
			Text:     text,
			BaseChat: tgbotapi.BaseChat{ReplyMarkup: keyboard},
		},
	}
}
