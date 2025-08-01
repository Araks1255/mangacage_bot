package helpers

import (
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StateMessageConfigTransition(state fsm.State, messageConfig tgbotapi.MessageConfig) fsm.Transition {
	return fsm.Transition{
		State: state,
		MessageConfig: fsm.MessageConfig{
			MessageConfig: messageConfig,
		},
	}
}
