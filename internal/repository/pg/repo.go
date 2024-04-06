package pg

import (
	"context"

	"github.com/oke11o/go-telegram-bot/internal/model"
)

const DBType = "pg"

func New() *Repo {
	return &Repo{}
}

type Repo struct {
}

func (r *Repo) SaveIncome(ctx context.Context, income model.IncomeRequest) (model.IncomeRequest, error) {
	panic("implement pg SaveIncome()")
}

func (r *Repo) SaveUser(ctx context.Context, income model.User) (model.User, error) {
	panic("implement pg SaveUser()")
}

func (r *Repo) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	panic("implement pg GetUserByUsername()")
}

func (r *Repo) SetUserIsManager(ctx context.Context, userID int64, isManager bool) error {
	panic("implement pg SetUserIsManager()")
}

func (r *Repo) SaveSession(ctx context.Context, session model.Session) (model.Session, error) {
	panic("implement pg SaveSession()")
}
