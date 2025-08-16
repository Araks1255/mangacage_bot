package main

import (
	"net/http"

	"github.com/Araks1255/mangacage_bot/pkg/common/db"
	"github.com/Araks1255/mangacage_bot/pkg/common/http/clients"
	"github.com/Araks1255/mangacage_bot/pkg/fsms"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	token := viper.Get("TOKEN").(string)
	dbUrl := viper.Get("DB_URL").(string)

	rateLimiter := clients.NewSendMessageRateLimitedRoundTripper(http.DefaultTransport, 25)
	defer rateLimiter.Stop()

	httpClient := &http.Client{Transport: rateLimiter}

	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, httpClient)
	if err != nil {
		panic(err)
	}

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	fsms.RegisterFSMs(bot, db)
}
