package fsms

import (
	"context"
	"log"

	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"

	"github.com/Araks1255/mangacage/pkg/auth/utils"
	"github.com/Araks1255/mangacage_bot/pkg/helpers"
	"github.com/Araks1255/mangacage_bot/pkg/keyboards"
	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

const GET_USER_NAME_STATE = "get-user-name"
const GET_PASSWORD_STATE = "get-password"

type GetUserNameHandler struct{ DB *gorm.DB }
type GetPasswordHandler struct{ DB *gorm.DB }

func NewGetUserNameHandler(db *gorm.DB) GetUserNameHandler { return GetUserNameHandler{DB: db} }
func NewGetPasswordHandler(db *gorm.DB) GetPasswordHandler { return GetPasswordHandler{DB: db} }

func (h GetUserNameHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	return fsm.MessageConfig{
		MessageConfig: tgbotapi.MessageConfig{
			Text:     "Введите ваше имя пользователя",
			BaseChat: tgbotapi.BaseChat{ReplyMarkup: keyboards.Cancel},
		},
	}
}

func (h GetUserNameHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "Отмена" { // Как я понял, хэндлеров высшего порядка или middleware в этой библиотеке нет, а commands нужны только для комманд (которые всегда начинаются со слэша). так что только так (свои middleware я писать желанием не горю для двух функций)
		return helpers.MainStateTextTransition("Ладно"), Data{}
	}

	data.user.UserName = update.Message.Text
	return fsm.StateTransition(GET_PASSWORD_STATE), data
}

func (h GetPasswordHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	return fsm.TextMessageConfig("Введите ваш пароль")
}

func (h GetPasswordHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "Отмена" {
		return helpers.MainStateTextTransition("Ладно"), Data{}
	}

	data.user.Password = update.Message.Text

	var existingUserPassword *string

	if err := h.DB.Raw("SELECT password FROM users WHERE user_name = ?", data.user.UserName).Scan(&existingUserPassword).Error; err != nil {
		log.Println(err)
		return helpers.MainStateTextTransition("Произошла непредвиденная ошибка"), Data{}
	}

	if existingUserPassword == nil {
		return helpers.MainStateTextTransition("Пользователь не найден"), Data{}
	}

	if !utils.CompareHashPassword(data.user.Password, *existingUserPassword) {
		return helpers.MainStateTextTransition("Неверный пароль"), Data{}
	}

	if err := h.DB.Exec("UPDATE users SET tg_user_id = ? WHERE user_name = ?", update.Message.Chat.ID, data.user.UserName).Error; err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniUsersTgUserID) {
			return helpers.MainStateTextTransition("Этот telegram аккаунт уже привязан к mangacage аккаунту"), Data{}
		}
		return helpers.MainStateTextTransition("Не удалось привязать аккаунт"), Data{}
	}

	return helpers.MainStateTextTransition("Аккаунт успешно привязан"), Data{}
}
