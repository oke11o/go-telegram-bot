package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/oke11o/go-telegram-bot/internal/app/bot"
	"github.com/oke11o/go-telegram-bot/internal/config"
	"github.com/oke11o/go-telegram-bot/internal/handler"
	"github.com/oke11o/go-telegram-bot/internal/log"
	"github.com/oke11o/go-telegram-bot/internal/model/iface"
	"github.com/oke11o/go-telegram-bot/internal/repository/mongo"
	"github.com/oke11o/go-telegram-bot/internal/repository/pg"
	"github.com/oke11o/go-telegram-bot/internal/repository/sqlite"
	"github.com/oke11o/go-telegram-bot/internal/service"
)

func Run(ctx context.Context, version string) error {
	l := log.New(true, slog.LevelDebug)
	ctx = log.AppendCtx(ctx, slog.String("version", version))
	err := config.InitDotEnv()
	if err != nil {
		l.ErrorContext(ctx, "error config.InitDotEnv()", slog.String("error", err.Error()))
		return fmt.Errorf("config.InitDotEnv() err: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		l.ErrorContext(ctx, "error config.Load()", slog.String("error", err.Error()))
		return fmt.Errorf("config.Load() err: %w", err)
	}
	var repo iface.Repo
	switch cfg.DBType {
	case sqlite.DBType:
		repo, err = sqlite.New(cfg.Sqlite)
		if err != nil {
			l.ErrorContext(ctx, "error sqlite.New()", slog.String("error", err.Error()))
			return fmt.Errorf("sqlite.New() err: %w", err)
		}
	case mongo.DBType:
		repo = mongo.New()
	case pg.DBType:
		repo = pg.New()
	default:
		l.ErrorContext(ctx, "error sqlite.New()", slog.String("error", err.Error()))
		return fmt.Errorf("unknown db_type %s", cfg.DBType)
	}
	income := service.NewIncomeServce(repo)
	b := bot.NewBot(cfg, l, handler.New(cfg, l, income, repo))
	err = b.Run(ctx)
	if err != nil {
		l.ErrorContext(ctx, "error bot.Run()", slog.String("error", err.Error()))
		return fmt.Errorf("bot.Run() err: %w", err)
	}

	return nil
}
