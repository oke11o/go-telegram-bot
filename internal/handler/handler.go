package handler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/oke11o/go-telegram-bot/internal/config"
	"github.com/oke11o/go-telegram-bot/internal/fsm"
	"github.com/oke11o/go-telegram-bot/internal/fsm/router"
	"github.com/oke11o/go-telegram-bot/internal/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/oke11o/go-telegram-bot/internal/model"
	"github.com/oke11o/go-telegram-bot/internal/model/iface"
	"github.com/oke11o/go-telegram-bot/pgk/utils/str"
)

type incomer interface {
	Income(ctx context.Context, requestID string, update tgbotapi.Update) (model.User, error)
}

func New(cfg config.Config, l *slog.Logger, income incomer, repo iface.Repo) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: l,
		income: income,
		repo:   repo,
	}
}

type Handler struct {
	logger *slog.Logger
	sender iface.Sender
	income incomer
	repo   iface.Repo
	cfg    config.Config
}

func (h *Handler) SetSender(sender iface.Sender) {
	h.sender = sender
}

func (h *Handler) HandleUpdate(ctx context.Context, update tgbotapi.Update) error {
	requestID := fmt.Sprintf("%s-%d", str.RandStringRunes(32, ""), update.UpdateID)
	ctx = log.AppendCtx(ctx, slog.String("request_id", requestID))

	user, err := h.income.Income(ctx, requestID, update)
	if err != nil {
		h.logger.ErrorContext(ctx, "income.Income Error", err)
		return fmt.Errorf("income.Income() err: %w", err)
	}
	ctx = log.AppendCtx(ctx, slog.Int64("user_id", user.ID))

	deps := fsm.NewDeps(h.cfg, h.repo, h.sender, h.logger)
	routr, err := router.NewRouter(deps)
	if err != nil {
		h.logger.ErrorContext(ctx, "fsm.NewRouter Error", err)
		return fmt.Errorf("fsm.NewRouter() err: %w", err)
	}
	machine, state, err := routr.GetMachine(ctx, user, update)
	if err != nil {
		h.logger.ErrorContext(ctx, "router.GetMachine Error", err)
		return fmt.Errorf("router.GetMachine() err: %w", err)
	}

	for machine != nil {
		ctx, machine, state, err = machine.Switch(ctx, state)
		if err != nil {
			h.logger.ErrorContext(ctx, "machine.Switch Error", err)
			return fmt.Errorf("machine.Switch() err: %w", err)
		}
	}

	return nil
}
