package tournament

import (
	"context"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/model"
	"log/slog"
)

func NewCreateTournamentSetTitle(deps *fsm.Deps) *CreateTournamentAskTitle {
	return &CreateTournamentAskTitle{
		Base: base.Base{Deps: deps},
	}
}

type CreateTournamentAskTitle struct {
	base.Base
}

func (m *CreateTournamentAskTitle) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if !state.User.IsManager && !state.User.IsMaintainer {
		smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}

	text := state.Update.Message.Text
	err := m.validate(text)
	if err != nil {
		smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "Invalid tournament title", 0)
		return ctx, smc, state, nil
	}
	state.Session.SetArg("title", text)
	state.Session.SetStatus(model.SessionCreateTournamentSetDate)
	ses, err := m.Deps.Repo.SaveSession(ctx, state.Session)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant save session", slog.String("error", err.Error()))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session for user %s", state.User.Username))
		return ctx, smc, state, nil //fmt.Errorf("repo.SaveSession() error: %w", err)
	}
	state.Session = ses
	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "Please text start date of the tournament", 0)
	return ctx, smc, state, nil
}

func (m *CreateTournamentAskTitle) validate(text string) error {
	if text == "" {
		return fmt.Errorf("empty title")
	}
	return nil
}
