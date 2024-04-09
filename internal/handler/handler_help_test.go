package handler

import (
	"context"
	"github.com/oke11o/go-telegram-bot/internal/fsm/help"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Suite) TestHandler_Help() {
	s.T().Run("/start", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(help.InstructionText, msg.Text)
			},
		})
		err := h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "/start"}, // Notice!!!
		})
		s.Require().NoError(err)
	})

	s.T().Run("/help", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(help.InstructionText, msg.Text)
			},
		})
		err := h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "/help"}, // Notice!!!
		})
		s.Require().NoError(err)
	})
}
