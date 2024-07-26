package tournament

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/log"
	"github.com/oke11o/go-telegram-bot/internal/model"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
)

func NewFinishChooseTournament(deps *fsm.Deps) *FinishChooseTournament {
	return &FinishChooseTournament{
		Base: Base{base.Base{Deps: deps}},
	}
}

type FinishChooseTournament struct {
	Base
}

func (m *FinishChooseTournament) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if !state.User.IsManager && !state.User.IsMaintainer {
		smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}
	tourMapping, ok := state.Session.GetArg("tourMapping")
	if !ok {
		return m.Base.StartSwitch(ctx, state, "Something wrong\nPlease try again\n\n")
	}
	var tours []model.Tournament
	err := json.Unmarshal([]byte(tourMapping), &tours)
	if err != nil {
		return m.Base.StartSwitch(ctx, state, "Something wrong\nPlease try again\n\n")
	}
	choose, err := strconv.ParseInt(strings.TrimSpace(state.Update.Message.Text), 10, 64)
	if err != nil {
		return m.Base.StartSwitch(ctx, state, fmt.Sprintf("Invalid input `%s`\nChoose one of:\n\n", state.Update.Message.Text))
	}
	if choose < 0 || int(choose) >= len(tours) {
		return m.Base.StartSwitch(ctx, state, fmt.Sprintf("Invalid input `%s`\nChoose one of:\n\n", state.Update.Message.Text))
	}
	choose--
	state.Session.SetArg("choose", strconv.FormatInt(choose, 10))
	state.Session, err = m.Deps.Repo.SaveSession(ctx, state.Session)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant save session", log.Err(err))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session for user %s", state.User.Username))
		return ctx, smc, state, nil
	}
	err = m.Deps.Repo.TournamentStart(ctx, tours[choose].ID)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant add player to tournament", log.Err(err))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant add player to tournament for user @%s", state.User.Username))
		return ctx, smc, state, nil
	}

	err = m.Deps.Repo.CloseSession(ctx, state.Session)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant close session", log.Err(err))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant close session for user @%s", state.User.Username))
		return ctx, smc, state, nil
	}

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, fmt.Sprintf("You are successfully leave to the tournament `%s - %s`", tours[choose].Title, tours[choose].Date), 0)

	return ctx, smc, state, nil
}
