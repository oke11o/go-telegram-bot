package model

type TournamentStatus string

const (
	TournamentStatusCreated    TournamentStatus = "created"
	TournamentStatusInProgress TournamentStatus = "in_progress"
	TournamentStatusFinished   TournamentStatus = "finished"
)

type Tournament struct {
	ID        int64            `db:"id"`
	Title     string           `db:"title"`
	Date      string           `db:"date"`
	Status    TournamentStatus `db:"status"`
	CreatedBy int64            `db:"created_by"`
	CreatedAt string           `db:"created_at"`
	UpdatedAt string           `db:"updated_at"`
}
