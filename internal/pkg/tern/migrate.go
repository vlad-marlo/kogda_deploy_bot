package tern

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/config"
	"github.com/vlad-marlo/kogda_deploy_bot/migrations"
	"go.uber.org/zap"
)

func Migrate(cfg *config.Config, log *zap.Logger) error {
	ctx := context.Background()
	log.Info("starting migration")

	conn, err := pgx.Connect(ctx, cfg.Postgres.URI())
	if err != nil {
		log.Info("failed to connect to database", zap.Error(err))
		return fmt.Errorf("pgx.Connect: %w", err)
	}

	defer func() {
		err = conn.Close(ctx)
		if err != nil {
			log.Error("failed to close connection", zap.Error(err))
		}

		log.Info("closed migrator pgx connection")
	}()

	log.Info("creating migrator")
	migrator, err := migrate.NewMigrator(ctx, conn, "versions")
	if err != nil {
		return fmt.Errorf("migrate.NewMigrator: %w", err)
	}

	//migrationRoot, err := fs.Sub(migrations.Migrations, "migrations")
	//if err != nil {
	//	log.Error("failed to sub migrations", zap.Error(err))
	//	return fmt.Errorf("fs.Sub: %w", err)
	//}

	log.Info("migrator created; loading migrations into migrator")
	err = migrator.LoadMigrations(migrations.Migrations)
	if err != nil {
		log.Info("migrator failed to load migrations", zap.Error(err))
		return fmt.Errorf("migrator.LoadMigrations: %w", err)
	}

	log.Info("migrator loaded; running migrations")

	err = migrator.Migrate(ctx)
	if err != nil {
		log.Info("migrator failed to run migrations", zap.Error(err))
		return fmt.Errorf("migrator.Migrate: %w", err)
	}
	log.Info("migrator successfully run migrations")
	return nil
}
