package logger_test

import (
	"testing"

	"github.com/calmonr/scaleid/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("could not parse level", func(t *testing.T) {
		t.Parallel()

		l, err := logger.New("invalid")
		assert.ErrorContains(t, err, "could not parse level")

		assert.Nil(t, l)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		l, err := logger.New("debug")
		assert.NoError(t, err)

		assert.Equal(t, zap.DebugLevel, zapcore.LevelOf(l.Core()))
	})
}
