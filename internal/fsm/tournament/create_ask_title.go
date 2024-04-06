package tournament

import (
	"context"
	"fmt"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
)

func NewCreateTournamentAskTitle(deps *fsm.Deps) *CreateTournamentAskTitle {
	return &CreateTournamentAskTitle{
		deps: deps,
	}
}

type CreateTournamentAskTitle struct {
	deps *fsm.Deps
}

func (s *CreateTournamentAskTitle) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if !state.User.IsManager {
		smc := sender.NewSenderMachine(s.deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}

	// Создать сессию. Сохранить в сессию состояние, что мы находимся в процессе создания турнира и ждем название

	text := strings.TrimPrefix(state.Update.Message.Text, CreateTournamentCommand)
	text = strings.TrimSpace(text)
	if text == "" {
		// Попросить пользователя ввести название турнира

		return ctx, nil, state, nil
	}
	// Сохранить название турнира в сессию
	// Попросить пользователя ввести дату начала турнира

	return ctx, nil, state, nil
}
