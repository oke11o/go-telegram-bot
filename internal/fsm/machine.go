package fsm

import (
	"context"
	"github.com/oke11o/go-telegram-bot/internal/config"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/oke11o/go-telegram-bot/internal/model"
	"github.com/oke11o/go-telegram-bot/internal/model/iface"
)

func NewDeps(cfg config.Config, repo iface.Repo, sender iface.Sender, logger *slog.Logger) *Deps {
	return &Deps{
		Cfg:    cfg,
		Repo:   repo,
		Sender: sender,
		Logger: logger,
	}

}

type Deps struct {
	Repo   iface.Repo
	Sender iface.Sender
	Logger *slog.Logger
	Cfg    config.Config
}

type State struct {
	User    model.User
	Session model.Session
	Update  tgbotapi.Update
}

type Machine interface {
	Switch(ctx context.Context, state State) (context.Context, Machine, State, error)
}
