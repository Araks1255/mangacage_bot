package fsms

import (
	"context"
	"log"

	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type User struct {
	UserName    string
	Password    string
	NewPassword string
}

type Data struct {
	user User
}

func RegisterFSMs(bot *tgbotapi.BotAPI, db *gorm.DB) {
	commands := make(map[string]fsm.TransitionProvider[Data])

	commands["start"] = StartCommandHandler{}
	commands["cancel"] = CancelHandler{}

	configs := make(map[fsm.State]fsm.StateHandler[Data])

	configs[fsm.UndefinedState] = MainStateHandler{}
	configs[MAIN_STATE] = MainStateHandler{}

	configs[GET_USER_NAME_STATE] = NewGetUserNameHandler(db)
	configs[GET_PASSWORD_STATE] = NewGetPasswordHandler(db)

	configs[GET_NEW_PASSWORD_STATE] = NewGetNewPasswordHandler(db)
	configs[VERIFY_PASSWORD_CHANGING] = NewVerifyPasswordChangingHandler(db)

	botFSM := fsm.NewBotFsm(bot, configs, fsm.WithCommands[Data](commands))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	ctx := context.Background()

	for update := range updates {
		err := botFSM.HandleUpdate(ctx, &update)
		if err != nil {
			log.Println(err)
		}
	}
}
