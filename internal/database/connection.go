package database

import (
	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

// NewConnection returns a new database connection.
func NewConnection(dsn, driver string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		zap.L().Error("failed to create a database connection", zap.Error(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		zap.L().Error("failed ping the database", zap.Error(err))
		return nil, err
	}

	return db, nil
}
