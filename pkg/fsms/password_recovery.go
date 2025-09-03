package fsms

import (
	"context"
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth/utils"
	"github.com/Araks1255/mangacage_bot/pkg/helpers"
	"github.com/Araks1255/mangacage_bot/pkg/keyboards"
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

const GET_NEW_PASSWORD_STATE = "get-new-password"
const VERIFY_PASSWORD_CHANGING = "verify-password-changing"

type GetNewPasswordHandler struct{ DB *gorm.DB }

type VerifyPasswordChangingHandler struct {
	DB  *gorm.DB
	bot *tgbotapi.BotAPI
}

func NewGetNewPasswordHandler(db *gorm.DB) GetNewPasswordHandler {
	return GetNewPasswordHandler{DB: db}
}
func NewVerifyPasswordChangingHandler(db *gorm.DB, bot *tgbotapi.BotAPI) VerifyPasswordChangingHandler {
	return VerifyPasswordChangingHandler{DB: db, bot: bot}
}

func (h GetNewPasswordHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	return fsm.MessageConfig{
		MessageConfig: tgbotapi.MessageConfig{
			Text:     "Введите новый пароль",
			BaseChat: tgbotapi.BaseChat{ReplyMarkup: keyboards.Cancel},
		},
	}
}

func (h GetNewPasswordHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "Отмена" {
		return helpers.MainStateTextTransition("Ладно"), Data{}
	}

	var userName *string

	if err := h.DB.Raw("SELECT user_name FROM users WHERE tg_user_id = ?", update.Message.Chat.ID).Scan(&userName).Error; err != nil {
		log.Println(err)
		return helpers.MainStateTextTransition("Произошла непредвиденная ошибка"), Data{}
	}

	if userName == nil {
		return helpers.MainStateTextTransition("Вы еще не вошли в аккаунт"), Data{}
	}

	data.user.UserName = *userName
	data.user.NewPassword = update.Message.Text

	data.messagesWithPasswordIDs = append(data.messagesWithPasswordIDs, update.Message.MessageID)

	return fsm.StateTransition(VERIFY_PASSWORD_CHANGING), data
}

func (h VerifyPasswordChangingHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	response := fmt.Sprintf("%s, вы уверены что хотите сменить пароль?", data.user.UserName)

	return fsm.MessageConfig{
		MessageConfig: tgbotapi.MessageConfig{
			Text:     response,
			BaseChat: tgbotapi.BaseChat{ReplyMarkup: keyboards.YesOrNot},
		},
	}
}

func (h VerifyPasswordChangingHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	defer func() {
		for i := 0; i < len(data.messagesWithPasswordIDs); i++ {
			h.bot.Request(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, data.messagesWithPasswordIDs[i]))
		}
		data.messagesWithPasswordIDs = []int{}
	}()
	if update.Message.Text == "Отмена" || update.Message.Text == "Нет" {
		return helpers.MainStateTextTransition("Ладно"), Data{}
	}

	if update.Message.Text == "Да" {
		hash, err := utils.GenerateHashPassword(data.user.NewPassword)

		if err != nil {
			log.Println(err)
			return helpers.MainStateTextTransition("Не удалось сгенерировать хэш для нового пароля"), Data{}
		}

		result := h.DB.Exec("UPDATE users SET password = ? WHERE user_name = ?", hash, data.user.UserName)

		if result.Error != nil {
			log.Println(err)
			return helpers.MainStateTextTransition("Произошла ошибка при изменении пароля"), Data{}
		}

		if result.RowsAffected == 0 {
			return helpers.MainStateTextTransition("Не удалось обновить пароль"), Data{}
		}

		return helpers.MainStateTextTransition("Пароль успешно изменен"), Data{}
	}

	return fsm.TextTransition("Отвечайте"), data
}
