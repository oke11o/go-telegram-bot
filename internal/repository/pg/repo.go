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

func (r *Repo) CloseSession(ctx context.Context, session model.Session) error {
	panic("implement pg CloseSession()")
}

func (r *Repo) GetOpenedSession(ctx context.Context, userID int64) (model.Session, error) {
	panic("implement pg GetOpenedSession()")
}

func (r *Repo) SaveTournament(ctx context.Context, tournament model.Tournament) (model.Tournament, error) {
	panic("implement pg SaveTournament()")
}

func (r *Repo) GetOpenedTournaments(ctx context.Context) ([]model.Tournament, error) {
	panic("implement pg GetOpenedTournaments()")
}

func (r *Repo) GetMemberTournaments(ctx context.Context, id int64) ([]model.Tournament, error) {
	panic("implement pg GetMemberTournaments()")
}

func (r *Repo) AddPlayerToTournament(ctx context.Context, userID int64, tournamentID int64) error {
	panic("implement pg AddPlayerToTournament()")
}

func (r *Repo) RemovePlayerFromTournament(ctx context.Context, userID int64, tournamentID int64) error {
	panic("implement pg RemovePlayerFromTournament()")
}

func (r *Repo) GetTournamentsPlayers(ctx context.Context, tournamentID int64) ([]model.User, error) {
	panic("implement pg GetTournamentsPlayers()")
}

func (r *Repo) TournamentOpenedAll(ctx context.Context) ([]model.Tournament, error) {
	panic("implement pg TournamentOpenedAll()")
}

func (r *Repo) TournamentOpenedByManager(ctx context.Context, userID int64) ([]model.Tournament, error) {
	panic("implement pg TournamentOpenedByManager()")
}

func (r *Repo) TournamentStart(ctx context.Context, id int64) error {
	panic("implement pg TournamentStart()")
}

func (r *Repo) TournamentFinish(ctx context.Context, id int64) error {
	panic("implement pg TournamentFinish()")
}
