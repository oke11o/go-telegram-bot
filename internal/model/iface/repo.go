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
}
