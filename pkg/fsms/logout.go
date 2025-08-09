package fsms

import (
	"context"
	"log"

	"github.com/Araks1255/mangacage_bot/pkg/helpers"
	"github.com/Araks1255/mangacage_bot/pkg/keyboards"
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

const LOGOUT_VERIFY = "logout-verify"

type LogoutVerifyHandler struct{ DB *gorm.DB }

func NewLogoutVerifyHandler(db *gorm.DB) LogoutVerifyHandler { return LogoutVerifyHandler{DB: db} }

func (h LogoutVerifyHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	return fsm.MessageConfig{
		MessageConfig: tgbotapi.MessageConfig{
			Text:     "Вы уверены, что хотите отвязать Telegram аккаунт от mangacage аккаунта?",
			BaseChat: tgbotapi.BaseChat{ReplyMarkup: keyboards.YesOrNot},
		},
	}
}

func (h LogoutVerifyHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "Нет" {
		return helpers.MainStateTextTransition("Ладно"), Data{}
	}

	if update.Message.Text == "Да" {
		result := h.DB.Exec("UPDATE users SET tg_user_id = NULL WHERE tg_user_id = ?", update.Message.Chat.ID)

		if result.Error != nil {
			log.Println(result.Error)
			return helpers.MainStateTextTransition("Произошла ошибка при отвязывании аккаунта"), Data{}
		}

		if result.RowsAffected == 0 {
			return helpers.MainStateTextTransition("Ваш Telegram аккаунт итак не привязан к mangacage аккаунту"), Data{}
		}

		return helpers.MainStateTextTransition("Ваш Telegram аккаунт успешно отвязан от mangacage аккаунта"), Data{}
	}

	return fsm.TextTransition("Отвечайте"), Data{}
}
