package tournament

import (
	"context"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"log/slog"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const CreateTournamentCommand = "/create"

func NewCreateTournament(deps *fsm.Deps) *CreateTournament {
	return &CreateTournament{
		Base: base.Base{Deps: deps},
	}
}

type CreateTournament struct {
	base.Base
}

func (m *CreateTournament) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}
	if !state.User.IsManager && !state.User.IsMaintainer {
		smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "You dont have enough permissions for this action.", 0)
		return ctx, smc, state, nil
	}

	ses := model.NewCreateTournamentSession(state.User.ID)

	split := strings.Split(state.Update.Message.Text, " ")
	var text string
	if len(split) > 1 {
		text = strings.Join(split[1:], " ")
	}
	if text == "" {
		ses.SetStatus(model.SessionCreateTournamentSetTitle)
		_, err := m.Deps.Repo.SaveSession(ctx, ses)
		if err != nil {
			m.Deps.Logger.ErrorContext(ctx, "Cant save session", slog.String("error", err.Error()))
			smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session for user %s", state.User.Username))
			return ctx, smc, state, nil
		}

		smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "Please text title of the tournament", 0)
		return ctx, smc, state, nil
	}
	ses.SetArg("title", text)
	ses.SetStatus(model.SessionCreateTournamentSetDate)
	_, err := m.Deps.Repo.SaveSession(ctx, ses)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant save session", slog.String("error", err.Error()))
		smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session for user %s", state.User.Username))
		return ctx, smc, state, nil
	}

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, "Please text start date of the tournament", 0)
	return ctx, smc, state, nil
}
