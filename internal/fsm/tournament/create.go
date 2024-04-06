package tournament

import (
	"context"
	"fmt"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const CreateTournamentCommand = "/create"

func NewCreateTournament(deps *fsm.Deps) *CreateTournament {
	return &CreateTournament{
		deps: deps,
	}
}

type CreateTournament struct {
	deps *fsm.Deps
}

func (s *CreateTournament) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if !state.User.IsManager {
		smc := sender.NewSenderMachine(s.deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}

	ses := model.NewCreateTournamentSession(state.User.ID)

	text := strings.TrimPrefix(state.Update.Message.Text, CreateTournamentCommand)
	text = strings.TrimSpace(text)
	if text == "" {
		ses.SetStatus(model.SessionCreateTournamentAskTitle)
		_, err := s.deps.Repo.SaveSession(ctx, ses)
		if err != nil {
			return ctx, nil, state, fmt.Errorf("Repo.SaveSession() error: %w", err)
		}

		smc := sender.NewSenderMachine(s.deps, state.Update.Message.Chat.ID, "Please type title of the tournament", 0)
		return ctx, smc, state, nil
	}
	ses.SetArg("title", text)
	ses.SetStatus(model.SessionCreateTournamentAskDate)
	_, err := s.deps.Repo.SaveSession(ctx, ses)
	if err != nil {
		return ctx, nil, state, fmt.Errorf("Repo.SaveSession() error: %w", err)
	}

	smc := sender.NewSenderMachine(s.deps, state.Update.Message.Chat.ID, "Please type start date of the tournament", 0)
	return ctx, smc, state, nil
}
