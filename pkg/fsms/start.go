package fsms

import (
	"context"

	fsm "github.com/Feolius/telegram-bot-fsm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartCommandHandler struct{}

func (h StartCommandHandler) TransitionFn(ctx context.Context, update *tgbotapi.Update, data Data) (fsm.Transition, Data) {
	return fsm.StateTransition(MAIN_STATE), data
}
