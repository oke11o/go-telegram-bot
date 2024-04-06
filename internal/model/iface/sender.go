package iface

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
