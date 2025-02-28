package fsms

import (
	"context"
	"log"

	"github.com/Araks1255/mangacage_bot/pkg/keyboards"

	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

const MAIN_STATE = "main"

var _db *gorm.DB

type User struct {
	UserName    string
	Password    string
	TgUserID    int64
	NewPassword string
}

type Data struct {
	user User
}

type StartCommandHandler struct{}

func (h StartCommandHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	data.user.TgUserID = update.Message.Chat.ID
	return fsm.StateTransition(MAIN_STATE), data
}

type MainStateHandler struct{}

func (h MainStateHandler) MessageFn(ctx context.Context, data Data) fsm.MessageConfig {
	msg := tgbotapi.NewMessage(data.user.TgUserID, "Здравствуйте, вас приветствует официальный бот mangacage! Все функции представлены на клавиатуре внизу\n(чтобы получать уведомления о выходе новых глав и иметь возможность сменить пароль войдите в аккаунт)")
	msg.ReplyMarkup = keyboards.Main
	return fsm.MessageConfig{
		MessageConfig: msg,
	}
}

func (h MainStateHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	if update.Message.Text == "Войти в аккаунт" {
		return fsm.StateTransition(GET_USER_NAME_STATE), data
	}

	if update.Message.Text == "Сменить пароль" {
		return fsm.StateTransition(GET_NEW_PASSWORD_STATE), data
	}

	return fsm.TextTransition("Выберите вариант с клавиатуры"), data
}

func RegisterFSMs(bot *tgbotapi.BotAPI, db *gorm.DB) {
	_db = db

	commands := make(map[string]fsm.TransitionProvider[Data])
	commands["start"] = StartCommandHandler{}
	commands["Войти в аккаунт"] = GetUserNameHandler{}
	commands["Сменить пароль"] = GetNewPasswordHandler{}

	configs := make(map[fsm.State]fsm.StateHandler[Data])
	configs[fsm.UndefinedState] = MainStateHandler{}
	configs[MAIN_STATE] = MainStateHandler{}

	configs[GET_USER_NAME_STATE] = GetUserNameHandler{}
	configs[GET_PASSWORD_STATE] = GetPasswordHandler{}

	configs[GET_NEW_PASSWORD_STATE] = GetNewPasswordHandler{}
	configs[VERIFY_PASSWORD_CHANGING] = VerifyPasswordChangingHandler{}

	botFSM := fsm.NewBotFsm(bot, configs, fsm.WithCommands[Data](commands))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	ctx := context.TODO()

	for update := range updates {
		err := botFSM.HandleUpdate(ctx, &update)
		if err != nil {
			log.Println(err)
		}
	}
}
