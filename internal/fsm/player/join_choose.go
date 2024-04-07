package player

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
)

func NewJoinChoose(deps *fsm.Deps) *JoinChoose {
	return &JoinChoose{
		Base: base.Base{Deps: deps},
	}
}

type JoinChoose struct {
	base.Base
}

func (m *JoinChoose) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	tourMapping, ok := state.Session.GetArg("tourMapping")
	if !ok {
		// TODO: FirstStep Join again

	}
	var mapping tournamentMapping
	err := json.Unmarshal([]byte(tourMapping), &mapping)
	if err != nil {
		// TODO: FirstStep Join again

	}
	choose, err := strconv.ParseInt(strings.TrimSpace(state.Update.Message.Text), 10, 64)
	if err != nil {
		//TODO: FirstStep Join again
	}
	if _, ok := mapping[choose]; !ok {
		//TODO: FirstStep Join again
	}
	state.Session.SetArg("choose", strconv.FormatInt(choose, 10))
	state.Session, err = m.Deps.Repo.SaveSession(ctx, state.Session)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant save session", slog.String("error", err.Error()))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session for user %s", state.User.Username))
		return ctx, smc, state, nil
	}

	err = m.Deps.Repo.AddPlayerToTournament(ctx, state.User.ID, mapping[choose].ID)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant add player to tournament", slog.String("error", err.Error()))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant add player to tournament for user @%s", state.User.Username))
		return ctx, smc, state, nil
	}

	err = m.Deps.Repo.CloseSession(ctx, state.Session)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant close session", slog.String("error", err.Error()))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant close session for user @%s", state.User.Username))
		return ctx, smc, state, nil
	}

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, fmt.Sprintf("You are successfully joined to the tournament `%s - %s`", mapping[choose].Title, mapping[choose].Date), 0)

	return ctx, smc, state, nil
}
