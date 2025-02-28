package fsms

import (
	"context"

	"github.com/Araks1255/mangacage_bot/pkg/common/utils"
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const GET_USER_NAME_STATE = "get-user-name"
const GET_PASSWORD_STATE = "get-password"

type GetUserNameHandler struct{}
type GetPasswordHandler struct{}

func (h GetUserNameHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	return fsm.TextMessageConfig("Введите ваше имя пользователя")
}

func (h GetUserNameHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "" {
		return fsm.TextTransition("Пожалуйста, введите ваше имя пользователя"), data
	}

	data.user.UserName = update.Message.Text
	data.user.TgUserID = update.Message.Chat.ID

	return fsm.StateTransition(GET_PASSWORD_STATE), data
}

func (h GetPasswordHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	return fsm.TextMessageConfig("Введите ваш пароль")
}

func (h GetPasswordHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "" {
		return fsm.TextTransition("Пожалуйста, введите ваш пароль"), data
	}

	data.user.Password = update.Message.Text

	var existingUserPassword string
	_db.Raw("SELECT password FROM users WHERE user_name = ?", data.user.UserName).Scan(&existingUserPassword)
	if existingUserPassword == "" {
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Пользователь не найден"),
			},
		}, data
	}

	if ok := utils.IsPasswordCorrect(data.user.Password, existingUserPassword); !ok {
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный пароль"),
			},
		}, data
	}

	if result := _db.Exec("UPDATE users SET tg_user_id = ? WHERE user_name = ?", data.user.TgUserID, data.user.UserName); result.Error != nil {
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось войти в аккаунт"),
			},
		}, data
	}

	return fsm.Transition{
		State: fsm.UndefinedState,
		MessageConfig: fsm.MessageConfig{
			MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Вход в аккаунт выполнен успешно"),
		},
	}, data
}
