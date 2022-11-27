package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) (*zap.Logger, error) {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("could not parse level: %w", err)
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		l,
	))

	return logger, nil
}
