package main

import (
	"github.com/Araks1255/mangacage_bot/pkg/common/db"
	"github.com/Araks1255/mangacage_bot/pkg/fsms"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	token := viper.Get("TOKEN").(string)
	dbUrl := viper.Get("DB_URL").(string)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	_db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	fsms.RegisterFSMs(bot, _db)
}
