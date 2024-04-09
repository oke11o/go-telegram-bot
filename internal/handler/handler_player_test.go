package handler

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

func (s *Suite) TestHandler_PlayerJoinToTournament() {
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
				s.Require().Equal(`Which tournament you want to join?
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
				s.Require().Equal(`Invalid input `+"`whaaaaat???`"+`
Choose one of:

Which tournament you want to join?
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
		s.Require().Len(sessions, 2)
		s.Require().Equal(model.SessionJoinProcess, sessions[1].Status)
		s.Require().Equal(int64(20), sessions[1].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":101,\"title\":\"tour1\",\"date\":\"2024-03-21\"},\"2\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"3\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[1].Data)
		s.Require().False(sessions[1].Closed)
	})

	s.T().Run("wrong text (number)", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(`Invalid input `+"` 5 `"+`
Choose one of:

Which tournament you want to join?
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
				Text: " 5 "}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		err = s.dbx.Select(&incomes, `select * from income_request`)
		s.Require().NoError(err)
		s.Len(incomes, 3)

		var sessions []model.Session
		err = s.dbx.Select(&sessions, `select * from session where user_id=20 order by id asc`)
		s.Require().NoError(err)
		s.Require().Len(sessions, 3)
		s.Require().Equal(model.SessionJoinProcess, sessions[2].Status)
		s.Require().Equal(int64(20), sessions[2].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":101,\"title\":\"tour1\",\"date\":\"2024-03-21\"},\"2\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"3\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[2].Data)
		s.Require().False(sessions[2].Closed)
	})

	s.T().Run("successful join", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal("You are successfully joined to the tournament `tour2 - 2024-03-22`", msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "   2 "}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		err = s.dbx.Select(&incomes, `select * from income_request`)
		s.Require().NoError(err)
		s.Len(incomes, 4)

		var sessions []model.Session
		err = s.dbx.Select(&sessions, `select * from session where user_id=20 order by id asc`)
		s.Require().NoError(err)
		s.Require().Len(sessions, 4)
		s.Require().Equal(model.SessionJoinProcess, sessions[3].Status)
		s.Require().Equal(int64(20), sessions[3].UserID)
		s.Require().Equal(`{"choose":"2","tourMapping":"{\"1\":{\"id\":101,\"title\":\"tour1\",\"date\":\"2024-03-21\"},\"2\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"3\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[3].Data)
		s.Require().True(sessions[0].Closed)
		s.Require().True(sessions[1].Closed)
		s.Require().True(sessions[2].Closed)
		s.Require().True(sessions[3].Closed)

		var id int64
		err = s.dbx.Get(&id, `select id from member where user_id=20 and tournament_id=102`)
		s.Require().NoError(err)
		s.Require().NotEqual(int64(0), id)
	})
}

