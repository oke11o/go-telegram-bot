package base

import (
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
)

type Base struct {
	Deps *fsm.Deps
}

func (m *Base) CombineSenderMachines(state fsm.State, userAnswer string, maintainerAnswer string) *fsm.Combine {
	combineMachine := fsm.NewCombine(nil,
		sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, userAnswer, 0),
		sender.NewSenderMachine(m.Deps, m.Deps.Cfg.MaintainerChatID, maintainerAnswer, 0),
	)
	return combineMachine
}
