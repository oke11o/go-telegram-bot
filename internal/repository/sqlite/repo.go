package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/oke11o/go-telegram-bot/internal/config"
	"github.com/oke11o/go-telegram-bot/internal/model"
)

const DBType = "sqlite"

func New(cfg config.SqliteConfig) (*Repo, error) {
	db, err := sql.Open("sqlite3", cfg.File)
	if err != nil {
		return nil, fmt.Errorf("sql.Open() err: %w", err)
	}
	dbx := sqlx.NewDb(db, "sqlite3")
	err = runMigrate(cfg.MigrationPath, cfg.File)
	if err != nil {
		return nil, fmt.Errorf("runMigrate() err: %w", err)
	}

	return &Repo{db: dbx}, nil
}

func NewWithDB(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) SaveIncome(ctx context.Context, income model.IncomeRequest) (model.IncomeRequest, error) {
	q := `insert into income_request (from_id,message_id,reply_to_message_id,request_id,message,username,text) 
values (:from_id,:message_id,:reply_to_message_id,:request_id, :message, :username, :text)`
	raw, err := r.db.NamedExecContext(ctx, q, income)
	if err != nil {
		return income, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := raw.LastInsertId()
	if err != nil {
		return income, fmt.Errorf("raw.LastInsertId() err: %w", err)
	}
	income.ID = id
	return income, nil
}

func (r *Repo) SetUserIsManager(ctx context.Context, userID int64, isManager bool) error {
	q := `update user set is_manager=? 
                where id=?`
	res, err := r.db.ExecContext(ctx, q, isManager, userID)
	if err != nil {
		return fmt.Errorf("db.ExecContext() err: %w", err)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("res.RowsAffected() err: %w", err)
	}
	if cnt == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repo) SaveUser(ctx context.Context, user model.User) (model.User, error) {
	u, err := r.GetUser(ctx, user.ID)
	if err != nil {
		return u, fmt.Errorf("db.Get() err: %w", err)
	}
	if u.ID == 0 {
		err = r.insertUser(ctx, user)
		if err != nil {
			return user, fmt.Errorf("insertUser() err: %w", err)
		}
		return user, nil
	}
	user.IsManager = u.IsManager
	user.IsMaintainer = u.IsMaintainer
	err = r.updateUser(ctx, user)
	if err != nil {
		return user, fmt.Errorf("updateUser() err: %w", err)
	}

	return user, nil
}

func (r *Repo) insertUser(ctx context.Context, user model.User) error {
	q := `insert into user (id,username,first_name,last_name,language_code,is_bot,is_maintainer,is_manager) 
values (:id,:username,:first_name,:last_name,:language_code,:is_bot,:is_maintainer,:is_manager)`
	_, err := r.db.NamedExecContext(ctx, q, user)
	if err != nil {
		return fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	return nil
}

func (r *Repo) updateUser(ctx context.Context, user model.User) error {
	q := `update user set username=:username,first_name=:first_name,last_name=:last_name,
                language_code=:language_code,is_bot=:is_bot,is_maintainer=:is_maintainer,is_manager=:is_manager 
                where id=:id`
	_, err := r.db.NamedExecContext(ctx, q, user)
	if err != nil {
		return fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	return nil
}

func (r *Repo) GetUser(ctx context.Context, id int64) (model.User, error) {
	user := model.User{}
	q := `select id,username,first_name,last_name,language_code,is_bot,is_maintainer,is_manager
from user where id=?`
	err := r.db.GetContext(ctx, &user, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return user, nil
	}
	return user, err
}

func (r *Repo) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	user := model.User{}
	q := `select id,username,first_name,last_name,language_code,is_bot,is_maintainer,is_manager
from user where username=?`
	err := r.db.GetContext(ctx, &user, q, username)
	if errors.Is(err, sql.ErrNoRows) {
		return user, nil
	}
	return user, err
}

func (r *Repo) SaveSession(ctx context.Context, session model.Session) (model.Session, error) {
	err := session.PrepareToSave()
	if err != nil {
		return session, fmt.Errorf("session.PrepareToSave() err: %w", err)
	}
	q := `insert into session (user_id,data,status,created_at,updated_at)
values (:user_id,:data,:status,:created_at,:updated_at)`
	raw, err := r.db.NamedExecContext(ctx, q, session)
	if err != nil {
		return session, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := raw.LastInsertId()
	if err != nil {
		return session, fmt.Errorf("raw.LastInsertId() err: %w", err)
	}
	session.ID = id

	return session, nil
}

func (r *Repo) CloseSession(ctx context.Context, session model.Session) error {
	q := `update session set closed=1,updated_at=? where user_id=?`
	_, err := r.db.ExecContext(ctx, q, time.Now().Format(time.RFC3339), session.UserID)
	if err != nil {
		return fmt.Errorf("db.ExecContext() err: %w", err)
	}
	return nil
}

func (r *Repo) GetSession(ctx context.Context, userID int64) (model.Session, error) {
	ses := model.Session{}
	q := `select id,user_id,data,status,created_at,updated_at from session where user_id=? order by id desc limit 1`
	err := r.db.GetContext(ctx, &ses, q, userID)
	if err == nil {
		err = ses.AfterGet()
		if err != nil {
			return ses, fmt.Errorf("session.AfterGet() err: %w", err)
		}
		return ses, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ses, nil
	}

	return ses, fmt.Errorf("db.GetContext() err: %w", err)
}

func (r *Repo) SaveTournament(ctx context.Context, tournament model.Tournament) (model.Tournament, error) {
	q := `insert into tournament (title,date,status,created_by,created_at,updated_at)
values (:title,:date,:status,:created_by,:created_at,:updated_at)`
	raw, err := r.db.NamedExecContext(ctx, q, tournament)
	if err != nil {
		return tournament, fmt.Errorf("db.NamedExecContext() err: %w", err)
	}
	id, err := raw.LastInsertId()
	if err != nil {
		return tournament, fmt.Errorf("raw.LastInsertId() err: %w", err)
	}
	tournament.ID = id

	return tournament, nil
}

func (r *Repo) GetOpenedTournaments(ctx context.Context) ([]model.Tournament, error) {
	q := `select id,title,date,status,created_by,created_at,updated_at from tournament where status in (?)`
	q, args, err := sqlx.In(q, []model.TournamentStatus{model.TournamentStatusCreated, model.TournamentStatusInProgress})
	if err != nil {
		return nil, fmt.Errorf("sqlx.In() err: %w", err)
	}
	q = r.db.Rebind(q)
	var tournaments []model.Tournament
	err = r.db.SelectContext(ctx, &tournaments, q, args...)
	if err != nil {
		return nil, fmt.Errorf("db.SelectContext() err: %w", err)
	}
	return tournaments, nil
}
