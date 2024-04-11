package tournament

import (
	"context"
	"fmt"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const FinishTournamentCommand = "/finish"

func NewFinishTournament(deps *fsm.Deps) *FinishTournament {
	return &FinishTournament{
		Base: Base{base.Base{Deps: deps}},
	}
}

type FinishTournament struct {
	Base
}

func (m *FinishTournament) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if !state.User.IsManager && !state.User.IsMaintainer {
		smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}
	state.Session = model.NewFinishTournamentSession(state.User.ID)

	return m.Base.StartSwitch(ctx, state, "")
}
