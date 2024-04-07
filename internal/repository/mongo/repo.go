package mongo

import (
	"context"

	"github.com/oke11o/go-telegram-bot/internal/model"
)

const DBType = "mongo"

func New() *Repo {
	return &Repo{}
}

type Repo struct {
}

func (r *Repo) SaveIncome(ctx context.Context, income model.IncomeRequest) (model.IncomeRequest, error) {
	panic("implement mongo SaveIncome()")
}

func (r *Repo) SaveUser(ctx context.Context, income model.User) (model.User, error) {
	panic("implement mongo SaveUser()")
}

func (r *Repo) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	panic("implement mongo GetUserByUsername()")
}

func (r *Repo) SetUserIsManager(ctx context.Context, userID int64, isManager bool) error {
	panic("implement mongo SetUserIsManager()")
}

func (r *Repo) SaveSession(ctx context.Context, session model.Session) (model.Session, error) {
	panic("implement mongo SaveSession()")
}

func (r *Repo) CloseSession(ctx context.Context, session model.Session) error {
	panic("implement mongo CloseSession()")
}

func (r *Repo) GetSession(ctx context.Context, userID int64) (model.Session, error) {
	panic("implement mongo GetSession()")
}

func (r *Repo) SaveTournament(ctx context.Context, tournament model.Tournament) (model.Tournament, error) {
	panic("implement mongo SaveTournament()")
}

func (r *Repo) GetOpenedTournaments(ctx context.Context) ([]model.Tournament, error) {
	panic("implement mongo GetOpenedTournaments()")
}
