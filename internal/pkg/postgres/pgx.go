package postgres

import (
	"context"
	"fmt"
	pgxzap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/config"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/pkg/retry"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	RetryAttempts = 4
	RetryDelay    = 500 * time.Millisecond
)

// New opens new postgres connection, configures it and return prepared client.
func New(connString string, log *zap.Logger) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	log.Info("initializing postgres client")

	c, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error while parsing db uri: %w", err)
	}

	var lvl = tracelog.LogLevelError
	c.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzap.NewLogger(log),
		LogLevel: lvl,
	}
	c.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("postgres: init pgxpool: %w", err)
	}

	log.Info("created postgres client")
	return pool, nil
}

func NewWithFx(lc fx.Lifecycle, cfg *config.Config, log *zap.Logger) (*pgxpool.Pool, error) {
	pool, err := New(cfg.Postgres.URI(), log)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return retry.TryWithAttemptsCtx(ctx, pool.Ping, RetryAttempts, RetryDelay)
		},
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

// NewTest prepares test client.
//
// If error occurred while creating connection then test will be skipped.
// Second argument cleanup function to close connection and rollback all changes.
func NewTest(t testing.TB) (*pgxpool.Pool, func()) {
	t.Helper()

	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DB_URI"))
	if err != nil {
		t.Skipf("can not create pool: %v", err)
	}

	if err = retry.TryWithAttemptsCtx(context.Background(), pool.Ping, RetryAttempts, RetryDelay); err != nil {
		t.Skipf("can not get access to db: %v", err)
	}
	return pool, func() {
		teardown(pool)()
	}
}

func NewTestWithTearDown(t testing.TB) (*pgxpool.Pool, func(...string)) {
	t.Helper()

	pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DB_URI"))
	if err != nil {
		t.Skipf("can not create pool: %v", err)
	}

	if err = retry.TryWithAttemptsCtx(context.Background(), pool.Ping, RetryAttempts, RetryDelay); err != nil {
		t.Skipf("can not get access to db: %v", err)
	}
	return pool, func(tables ...string) {
		teardown(pool, tables...)()
	}
}

func BadCli(t testing.TB) *pgxpool.Pool {
	t.Helper()

	pool, err := pgxpool.New(context.Background(), "postgresql://postgres:postgres@localhost:4321/unknown_db")
	if err != nil {
		t.Skipf("can not create pool: %v", err)
	}

	if err = retry.TryWithAttemptsCtx(context.Background(), pool.Ping, RetryAttempts, RetryDelay); err == nil {
		t.Skip("must have no connection to database")
	}
	return pool
}

// teardown return func for defer it to clear tables.
//
// Always pass one or more tables in it.
func teardown(pool *pgxpool.Pool, tables ...string) func() {
	return func() {
		_, _ = pool.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s CASCADE;", strings.Join(tables, ", ")))
		pool.Close()
	}
}
