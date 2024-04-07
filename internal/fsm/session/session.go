package session

import (
	"context"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/fsm/tournament"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

func NewSessionMachine(deps *fsm.Deps) *SessionMachine {
	return &SessionMachine{
		deps: deps,
	}
}

type SessionMachine struct {
	deps *fsm.Deps
}

func (s *SessionMachine) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	var scm fsm.Machine
	switch state.Session.Status {
	case model.SessionCreateTournamentProcess:
		scm = tournament.NewCreateTournament(s.deps)
	case model.SessionCreateTournamentSetTitle:
		scm = tournament.NewCreateTournamentSetTitle(s.deps)
	case model.SessionCreateTournamentSetDate:
		scm = tournament.NewCreateTournamenSetDate(s.deps)
	default:
		scm = sender.NewSenderMachine(s.deps, state.User.ID, "Choose action", 0)
	}

	return ctx, scm, state, nil
}
