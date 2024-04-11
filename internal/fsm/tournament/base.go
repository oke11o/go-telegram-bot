package tournament

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

type Base struct {
	base.Base
}

func (m *Base) StartSwitch(ctx context.Context, state fsm.State, errorMessage string) (context.Context, fsm.Machine, fsm.State, error) {
	return m.defaultSwitch(ctx, state, errorMessage, "start")
}

func (m *Base) defaultSwitch(ctx context.Context, state fsm.State, errorMessage string, verb string) (context.Context, fsm.Machine, fsm.State, error) {
	var err error
	var tours []model.Tournament
	if state.User.IsMaintainer {
		tours, err = m.Deps.Repo.TournamentOpenedAll(ctx)
	} else {
		tours, err = m.Deps.Repo.TournamentOpenedByManager(ctx, state.User.ID)
	}
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "cant get OpenedTournaments", slog.String("error", err.Error()))
		combineMachine := m.CombineSenderMachines(state, "Something wrong. Try again latter", "Cant get tournament list")
		return ctx, combineMachine, state, nil
	}

	tourMappingStr, err := json.Marshal(tours)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "cant marshal tourMapping", slog.String("error", err.Error()))
		combineMachine := m.CombineSenderMachines(state, "Something wrong. Try again latter", "Cant get tournament list")
		return ctx, combineMachine, state, nil
	}

	state.Session.SetArg("tourMapping", string(tourMappingStr))
	state.Session, err = m.Deps.Repo.SaveSession(ctx, state.Session)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant save session", slog.String("error", err.Error()))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session for user %s", state.User.Username))
		return ctx, smc, state, nil
	}
	toursTexts := make([]string, 0, len(tours))
	for i, tour := range tours {
		toursTexts = append(toursTexts, fmt.Sprintf("%d. %s [%s]", i+1, tour.Title, tour.Date))
	}
	text := fmt.Sprintf("Which tournament you want to start?\n%s", strings.Join(toursTexts, "\n"))

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, text, 0)
	return ctx, smc, state, nil
}
