package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Main = tgbotapi.NewOneTimeReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Войти в аккаунт"),
	), tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Сменить пароль"),
	),
)