func (s *Suite) TestHandler_PlayerLeaveTournament() {
	ctx := context.Background()
	// arrange
	mainAdmin := model.User{ID: 111, Username: "main_admin", FirstName: "Main", LastName: "Admin", LanguageCode: "en", IsMaintainer: true}
	_, err := s.repo.SaveUser(ctx, mainAdmin)
	s.Require().NoError(err)
	tmpUser := model.User{ID: 20, Username: "tmp_user", FirstName: "Tmp", LastName: "User", LanguageCode: "en"}
	_, err = s.repo.SaveUser(ctx, tmpUser)
	s.Require().NoError(err)
	tmpUser2 := model.User{ID: 21, Username: "tmp_user2", FirstName: "Tmp2", LastName: "User2", LanguageCode: "en"}
	_, err = s.repo.SaveUser(ctx, tmpUser2)
	s.Require().NoError(err)
	q := `insert into tournament (id,title,date,status,created_by,created_at,updated_at) 
values (101,'tour1', '2024-03-21', 'created', 111, '2024-03-21 00:00:00', '2024-03-21 00:00:00'),
       (102,'tour2', '2024-03-22', 'created', 111, '2024-03-22 00:00:00', '2024-03-22 00:00:00'),
       (103,'tour3', '2024-03-23', 'in_progress', 111, '2024-03-23 00:00:00', '2024-03-23 00:00:00'),
       (104,'tour4', '2024-03-24', 'finished', 111, '2024-03-24 00:00:00', '2024-03-24 00:00:00')`
	_, err = s.dbx.Exec(q)
	s.Require().NoError(err)
	q = `insert into member (user_id,tournament_id,created_at)
values (20, 102, '2024-03-22 00:00:00'),
       (21, 102, '2024-03-22 00:00:00'),
       (111, 102, '2024-03-23 00:00:00'),
       (21, 103, '2024-03-22 00:00:00'),
       (20, 103, '2024-03-23 00:00:00')`
	_, err = s.dbx.Exec(q)
	s.Require().NoError(err)

	s.T().Run("/leave", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(`Which tournament you want to leave?
1. tour2 [2024-03-22]
2. tour3 [2024-03-23]`, msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "/leave"}, // Notice!!!
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
		s.Require().Equal(model.SessionLeaveProcess, sessions[0].Status)
		s.Require().Equal(int64(20), sessions[0].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"2\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[0].Data)
		s.Require().False(sessions[0].Closed)
	})

	s.T().Run("wrong text", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(`Invalid input `+"`whaaaaat???`"+`
Choose one of:

Which tournament you want to leave?
1. tour2 [2024-03-22]
2. tour3 [2024-03-23]`, msg.Text)
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
		s.Require().Len(sessions, 2)
		s.Require().Equal(model.SessionLeaveProcess, sessions[1].Status)
		s.Require().Equal(int64(20), sessions[1].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"2\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[1].Data)
		s.Require().False(sessions[1].Closed)
	})

	s.T().Run("wrong text (number)", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(`Invalid input `+"` 5 `"+`
Choose one of:

Which tournament you want to leave?
1. tour2 [2024-03-22]
2. tour3 [2024-03-23]`, msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: " 5 "}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		err = s.dbx.Select(&incomes, `select * from income_request`)
		s.Require().NoError(err)
		s.Len(incomes, 3)

		var sessions []model.Session
		err = s.dbx.Select(&sessions, `select * from session where user_id=20 order by id asc`)
		s.Require().NoError(err)
		s.Require().Len(sessions, 3)
		s.Require().Equal(model.SessionLeaveProcess, sessions[2].Status)
		s.Require().Equal(int64(20), sessions[2].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"2\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[2].Data)
		s.Require().False(sessions[2].Closed)
	})

	s.T().Run("successful leave", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal("You are successfully leave to the tournament `tour3 - 2024-03-23`", msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "   2 "}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		err = s.dbx.Select(&incomes, `select * from income_request`)
		s.Require().NoError(err)
		s.Len(incomes, 4)

		var sessions []model.Session
		err = s.dbx.Select(&sessions, `select * from session where user_id=20 order by id asc`)
		s.Require().NoError(err)
		s.Require().Len(sessions, 4)
		s.Require().Equal(model.SessionLeaveProcess, sessions[3].Status)
		s.Require().Equal(int64(20), sessions[3].UserID)
		s.Require().Equal(`{"choose":"2","tourMapping":"{\"1\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"2\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[3].Data)
		s.Require().True(sessions[0].Closed)
		s.Require().True(sessions[1].Closed)
		s.Require().True(sessions[2].Closed)
		s.Require().True(sessions[3].Closed)

		var id int64
		err = s.dbx.Get(&id, `select id from member where user_id=20 and tournament_id=103`)
		s.Require().True(errors.Is(err, sql.ErrNoRows))
	})
}

func (s *Suite) TestHandler_PlayerInTournament() {
	ctx := context.Background()
	// arrange
	mainAdmin := model.User{ID: 111, Username: "main_admin", FirstName: "Main", LastName: "Admin", LanguageCode: "en", IsMaintainer: true}
	_, err := s.repo.SaveUser(ctx, mainAdmin)
	s.Require().NoError(err)
	tmpUser := model.User{ID: 20, Username: "tmp_user", FirstName: "Tmp", LastName: "User", LanguageCode: "en"}
	_, err = s.repo.SaveUser(ctx, tmpUser)
	s.Require().NoError(err)
	tmpUser2 := model.User{ID: 21, Username: "tmp_user2", FirstName: "Tmp2", LastName: "User2", LanguageCode: "en"}
	_, err = s.repo.SaveUser(ctx, tmpUser2)
	s.Require().NoError(err)
	q := `insert into tournament (id,title,date,status,created_by,created_at,updated_at) 
values (101,'tour1', '2024-03-21', 'created', 111, '2024-03-21 00:00:00', '2024-03-21 00:00:00'),
       (102,'tour2', '2024-03-22', 'created', 111, '2024-03-22 00:00:00', '2024-03-22 00:00:00'),
       (103,'tour3', '2024-03-23', 'in_progress', 111, '2024-03-23 00:00:00', '2024-03-23 00:00:00'),
       (104,'tour4', '2024-03-24', 'finished', 111, '2024-03-24 00:00:00', '2024-03-24 00:00:00')`
	_, err = s.dbx.Exec(q)
	s.Require().NoError(err)
	q = `insert into member (user_id,tournament_id,created_at)
values (20, 102, '2024-03-22 00:00:00'),
       (21, 102, '2024-03-22 00:00:00'),
       (111, 102, '2024-03-23 00:00:00'),
       (21, 103, '2024-03-22 00:00:00'),
       (20, 103, '2024-03-23 00:00:00')`
	_, err = s.dbx.Exec(q)
	s.Require().NoError(err)

	s.T().Run("/members", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(`Which tournament you want to show members?
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
				Text: "/members"}, // Notice!!!
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
		s.Require().Equal(model.SessionMembersProcess, sessions[0].Status)
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
				s.Require().Equal(`Invalid input `+"`whaaaaat???`"+`
Choose one of:

Which tournament you want to show members?
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
		s.Require().Len(sessions, 2)
		s.Require().Equal(model.SessionMembersProcess, sessions[1].Status)
		s.Require().Equal(int64(20), sessions[1].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":101,\"title\":\"tour1\",\"date\":\"2024-03-21\"},\"2\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"3\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[1].Data)
		s.Require().False(sessions[1].Closed)
	})

	s.T().Run("wrong text (number)", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(`Invalid input `+"` 5 `"+`
Choose one of:

Which tournament you want to show members?
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
				Text: " 5 "}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		err = s.dbx.Select(&incomes, `select * from income_request`)
		s.Require().NoError(err)
		s.Len(incomes, 3)

		var sessions []model.Session
		err = s.dbx.Select(&sessions, `select * from session where user_id=20 order by id asc`)
		s.Require().NoError(err)
		s.Require().Len(sessions, 3)
		s.Require().Equal(model.SessionMembersProcess, sessions[2].Status)
		s.Require().Equal(int64(20), sessions[2].UserID)
		s.Require().Equal(`{"tourMapping":"{\"1\":{\"id\":101,\"title\":\"tour1\",\"date\":\"2024-03-21\"},\"2\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"3\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[2].Data)
		s.Require().False(sessions[2].Closed)
	})

	s.T().Run("successful members", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal("Players%\n0. tmp_user\n1. tmp_user2\n2. main_admin", msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "   2 "}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		err = s.dbx.Select(&incomes, `select * from income_request`)
		s.Require().NoError(err)
		s.Len(incomes, 4)

		var sessions []model.Session
		err = s.dbx.Select(&sessions, `select * from session where user_id=20 order by id asc`)
		s.Require().NoError(err)
		s.Require().Len(sessions, 4)
		s.Require().Equal(model.SessionMembersProcess, sessions[3].Status)
		s.Require().Equal(int64(20), sessions[3].UserID)
		s.Require().Equal(`{"choose":"2","tourMapping":"{\"1\":{\"id\":101,\"title\":\"tour1\",\"date\":\"2024-03-21\"},\"2\":{\"id\":102,\"title\":\"tour2\",\"date\":\"2024-03-22\"},\"3\":{\"id\":103,\"title\":\"tour3\",\"date\":\"2024-03-23\"}}"}`, sessions[3].Data)
		s.Require().True(sessions[0].Closed)
		s.Require().True(sessions[1].Closed)
		s.Require().True(sessions[2].Closed)
		s.Require().True(sessions[3].Closed)
	})
}
