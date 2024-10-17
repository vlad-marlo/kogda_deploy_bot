package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/config"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
	"os"
	"testing"
)

func testCfg(t testing.TB) *config.Config {
	connString := os.Getenv("TEST_DB_URI")
	if connString == "" {
		t.Skip("db uri was not provided")
	}
	return &config.Config{
		Postgres: config.PostgresConfig{
			TestURI: connString,
		},
	}
}

func TestClient_P(t *testing.T) {
	t.Run("non nil client", func(t *testing.T) {
		cli := &pgxpool.Pool{}
		assert.Empty(t, cli)
	})
}

func TestNew_Positive(t *testing.T) {
	cfg := testCfg(t)
	lc := fxtest.NewLifecycle(t)
	cli, err := NewWithFx(lc, cfg, zap.L())
	assert.NoError(t, err)
	assert.NotNil(t, cli)
}

func TestNew_Negative_BadConfig(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	cli, err := NewWithFx(
		lc,
		&config.Config{
			Postgres: config.PostgresConfig{
				TestURI: "bad uri",
			},
		},
		zap.L(),
	)
	assert.Nil(t, cli)
	assert.Error(t, err)
}

func TestNewTest_DefaultClient(t *testing.T) {
	cli, td := NewTest(t)
	defer td()
	assert.NoError(t, cli.Ping(context.Background()))
}

func TestNewTest_BadClient(t *testing.T) {
	require.NoError(t, os.Setenv("TEST_DB_URI", "postgres://postgres:password@localhost:5432/unknonnnonnononnnononno"))
	_, _ = NewTest(t)
	t.Log("test is unexpectedly was not skipped")
	t.Fail()
}

func TestTD(t *testing.T) {
	td := teardown(&pgxpool.Pool{}, "")
	assert.Panics(t, td)
}

func TestBadCli(t *testing.T) {
	cli := BadCli(t)
	assert.Error(t, cli.Ping(context.Background()))
}
