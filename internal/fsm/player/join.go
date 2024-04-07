package player

import (
	"context"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/base"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const JoinCommand = "/join"

type tournamentMappingStr struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}
type tournamentMapping map[int64]tournamentMappingStr

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

	return m.Base.DefaultSwitch(ctx, state, "")
	//tours, tourMapping, err := getTournaments(ctx, m.Deps.Repo)
	//if err != nil {
	//	m.Deps.Logger.ErrorContext(ctx, "cant get opened tournaments", slog.String("error", err.Error()))
	//	smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("cant get tournaments for user %s", state.User.Username))
	//	return ctx, smc, state, nil
	//}
	//
	//state.Session.SetArg("tourMapping", tourMapping)
	//state.Session, err = m.Deps.Repo.SaveSession(ctx, state.Session)
	//if err != nil {
	//	m.Deps.Logger.ErrorContext(ctx, "Cant save session", slog.String("error", err.Error()))
	//	smc := m.CombineSenderMachines(state, "Something wrong. Try again latter", fmt.Sprintf("Cant save session for user %s", state.User.Username))
	//	return ctx, smc, state, nil //fmt.Errorf("repo.SaveSession() error: %w", err)
	//}
	//
	//toursTexts := make([]string, 0, len(tours))
	//for i, tour := range tours {
	//	toursTexts = append(toursTexts, fmt.Sprintf("%d. %s [%s]", i+1, tour.Title, tour.Date))
	//}
	//text := "For which tournament you want to join?\n" + strings.Join(toursTexts, "\n")
	//
	//smc := sender.NewSenderMachine(m.Deps, state.Update.Message.Chat.ID, text, 0)
	//
	//return ctx, smc, state, nil
}
