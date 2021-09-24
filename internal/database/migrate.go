package database

import (
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"

	"go.uber.org/zap"
)

func Migrate(db *sql.DB, command string) error {
	var err error

	switch command {
	case "up":
		if err = goose.Up(db, "migrations"); err != nil {
			zap.L().Error("migration failed", zap.Error(err))
		}
	case "up-by-one":
		if err = goose.UpByOne(db, "migrations"); err != nil {
			zap.L().Error("migration failed", zap.Error(err))
		}
	case "down":
		if err = goose.Down(db, "migrations"); err != nil {
			zap.L().Error("migration failed", zap.Error(err))
		}
	default:
		zap.L().Error("Invalid command for 'migration'", zap.String("flag", command))
		return errors.New("invalid command")
	}
	return err
}
