package application

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/frolmr/metrics/internal/server/config"
	"github.com/frolmr/metrics/internal/server/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplication(t *testing.T) {
	t.Run("http server", func(t *testing.T) {
		cfg := &config.Config{
			Scheme:      "http",
			HTTPAddress: "localhost:0",
		}
		log, logErr := logger.NewLogger()
		assert.NoError(t, logErr)
		app := NewApplication(cfg, log)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := app.RunServer()
			if err != nil && err != http.ErrServerClosed {
				require.NoError(t, err)
			}
		}()

		time.Sleep(100 * time.Millisecond)

		{
			err := app.Shutdown(ctx)
			require.NoError(t, err)
		}
	})

	t.Run("grpc server", func(t *testing.T) {
		cfg := &config.Config{
			Scheme:      "grpc",
			HTTPAddress: "localhost:0",
		}
		log, logErr := logger.NewLogger()
		assert.NoError(t, logErr)
		app := NewApplication(cfg, log)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := app.RunServer()
			require.NoError(t, err)
		}()

		time.Sleep(100 * time.Millisecond)

		{
			err := app.Shutdown(ctx)
			require.NoError(t, err)
		}
	})
}
