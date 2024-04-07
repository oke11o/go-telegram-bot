package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/oke11o/go-telegram-bot/internal/config"
	"github.com/oke11o/go-telegram-bot/pgk/utils/str"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oke11o/go-telegram-bot/internal/model"
	"github.com/oke11o/go-telegram-bot/internal/repository/sqlite"
	"github.com/oke11o/go-telegram-bot/internal/service"
)

type testSender struct {
	assert    func(c tgbotapi.Chattable)
	returnMsg tgbotapi.Message
	returnErr error
}

func (t testSender) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if t.assert != nil {
		t.assert(c)
	}
	return t.returnMsg, t.returnErr
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {}

func (s *Suite) SetupTest() {
	s.dbCfg = config.SqliteConfig{
		File:          fmt.Sprintf("../../tests/db/test-%s.sqlite", str.RandStringRunes(8, "")),
		MigrationPath: "../../migrations/sqlite",
	}
	dbx, err := s.createDB(s.dbCfg)
	s.Require().NoError(err)
	s.dbx = dbx
	repo, err := sqlite.New(s.dbCfg)
	s.Require().NoError(err)
	s.repo = repo
	s.T().Logf("Start testing with db in file `%s`", s.dbCfg.File)
}

func (s *Suite) TearDownTest() {
	s.dbx.Close()
	//os.Remove(s.dbCfg.File)
}

func (s *Suite) TearDownSuite() {}

func (s *Suite) TestHandler_AddRemoveAdminCmd() {
	// arrange
	mainAdmin := model.User{ID: 1, Username: "main_admin", FirstName: "Main", LastName: "Admin", LanguageCode: "en", IsMaintainer: true}
	ctx := context.Background()
	_, err := s.repo.SaveUser(ctx, mainAdmin)
	s.Require().NoError(err)
	targetUser := model.User{ID: 2, Username: "target_user", FirstName: "Target", LastName: "User", LanguageCode: "ru"}
	_, err = s.repo.SaveUser(ctx, targetUser)
	s.Require().NoError(err)

	s.T().Run("cant find user", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal("I don't know the user unknown_user", msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 1, FirstName: "Main", LastName: "Admin", UserName: "main_admin", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 1, Type: "private", UserName: "main_admin", FirstName: "Main", LastName: "Admin"},
				Text: "/addAdmin @unknown_user"},
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		q := `select * from income_request`
		err = s.dbx.Select(&incomes, q)
		s.Require().NoError(err)
		s.Len(incomes, 1)

		var user model.User
		q = `select * from user where username=?`
		err = s.dbx.Get(&user, q, "unknown_user")
		s.Require().True(errors.Is(err, sql.ErrNoRows))
	})

	s.T().Run("successful find user", func(t *testing.T) {
		h := s.createHandler()
		responseMessages := []string{}
		responseIds := []int64{}
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				responseMessages = append(responseMessages, msg.Text)
				responseIds = append(responseIds, msg.ChatID)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 1, FirstName: "Main", LastName: "Admin", UserName: "main_admin", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 1, Type: "private", UserName: "main_admin", FirstName: "Main", LastName: "Admin"},
				Text: "/addAdmin @target_user"},
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		q := `select * from income_request`
		err = s.dbx.Select(&incomes, q)
		s.Require().NoError(err)
		s.Len(incomes, 2)

		var user model.User
		q = `select * from user where username=?`
		err = s.dbx.Get(&user, q, "target_user")
		s.Require().NoError(err)
		s.Require().Equal("Target", user.FirstName)
		s.Require().True(user.IsManager)

		s.Require().Len(responseMessages, 2)
		s.Require().ElementsMatch([]string{"Successful give permissions to user target_user", "Maintainer give you manager permissions"}, responseMessages)
		s.Require().ElementsMatch([]int64{1, 2}, responseIds)
	})

	s.T().Run("successful remove manager permissions user", func(t *testing.T) {
		h := s.createHandler()
		responseMessages := []string{}
		responseIds := []int64{}
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				responseMessages = append(responseMessages, msg.Text)
				responseIds = append(responseIds, msg.ChatID)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 1, FirstName: "Main", LastName: "Admin", UserName: "main_admin", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 1, Type: "private", UserName: "main_admin", FirstName: "Main", LastName: "Admin"},
				Text: "/removeAdmin @target_user"},
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		q := `select * from income_request`
		err = s.dbx.Select(&incomes, q)
		s.Require().NoError(err)
		s.Len(incomes, 3)

		var user model.User
		q = `select * from user where username=?`
		err = s.dbx.Get(&user, q, "target_user")
		s.Require().NoError(err)
		s.Require().Equal("Target", user.FirstName)
		s.Require().False(user.IsManager)

		s.Require().Len(responseMessages, 2)
		s.Require().ElementsMatch([]string{"Successful remove permissions to user target_user", "Maintainer remove your manager permissions"}, responseMessages)
		s.Require().ElementsMatch([]int64{1, 2}, responseIds)
	})
}

