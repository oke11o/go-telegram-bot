package sender

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
)

func NewSenderMachine(deps *fsm.Deps, toChatID int64, text string, replyMsgID int) *SenderMachine {
	return &SenderMachine{
		deps:       deps,
		toChatID:   toChatID,
		text:       text,
		replyMsgID: replyMsgID,
	}
}

type SenderMachine struct {
	deps       *fsm.Deps
	toChatID   int64
	text       string
	replyMsgID int
}

func (s *SenderMachine) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	msg := tgbotapi.NewMessage(s.toChatID, s.text)
	if s.replyMsgID != 0 {
		msg.ReplyToMessageID = s.replyMsgID
	}

	respMsg, err := s.deps.Sender.Send(msg)
	if err != nil {
		s.deps.Logger.ErrorContext(ctx, "sender.Send", slog.String("error", err.Error()))
	} else {
		s.deps.Logger.DebugContext(ctx, "sender.Send", slog.Any("response", respMsg))
	}

	return ctx, nil, state, nil
}
