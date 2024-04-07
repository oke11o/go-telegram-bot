package model

import "time"

type TournamentStatus string

const (
	TournamentStatusCreated    TournamentStatus = "created"
	TournamentStatusInProgress TournamentStatus = "in_progress"
	TournamentStatusFinished   TournamentStatus = "finished"
)

func NewTournament(title, date string, createdBy int64) Tournament {
	return Tournament{
		Title:     title,
		Date:      date,
		Status:    TournamentStatusCreated,
		CreatedBy: createdBy,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}

type Tournament struct {
	ID        int64            `db:"id"`
	Title     string           `db:"title"`
	Date      string           `db:"date"`
	Status    TournamentStatus `db:"status"`
	CreatedBy int64            `db:"created_by"`
	CreatedAt string           `db:"created_at"`
	UpdatedAt string           `db:"updated_at"`
}