func (s *Suite) TestHandler_JustText() {
	h := s.createHandler()
	h.SetSender(testSender{
		assert: func(c tgbotapi.Chattable) {
			msg, ok := c.(tgbotapi.MessageConfig)
			s.Require().True(ok)
			s.Require().Equal("Choose action", msg.Text)
		},
	})
	err := h.HandleUpdate(context.Background(), tgbotapi.Update{
		UpdateID: 1,
		Message: &tgbotapi.Message{MessageID: 44,
			From: &tgbotapi.User{ID: 1, FirstName: "Main", LastName: "Admin", UserName: "main_admin", LanguageCode: "en"},
			Date: 1712312739,
			Chat: &tgbotapi.Chat{ID: 1, Type: "private", UserName: "main_admin", FirstName: "Main", LastName: "Admin"},
			Text: "🤟"},
	})
	s.Require().NoError(err)

	var incomes []model.IncomeRequest
	q := `select * from income_request`
	err = s.dbx.Select(&incomes, q)
	s.Require().NoError(err)
	s.Len(incomes, 1)

	q = `select * from user where username=?`
	var user model.User
	err = s.dbx.Get(&user, q, "main_admin")
	s.Require().NoError(err)
	s.Require().Equal("Admin", user.LastName)
}

func (s *Suite) TestHandler_Tournament() {
	// arrange
	mainAdmin := model.User{ID: 111, Username: "main_admin", FirstName: "Main", LastName: "Admin", LanguageCode: "en", IsMaintainer: true}
	ctx := context.Background()
	_, err := s.repo.SaveUser(ctx, mainAdmin)
	s.Require().NoError(err)

	s.T().Run("choose tg command", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal("Please text title of the tournament", msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 111, FirstName: "Main", LastName: "Admin", UserName: "main_admin", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 111, Type: "private", UserName: "main_admin", FirstName: "Main", LastName: "Admin"},
				Text: "/create"}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		q := `select * from income_request`
		err = s.dbx.Select(&incomes, q)
		s.Require().NoError(err)
		s.Len(incomes, 1)

		var sessions []model.Session
		q = `select * from session where user_id=?`
		err = s.dbx.Select(&sessions, q, mainAdmin.ID)
		s.Require().NoError(err)
		s.Require().Len(sessions, 1)
		s.Require().Equal(model.SessionCreateTournamentSetTitle, sessions[0].Status)
		s.Require().Equal(int64(111), sessions[0].UserID)
		s.Require().Equal("{}", sessions[0].Data)
		s.Require().False(sessions[0].Closed)
	})

	s.T().Run("text title", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal("Please text start date of the tournament", msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 111, FirstName: "Main", LastName: "Admin", UserName: "main_admin", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 111, Type: "private", UserName: "main_admin", FirstName: "Main", LastName: "Admin"},
				Text: "My new tournament"}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		q := `select * from income_request`
		err = s.dbx.Select(&incomes, q)
		s.Require().NoError(err)
		s.Len(incomes, 2)

		var sessions []model.Session
		q = `select * from session where user_id=? order by id asc`
		err = s.dbx.Select(&sessions, q, mainAdmin.ID)
		s.Require().NoError(err)
		s.Require().Len(sessions, 2)
		s.Require().Equal(model.SessionCreateTournamentSetDate, sessions[1].Status)
		s.Require().Equal(int64(111), sessions[1].UserID)
		s.Require().Equal(`{"title":"My new tournament"}`, sessions[1].Data)
		s.Require().False(sessions[0].Closed)
	})

	s.T().Run("text date", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal("Please text start date of the tournament", msg.Text)
			},
		})
		err = h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 111, FirstName: "Main", LastName: "Admin", UserName: "main_admin", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 111, Type: "private", UserName: "main_admin", FirstName: "Main", LastName: "Admin"},
				Text: "21.03.2024"}, // Notice!!!
		})
		s.Require().NoError(err)

		var incomes []model.IncomeRequest
		q := `select * from income_request`
		err = s.dbx.Select(&incomes, q)
		s.Require().NoError(err)
		s.Len(incomes, 3)

		var sessions []model.Session
		q = `select * from session where user_id=? order by id asc`
		err = s.dbx.Select(&sessions, q, mainAdmin.ID)
		s.Require().NoError(err)
		s.Require().Len(sessions, 3)
		s.Require().Equal(model.SessionCreateTournamentSetDate, sessions[2].Status)
		s.Require().Equal(int64(111), sessions[2].UserID)
		s.Require().Equal(`{"date":"21.03.2024","title":"My new tournament"}`, sessions[2].Data)
		s.Require().True(sessions[0].Closed)
		s.Require().True(sessions[1].Closed)
		s.Require().True(sessions[2].Closed)

		var tournament model.Tournament
		err = s.dbx.Get(&tournament, `select * from tournament where title=?`, "My new tournament")
		s.Require().NoError(err)
		s.Require().Equal("My new tournament", tournament.Title)
		s.Require().Equal("21.03.2024", tournament.Date)
		s.Require().Equal(int64(111), tournament.CreatedBy)
		s.Require().Equal(model.TournamentStatusCreated, tournament.Status)
	})
}

func (s *Suite) createHandler() *Handler {
	income := service.NewIncomeServce(s.repo)
	h := &Handler{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		income: income,
		repo:   s.repo,
		cfg:    s.cfg,
	}
	return h
}
