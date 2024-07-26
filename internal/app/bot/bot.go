package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/oke11o/wslog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/oke11o/go-telegram-bot/internal/config"
	"github.com/oke11o/go-telegram-bot/internal/model/iface"
)

func NewBot(cfg config.Config, logger *slog.Logger, handler handler) *Bot {
	return &Bot{
		cfg:     cfg,
		logger:  logger,
		handler: handler,
	}
}

type handler interface {
	SetSender(bot iface.Sender)
	HandleUpdate(ctx context.Context, update tgbotapi.Update) error
}

type Bot struct {
	cfg     config.Config
	logger  *slog.Logger
	bot     *tgbotapi.BotAPI
	handler handler
}

func (b *Bot) Run(ctx context.Context) error {
	var err error
	b.bot, err = tgbotapi.NewBotAPI(b.cfg.TgToken)
	if err != nil {
		return fmt.Errorf("tgbotapi.NewBotAPI err: %w", err)
	}

	b.bot.Debug = false

	b.handler.SetSender(b.bot)

	b.logger.InfoContext(ctx, "Authorized on account", slog.String("username", b.bot.Self.UserName))

	b.bot.Send(tgbotapi.NewMessage(b.cfg.MaintainerChatID, "бот стартанул"))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	//TODO:
	/**
	for {
		select {
			case update := <- updates
			case <- ctx.Done
		}
	}
	*/
	for update := range updates {
		upCtx := wslog.AppendCtx(ctx, slog.Int("updateId", update.UpdateID))
		b.logger.DebugContext(upCtx, "got tg update")
		err = b.handler.HandleUpdate(upCtx, update)
		if err != nil {
			return fmt.Errorf("handler.HandleUpdate() err: %w", err)
		}

	}
	return nil
}
