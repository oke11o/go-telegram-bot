package player

import (
	"context"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const MembersCommand = "/members"

func NewMembers(deps *fsm.Deps) *Members {
	return &Members{
		Base: Base{Base: base.Base{Deps: deps}},
	}
}

type Members struct {
	Base
}

func (m *Members) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	state.Session = model.NewMembersSession(state.User.ID)

	return m.Base.MembersSwitch(ctx, state, "")
}
