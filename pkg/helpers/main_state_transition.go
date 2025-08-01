package helpers

import (
	"github.com/Araks1255/mangacage_bot/pkg/keyboards"
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainStateTextTransition(text string) fsm.Transition {
	return fsm.Transition{
		State: fsm.UndefinedState,
		MessageConfig: fsm.MessageConfig{
			MessageConfig: tgbotapi.MessageConfig{
				Text:     text,
				BaseChat: tgbotapi.BaseChat{ReplyMarkup: keyboards.Main},
			},
		},
	}
}
