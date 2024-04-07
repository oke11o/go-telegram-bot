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

func (m *CreateTournament) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if !state.User.IsManager && !state.User.IsMaintainer {
		smc := sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}

	ses := model.NewCreateTournamentSession(state.User.ID)

	text := strings.TrimPrefix(state.Update.Message.Text, CreateTournamentCommand)
	text = strings.TrimSpace(text)
	if text == "" {
		ses.SetStatus(model.SessionCreateTournamentAskTitle)
		_, err := m.deps.Repo.SaveSession(ctx, ses)
		if err != nil {
			combineMachine := fsm.NewCombine(nil,
				sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "Something wrong. Try again latter", 0),
				sender.NewSenderMachine(m.deps, m.deps.Cfg.MaintainerChatID, fmt.Sprintf("Cant save session fog user %s", state.User.Username), 0),
			)
			return ctx, combineMachine, state, nil
		}

		smc := sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "Please text title of the tournament", 0)
		return ctx, smc, state, nil
	}
	ses.SetArg("title", text)
	ses.SetStatus(model.SessionCreateTournamentAskDate)
	_, err := m.deps.Repo.SaveSession(ctx, ses)
	if err != nil {
		combineMachine := fsm.NewCombine(nil,
			sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "Something wrong. Try again latter", 0),
			sender.NewSenderMachine(m.deps, m.deps.Cfg.MaintainerChatID, fmt.Sprintf("Cant save session fog user %s", state.User.Username), 0),
		)
		return ctx, combineMachine, state, nil
	}

	smc := sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "Please text start date of the tournament", 0)
	return ctx, smc, state, nil
}
