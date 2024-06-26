package player

import (
	"context"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const JoinCommand = "/join"

func NewJoin(deps *fsm.Deps) *Join {
	return &Join{
		Base: Base{Base: base.Base{Deps: deps}},
	}
}

type Join struct {
	Base
}

func (m *Join) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	state.Session = model.NewJoinSession(state.User.ID)

	return m.Base.JoinSwitch(ctx, state, "")
}
