package model

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	ID           int64  `db:"id"`
	Username     string `db:"username"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	LanguageCode string `db:"language_code"`
	IsBot        bool   `db:"is_bot"`
	IsMaintainer bool   `db:"is_maintainer"`
	IsManager    bool   `db:"is_manager"`
}

func NewUserFromTgUpdate(update tgbotapi.Update) (User, error) {
	if update.Message == nil {
		return User{}, fmt.Errorf("unknown to parse User (empty Message) from Update: %+v", update)
	}
	if update.Message.From == nil {
		return User{}, fmt.Errorf("unknown to parse User (empty Message.From) from Update: %+v", update)
	}
	return User{
		ID:           update.Message.From.ID,
		Username:     update.Message.From.UserName,
		FirstName:    update.Message.From.FirstName,
		LastName:     update.Message.From.LastName,
		LanguageCode: update.Message.From.LanguageCode,
	}, nil
}
