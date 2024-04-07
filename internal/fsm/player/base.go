package player

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"log/slog"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/model"
	"github.com/oke11o/go-telegram-bot/internal/model/iface"
)

type Base struct {
	base.Base
}

func (m *Base) DefaultSwitch(ctx context.Context, state fsm.State, errorMessage string) (context.Context, fsm.Machine, fsm.State, error) {
	tours, tourMapping, err := getTournaments(ctx, m.Deps.Repo)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "cant get opened tournaments", slog.String("error", err.Error()))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("cant get tournaments for user %s", state.User.Username))
		return ctx, smc, state, nil
	}

	state.Session.SetArg("tourMapping", tourMapping)
	state.Session, err = m.Deps.Repo.SaveSession(ctx, state.Session)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant save session", slog.String("error", err.Error()))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session for user %s", state.User.Username))
		return ctx, smc, state, nil //fmt.Errorf("repo.SaveSession() error: %w", err)
	}

	toursTexts := make([]string, 0, len(tours))
	for i, tour := range tours {
		toursTexts = append(toursTexts, fmt.Sprintf("%d. %s [%s]", i+1, tour.Title, tour.Date))
	}
	text := errorMessage + "For which tournament you want to join?\n" + strings.Join(toursTexts, "\n")

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, text, 0)

	return ctx, smc, state, nil
}

func getTournaments(ctx context.Context, repo iface.Repo) ([]model.Tournament, string, error) {
	tours, err := repo.GetOpenedTournaments(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("repo.GetOpenedTournaments() error: %w", err)
	}
	tourMapping := make(tournamentMapping)
	for i, tour := range tours {
		tourMapping[int64(i+1)] = tournamentMappingStr{
			ID:    tour.ID,
			Title: tour.Title,
			Date:  tour.Date,
		}
	}
	b, err := json.Marshal(tourMapping)
	if err != nil {
		return nil, "", fmt.Errorf("json.Marshal() error: %w", err)
	}
	return tours, string(b), err
}
