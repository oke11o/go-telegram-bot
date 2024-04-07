package tournament

import (
	"context"
	"fmt"
	"strings"

	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/sender"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const ListTournamentCommand = "/list"

func NewListTournament(deps *fsm.Deps) *ListTournament {
	return &ListTournament{
		deps: deps,
	}
}

type ListTournament struct {
	deps *fsm.Deps
}

func (m *ListTournament) Switch(ctx context.Context, state fsm.State) (context.Context, fsm.Machine, fsm.State, error) {
	if state.Update.Message == nil {
		return ctx, nil, state, fmt.Errorf("unexpected part. ")
	}

	tournaments, err := m.deps.Repo.GetOpenedTournaments(ctx)
	if err != nil {
		combineMachine := fsm.NewCombine(nil,
			sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, "Something wrong. Try again latter", 0),
			sender.NewSenderMachine(m.deps, m.deps.Cfg.MaintainerChatID, fmt.Sprintf("Cant get tournament list"), 0),
		)
		return ctx, combineMachine, state, nil
	}
	text := m.PrintTournaments(tournaments)

	smc := sender.NewSenderMachine(m.deps, state.Update.Message.Chat.ID, text, 0)
	return ctx, smc, state, nil
}

func (m *ListTournament) PrintTournaments(tournaments []model.Tournament) string {
	tours := make([]string, len(tournaments))
	for i := range tournaments {
		tours[i] = fmt.Sprintf("- %s [%s]", tournaments[i].Title, tournaments[i].Date)
	}
	return "List of opened tournaments:\n" + strings.Join(tours, "\n")
}
