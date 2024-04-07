package app

import (
	"context"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/model/iface"
	"github.com/oke11o/go-telegram-bot/internal/service"
	"log/slog"

	"github.com/oke11o/go-telegram-bot/internal/app/bot"
	"github.com/oke11o/go-telegram-bot/internal/config"
	handler2 "github.com/oke11o/go-telegram-bot/internal/handler"
	"github.com/oke11o/go-telegram-bot/internal/logger"
	"github.com/oke11o/go-telegram-bot/internal/repository/mongo"
	"github.com/oke11o/go-telegram-bot/internal/repository/pg"
	"github.com/oke11o/go-telegram-bot/internal/repository/sqlite"
)

func Run(ctx context.Context) error {
	l := logger.New(true, slog.LevelDebug)
	err := config.InitDotEnv()
	if err != nil {
		panic(err)
	}

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	var repo iface.Repo
	switch cfg.DBType {
	case sqlite.DBType:
		repo, err = sqlite.New(cfg.Sqlite)
		if err != nil {
			return fmt.Errorf("sqlite.New() err: %w", err)
		}
	case mongo.DBType:
		repo = mongo.New()
	case pg.DBType:
		repo = pg.New()
	default:
		return fmt.Errorf("unknown db_type %s", cfg.DBType)
	}
	income := service.NewIncomeServce(repo)
	b := bot.NewBot(cfg, l, handler2.New(cfg, l, income, repo))
	return b.Run(ctx)
}
