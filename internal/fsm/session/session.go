package session

import (
	"context"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
)

func NewSessionMachine(deps *fsm.Deps) *SessionMachine {
	return &SessionMachine{
		deps: deps,
	}
}

type SessionMachine struct {
	deps *fsm.Deps
}

func (s *SessionMachine) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	// Получить сессию пользователя
	// Сохранить сессию в state
	// Под длинному switch определить нужную машину
	// Если машину определить не удалость - вернуть ответ - выберете действие - или меню

	return ctx, nil, state, nil
}
