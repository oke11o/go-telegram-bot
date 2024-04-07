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

const ListTournamentCommand = "/list"

func NewListTournament(deps *fsm.Deps) *ListTournament {
	return &ListTournament{
		Base: base.Base{Deps: deps},
	}
}

type ListTournament struct {
	base.Base
}

func (m *ListTournament) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}

	tournaments, err := m.Deps.Repo.GetOpenedTournaments(ctx)
	if err != nil {
		m.Deps.Logger.ErrorContext(ctx, "Cant GetOpenedTournaments", slog.String("error", err.Error()))
		combineMachine := m.CombineSenderMachines(state, "Something wrong. Try again latter", "Cant get tournament list")
		return ctx, combineMachine, state, nil
	}
	text := m.PrintTournaments(tournaments)

	smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, text, 0)
	return ctx, smc, state, nil
}

func (m *ListTournament) PrintTournaments(tournaments []model.Tournament) string {
	tours := make([]string, len(tournaments))
	for i := range tournaments {
		tours[i] = fmt.Sprintf("- %s [%s]", tournaments[i].Title, tournaments[i].Date)
	}
	return "List of opened tournaments:\n" + strings.Join(tours, "\n")
}
