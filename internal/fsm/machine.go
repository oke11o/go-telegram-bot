package fsm

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/oke11o/go-telegram-bot/internal/model"
	"github.com/oke11o/go-telegram-bot/internal/model/iface"
)

func NewDeps(repo iface.Repo, sender iface.Sender, logger *slog.Logger) *Deps {
	return &Deps{
		Repo:   repo,
		Sender: sender,
		Logger: logger,
	}

}

type Deps struct {
	Repo   iface.Repo
	Sender iface.Sender
	Logger *slog.Logger
}

type State struct {
	User   model.User
	Update tgbotapi.Update
}

type Machine interface {
	Switch(ctx context.Context, state State) (context.Context, Machine, State, error)
}
