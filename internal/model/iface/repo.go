package iface

import (
	"context"

	"github.com/oke11o/go-telegram-bot/internal/model"
)

type Repo interface {
	SaveIncome(ctx context.Context, income model.IncomeRequest) (model.IncomeRequest, error)
	SaveUser(ctx context.Context, user model.User) (model.User, error)
	SetUserIsManager(ctx context.Context, userID int64, isManager bool) error
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	SaveSession(ctx context.Context, session model.Session) (model.Session, error)
	GetSession(ctx context.Context, userID int64) (model.Session, error)
	SaveTournament(ctx context.Context, tournament model.Tournament) (model.Tournament, error)
	CloseSession(ctx context.Context, session model.Session) error
}
