package fsms

import (
	"context"

	"github.com/Araks1255/mangacage_bot/pkg/helpers"
	"github.com/Araks1255/mangacage_bot/pkg/keyboards"

	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const MAIN_STATE = "main"

type MainStateHandler struct{}

func (h MainStateHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	return helpers.MessageConfigWithKeyboard(
		"Здравствуйте, вас приветствует официальный бот mangacage! Все функции представлены на клавиатуре внизу\n(чтобы получать уведомления и иметь возможность сменить пароль войдите в аккаунт)",
		keyboards.Main,
	)
}

func (h MainStateHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	switch update.Message.Text {
	case "Войти в аккаунт":
		return fsm.StateTransition(GET_USER_NAME_STATE), Data{}
	case "Сменить пароль":
		return fsm.StateTransition(GET_NEW_PASSWORD_STATE), Data{}
	}
	return fsm.TextTransition("Выберите вариант с клавиатуры"), Data{}
}
