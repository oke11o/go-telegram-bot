package model

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IncomeRequest struct {
	ID               int64  `db:"id"`
	FromID           int64  `db:"from_id"`
	MessageID        int64  `db:"message_id"`
	ReplyToMessageID int64  `db:"reply_to_message_id"`
	RequestID        string `db:"request_id"`
	Message          []byte `db:"message"`
	Username         string `db:"username"`
	Text             string `db:"text"`
}

func NewIncomeRequestFromTgUpdate(requestID string, update tgbotapi.Update) (IncomeRequest, error) {
	message, err := json.Marshal(update)
	if err != nil {
		return IncomeRequest{}, fmt.Errorf("json.Marshal(update) err: %w", err)
	}
	income := IncomeRequest{
		RequestID: requestID,
		Message:   message,
	}
	if update.Message != nil {
		income.Text = update.Message.Text
		if update.Message.From != nil {
			income.FromID = update.Message.From.ID
			income.Username = update.Message.From.UserName
		}
		income.MessageID = int64(update.Message.MessageID)
		if update.Message.ReplyToMessage != nil {
			income.ReplyToMessageID = int64(update.Message.ReplyToMessage.MessageID)
		}
	}

	return income, nil
}
