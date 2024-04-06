package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type SessionStatus string

const (
	SessionCreateTournamentProcess  SessionStatus = "create_tournament_process"
	SessionCreateTournamentAskTitle SessionStatus = "create_tournament_ask_title"
	SessionCreateTournamentAskDate  SessionStatus = "create_tournament_ask_date"
	SessionStatusClosed             SessionStatus = "closed"
)

type Session struct {
	ID        int64             `db:"id"`
	UserID    int64             `db:"user_id"`
	Data      string            `db:"data"`
	dataArgs  map[string]string `db:"-"`
	Status    SessionStatus     `db:"status"`
	CreatedAt string            `db:"created_at"`
	UpdatedAt string            `db:"updated_at"`
}

func (s *Session) GetArg(key string) (string, bool) {
	value, ok := s.dataArgs[key]
	return value, ok
}

func (s *Session) SetArg(key string, value string) {
	s.dataArgs[key] = value
}

func (s *Session) RemoveArg(key string) {
	delete(s.dataArgs, key)
}

func (s *Session) SetStatus(status SessionStatus) {
	s.Status = status
}

func (s *Session) PrepareToSave() error {
	b, err := json.Marshal(s.dataArgs)
	if err != nil {
		return fmt.Errorf("json.Marshal() err: %w", err)
	}
	s.Data = string(b)
	return nil
}

func NewCreateTournamentSession(userID int64) Session {
	return Session{
		UserID:    userID,
		Status:    SessionCreateTournamentProcess,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}
