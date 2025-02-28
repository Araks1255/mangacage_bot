package fsms

import (
	"context"
	"fmt"
	"log"

	"github.com/Araks1255/mangacage_bot/pkg/common/utils"
	"github.com/Araks1255/mangacage_bot/pkg/keyboards"
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const GET_NEW_PASSWORD_STATE = "get-new-password"
const VERIFY_PASSWORD_CHANGING = "verify-password-changing"

type GetNewPasswordHandler struct{}
type VerifyPasswordChangingHandler struct{}

func (h GetNewPasswordHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	msg := tgbotapi.NewMessage(data.user.TgUserID, "Введите новый пароль")
	msg.ReplyMarkup = keyboards.Canel
	return fsm.MessageConfig{
		MessageConfig: msg,
	}
}

func (h GetNewPasswordHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "Отмена" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ладно")
		msg.ReplyMarkup = keyboards.Main
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: msg,
			},
		}, data
	}

	var userName string
	_db.Raw("SELECT user_name FROM users WHERE tg_user_id = ?", update.Message.Chat.ID).Scan(&userName)
	if userName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы ещё не вошли в аккаунт")
		msg.ReplyMarkup = keyboards.Main
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: msg,
			},
		}, data
	}

	data.user.TgUserID = update.Message.Chat.ID

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

	msg := tgbotapi.NewMessage(data.user.TgUserID, reponse)
	msg.ReplyMarkup = keyboards.YesOrNot
	return fsm.MessageConfig{
		MessageConfig: msg,
	}
}

func (h VerifyPasswordChangingHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "Да" {
		hash, err := utils.GenerateHashPassword(data.user.NewPassword)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось сгенерировать хэш для нового пароля")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			return fsm.Transition{
				State: fsm.UndefinedState,
				MessageConfig: fsm.MessageConfig{
					MessageConfig: msg,
				},
			}, data
		}

		if result := _db.Exec("UPDATE users SET password = ? WHERE user_name = ?", hash, data.user.UserName); result.RowsAffected == 0 {
			log.Println(result.Error)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось обновить пароль")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			return fsm.Transition{
				State: fsm.UndefinedState,
				MessageConfig: fsm.MessageConfig{
					MessageConfig: msg,
				},
			}, data
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пароль успешно обновлён")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: msg,
			},
		}, data
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ладно")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	if update.Message.Text == "Нет" {
		return fsm.Transition{
			State: fsm.UndefinedState,
			MessageConfig: fsm.MessageConfig{
				MessageConfig: msg,
			},
		}, data
	}

	return fsm.TextTransition("Отвечай."), data
}
