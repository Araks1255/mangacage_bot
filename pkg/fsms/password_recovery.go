package fsms

import (
	"context"
	"fmt"
	"log"

	"github.com/Araks1255/mangacage_bot/pkg/common/utils"
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const GET_NEW_PASSWORD_STATE = "get-new-password"
const VERIFY_PASSWORD_CHANGING = "verify-password-changing"

type GetNewPasswordHandler struct{}
type VerifyPasswordChangingHandler struct{}

func (h GetNewPasswordHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	return fsm.TextMessageConfig("Введите новый пароль")
}

func (h GetNewPasswordHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	var userName string
	_db.Raw("SELECT user_name FROM users WHERE tg_user_id = ?", update.Message.Chat.ID).Scan(&userName)
	if userName == "" {
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ещё не вошли в аккаунт"),
			},
		}, data
	}

	data.user.UserName = userName
	data.user.NewPassword = update.Message.Text

	return fsm.StateTransition(VERIFY_PASSWORD_CHANGING), data
}

func (h VerifyPasswordChangingHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	reponse := fmt.Sprintf(
		"%s, вы уверены что хотите сменить пароль на %s?",
		data.user.UserName,
		data.user.NewPassword,
	)
	return fsm.TextMessageConfig(reponse)
}

func (h VerifyPasswordChangingHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "Да" {
		hash, err := utils.GenerateHashPassword(data.user.NewPassword)
		if err != nil {
			return fsm.Transition{
				State: fsm.UndefinedState,
				MessageConfig: fsm.MessageConfig{
					MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать хэш для нового пароля"),
				},
			}, data
		}

		if result := _db.Exec("UPDATE users SET password = ? WHERE user_name = ?", hash, data.user.UserName); result.RowsAffected == 0 {
			log.Println(result.Error)
			return fsm.Transition{
				State: fsm.UndefinedState,
				MessageConfig: fsm.MessageConfig{
					MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обновить пароль"),
				},
			}, data
		}

		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Пароль успешно обновлён"),
			},
		}, data
	}

	if update.Message.Text == "Нет" {
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: tgbotapi.NewMessage(update.Message.Chat.ID, "Ладно"),
			},
		}, data
	}

	return fsm.TextTransition("Отвечай."), data
}
