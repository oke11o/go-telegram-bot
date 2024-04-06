package router

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/maintainer"
	"github.com/oke11o/go-telegram-bot/internal/fsm/session"
	"github.com/oke11o/go-telegram-bot/internal/fsm/tournament"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

type Router struct {
	deps *fsm.Deps
}

func NewRouter(deps *fsm.Deps) (*Router, error) {
	return &Router{
		deps: deps,
	}, nil
}

func (r *Router) GetMachine(ctx context.Context, user model.User, update tgbotapi.Update) (fsm.Machine, fsm.State, error) {
	state := fsm.State{
		User:   user,
		Update: update,
	}
	// Получить текущую сессию из базы, если есть
	if update.Message != nil {
		if strings.HasPrefix(update.Message.Text, maintainer.AddAdminCommand) {
			return maintainer.NewAddAdmin(r.deps), state, nil
		}
		if strings.HasPrefix(update.Message.Text, maintainer.RemoveAdminCommand) {
			return maintainer.NewRemoveAdmin(r.deps), state, nil
		}
		if strings.HasPrefix(update.Message.Text, tournament.CreateTournamentCommand) {
			return tournament.NewCreateTournament(r.deps), state, nil
		}

		return session.NewSessionMachine(r.deps), state, nil
	}

	return nil, state, fmt.Errorf("unknown state machine")
}
