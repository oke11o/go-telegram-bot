package tournament

import (
	"context"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

func NewCreateTournamenSetDate(deps *fsm.Deps) *CreateTournamentAskDate {
	return &CreateTournamentAskDate{
		deps: deps,
	}
}

type CreateTournamentAskDate struct {
	deps *fsm.Deps
}

func (m *CreateTournamentAskDate) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if !state.User.IsManager && !state.User.IsMaintainer {
		smc := sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}

	date := state.Update.Message.Text
	err := m.validateDate(date)
	if err != nil {
		smc := sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "Invalid tournament title. Text date again", 0)
		return ctx, smc, state, nil
	}
	state.Session.SetArg("date", date)

	ses, err := m.deps.Repo.SaveSession(ctx, state.Session)
	if err != nil {
		smc := m.combineMachine(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session %d", state.Session.ID))
		return ctx, smc, state, nil
	}
	state.Session = ses

	title, ok := state.Session.GetArg("title")
	if !ok {
		smc := m.combineMachine(state, "Something wrong. Try again latter", fmt.Sprintf("Unexpected behaviour for session %d", state.Session.ID))
		return ctx, smc, state, nil
	}
	tournament := model.NewTournament(title, date, state.User.ID)
	tournament, err = m.deps.Repo.SaveTournament(ctx, tournament)
	if err != nil {
		smc := m.combineMachine(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save tournament for session %d", state.Session.ID))
		return ctx, smc, state, nil
	}

	err = m.deps.Repo.CloseSession(ctx, state.Session)
	if err != nil {
		smc := m.combineMachine(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session %d", state.Session.ID))
		return ctx, smc, state, nil
	}
	state.Session.Closed = true

	smc := sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "The tournament was successfully created", 0)
	return ctx, smc, state, nil
}

func (m *CreateTournamentAskDate) combineMachine(state fsm.State, userAnswer string, maintainerAnswer string) *fsm.Combine {
	combineMachine := fsm.NewCombine(nil,
		sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, userAnswer, 0),
		sender.NewSenderMachine(m.deps, m.deps.Cfg.MaintainerChatID, maintainerAnswer, 0),
	)
	return combineMachine
}

func (m *CreateTournamentAskDate) validateDate(text string) error {
	if text == "" {
		return fmt.Errorf("empty title")
	}
	return nil
}
