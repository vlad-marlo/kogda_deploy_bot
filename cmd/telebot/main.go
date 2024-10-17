package main

import (
	"github.com/vlad-marlo/kogda_deploy_bot/internal/config"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/controller"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/controller/telebot"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/pkg/postgres"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/pkg/tern"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		NewOptions(),
	).Run()
}

func NewOptions() fx.Option {
	return fx.Options(
		fx.NopLogger,
		fx.Provide(
			config.New,
			zap.NewProduction,
			fx.Annotate(telebot.New, fx.As(new(controller.Controller))),
			postgres.NewWithFx,
		),
		fx.Invoke(
			controller.RunFx,
			tern.Migrate,
		),
	)
}
