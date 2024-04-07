package player

import (
	"context"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const LeaveCommand = "/leave"

func NewLeave(deps *fsm.Deps) *Leave {
	return &Leave{
		Base: Base{Base: base.Base{Deps: deps}},
	}
}

type Leave struct {
	Base
}

func (m *Leave) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	state.Session = model.NewLeaveSession(state.User.ID)

	return m.Base.LeaveSwitch(ctx, state, "")
}
