package database

import (
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

func Migrate(db *sql.DB, command string) error {
	var err error

	switch command {
	case "up":
		if err = goose.Up(db, "migrations"); err != nil {
			log.Error().Err(err).Msg("Migration failed")
		}
	case "up-by-one":
		if err = goose.UpByOne(db, "migrations"); err != nil {
			log.Error().Err(err).Msg("Migration failed")
		}
	case "down":
		if err = goose.Down(db, "migrations"); err != nil {
			log.Error().Err(err).Msg("Migration failed")
		}
	default:
		log.Warn().Msgf("Invalid command for 'migration' flag: '%v'", command)
		return errors.New("invalid command")
	}
	return err
}
