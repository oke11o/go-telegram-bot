package router

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/help"
	"github.com/oke11o/go-telegram-bot/internal/fsm/maintainer"
	"github.com/oke11o/go-telegram-bot/internal/fsm/player"
	"github.com/oke11o/go-telegram-bot/internal/fsm/session"
	"github.com/oke11o/go-telegram-bot/internal/fsm/tournament"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const StartCommand = "/start"

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
	ses, err := r.deps.Repo.GetOpenedSession(ctx, user.ID)
	if err != nil {
		return nil, state, fmt.Errorf("repo.GetOpenedSession(%d) err: %w", user.ID, err)
	}
	state.Session = ses

	// Получить текущую сессию из базы, если есть
	if update.Message != nil {
		cmdMachine := r.resolveCommandMachine(update)
		if cmdMachine != nil {
			return cmdMachine, state, nil
		}

		return session.NewSessionMachine(r.deps), state, nil
	}

	return nil, state, fmt.Errorf("unknown state machine")
}

func (r *Router) resolveCommandMachine(update tgbotapi.Update) fsm.Machine {
	if strings.HasPrefix(update.Message.Text, maintainer.AddAdminCommand) {
		return maintainer.NewAddAdmin(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, maintainer.RemoveAdminCommand) {
		return maintainer.NewRemoveAdmin(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, tournament.StartTournamentCommand) {
		return tournament.NewStartTournament(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, tournament.CreateTournamentCommand) {
		return tournament.NewCreateTournament(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, tournament.ListTournamentCommand) {
		return tournament.NewListTournament(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, player.JoinCommand) {
		return player.NewJoin(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, player.LeaveCommand) {
		return player.NewLeave(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, player.MembersCommand) {
		return player.NewMembers(r.deps)
	}
	if strings.HasPrefix(update.Message.Text, help.HelpCommand) || strings.HasPrefix(update.Message.Text, StartCommand) {
		return help.NewHelp(r.deps)
	}
	return nil
}
