package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/go-telegram-bot/internal/model"
	"testing"
)

func (s *Suite) TestHandler_Player() {
	// arrange
	mainAdmin := model.User{ID: 111, Username: "main_admin", FirstName: "Main", LastName: "Admin", LanguageCode: "en", IsMaintainer: true}
	ctx := context.Background()
	_, err := s.repo.SaveUser(ctx, mainAdmin)
	s.Require().NoError(err)
	q := `insert into tournament (id,title,date,status,created_by,created_at,updated_at) 
values (101,'tour1', '2024-03-21', 'created', 111, '2024-03-21 00:00:00', '2024-03-21 00:00:00'),
       (102,'tour2', '2024-03-22', 'created', 111, '2024-03-22 00:00:00', '2024-03-22 00:00:00'),
       (103,'tour3', '2024-03-23', 'in_progress', 111, '2024-03-23 00:00:00', '2024-03-23 00:00:00'),
       (104,'tour4', '2024-03-24', 'finished', 111, '2024-03-24 00:00:00', '2024-03-24 00:00:00')`
	_, err = s.dbx.Exec(q)
	s.Require().NoError(err)

	s.T().Run("/join", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(`For which tournament you want to join?
1. tour1 [2024-03-21]
2. tour2 [2024-03-22]
3. tour3 [2024-03-23]`, msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "/join"}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		err = s.dbx.Select(&incomes, `select * from income_request`)
		s.Require().NoError(err)
		s.Len(incomes, 1)

		var sessions []model.Session
		err = s.dbx.Select(&sessions, `select * from session where user_id=20 order by id asc`)
		s.Require().NoError(err)
		s.Require().Len(sessions, 1)
		s.Require().Equal(model.SessionJoinProcess, sessions[0].Status)
		s.Require().Equal(int64(20), sessions[0].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":101,\"title\":\"tour1\",\"date\":\"2024-03-21\"},\"2\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"3\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[0].Data)
		s.Require().False(sessions[0].Closed)
	})

	s.T().Run("wrong text", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(`Sorry, I got wrong answer.
For which tournament you want to join?
1. tour1 [2024-03-21]
2. tour2 [2024-03-22]
3. tour3 [2024-03-23]`, msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "whaaaaat???"}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		err = s.dbx.Select(&incomes, `select * from income_request`)
		s.Require().NoError(err)
		s.Len(incomes, 2)

		var sessions []model.Session
		err = s.dbx.Select(&sessions, `select * from session where user_id=20 order by id asc`)
		s.Require().NoError(err)
		s.Require().Len(sessions, 1)
		s.Require().Equal(model.SessionJoinProcess, sessions[0].Status)
		s.Require().Equal(int64(20), sessions[0].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":101,\"title\":\"tour1\",\"date\":\"2024-03-21\"},\"2\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"3\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[0].Data)
		s.Require().False(sessions[0].Closed)
	})

}
