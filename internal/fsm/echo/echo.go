package echo

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
)

func New(deps *fsm.Deps) *Echo {
	return &Echo{
		deps: deps,
	}
}

type Echo struct {
	deps *fsm.Deps
}

func (s *Echo) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message != nil {
		s.deps.Logger.DebugContext(ctx, "update.Message", slog.String("username", state.Update.Message.From.UserName), slog.String("text", state.Update.Message.Text))

		msg := tgbotapi.NewMessage(state.Update.Message.Chat.ID, state.Update.Message.Text)
		msg.ReplyToMessageID = state.Update.Message.MessageID

		respMsg, err := s.deps.Sender.Send(msg)
		if err != nil {
			s.deps.Logger.ErrorContext(ctx, "sender.Send", slog.String("error", err.Error()))
		} else {
			s.deps.Logger.DebugContext(ctx, "sender.Send", slog.Any("response", respMsg))
		}

	}
	return ctx, nil, state, nil
}
