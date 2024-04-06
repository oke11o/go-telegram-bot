package sqlite

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrate(source string, dbpath string) error {

	m, err := migrate.New("file://"+source, "sqlite3://"+dbpath)
	if err != nil {
		return fmt.Errorf("cant migrate.New, err: %w", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
