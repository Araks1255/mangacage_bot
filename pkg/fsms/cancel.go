package fsms

import (
	"context"

	"github.com/Araks1255/mangacage_bot/pkg/helpers"

	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CancelHandler struct{}

func (h CancelHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	return helpers.MainStateTextTransition("Ладно"), Data{}
}
