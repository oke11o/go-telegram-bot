package model

type Chat struct {
	ID       int64   `db:"id"`
	Type     string  `db:"type"`
	Title    string  `db:"title"`
	Photo    *string `db:"photo"`
	Location *string `db:"location"`
}
