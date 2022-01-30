package telemetry

import (
	"go.uber.org/zap"
)

// NewLogger returns a new logger.
func NewLogger(cfg *zap.Config) (*zap.Logger, error) {
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)

	return logger, nil
}
