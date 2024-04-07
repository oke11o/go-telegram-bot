package player

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/logger"
	"github.com/oke11o/go-telegram-bot/internal/model"
	"github.com/oke11o/go-telegram-bot/internal/model/iface"
)

type tournamentMappingStr struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}
type tournamentMapping map[int64]tournamentMappingStr

type Base struct {
	base.Base
}

func (m *Base) JoinSwitch(ctx context.Context, state fsm.State, errorMessage string) (context.Context, fsm.Machine, fsm.State, error) {
	return m.defaultSwitch(ctx, state, errorMessage, "join", tournamentForJoin)
}

func (m *Base) LeaveSwitch(ctx context.Context, state fsm.State, errorMessage string) (context.Context, fsm.Machine, fsm.State, error) {
	return m.defaultSwitch(ctx, state, errorMessage, "leave", tournamentForLeave)
}

func (m *Base) defaultSwitch(ctx context.Context, state fsm.State, errorMessage string, verb string, getter tournamentGetter) (context.Context, fsm.Machine, fsm.State, error) {
	ctx = logger.AppendCtx(ctx, slog.String("verb", verb), slog.String("errorMessage", errorMessage))
	tours, tourMapping, err := getter(ctx, m.Deps.Repo, state.User.ID)
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
		return ctx, smc, state, nil
	}

	toursTexts := make([]string, 0, len(tours))
	for i, tour := range tours {
		toursTexts = append(toursTexts, fmt.Sprintf("%d. %s [%s]", i+1, tour.Title, tour.Date))
	}
	text := fmt.Sprintf("%sWhich tournament you want to %s?\n%s", errorMessage, verb, strings.Join(toursTexts, "\n"))

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, text, 0)

	return ctx, smc, state, nil
}

type tournamentGetter func(ctx context.Context, repo iface.Repo, userID int64) ([]model.Tournament, string, error)

func tournamentForJoin(ctx context.Context, repo iface.Repo, _ int64) ([]model.Tournament, string, error) {
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

func tournamentForLeave(ctx context.Context, repo iface.Repo, userID int64) ([]model.Tournament, string, error) {
	tours, err := repo.GetMemberTournaments(ctx, userID)
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
