package telebot

import (
	"context"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/config"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/controller"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
)

type Controller struct {
	bot      *telebot.Bot
	log      *zap.Logger
	cfg      *config.Config
	stopChan chan struct{}
}

var _ controller.Controller = (*Controller)(nil)

func New(log *zap.Logger, cfg *config.Config) (*Controller, error) {
	botSettings := telebot.Settings{
		Token: cfg.Telegram.TelegramSecretKey,
		Poller: &telebot.LongPoller{
			Timeout: cfg.Telegram.Timeout(),
		},
	}

	bot, err := telebot.NewBot(botSettings)
	if err != nil {
		return nil, err
	}

	ctrl := &Controller{
		log: log,
		cfg: cfg,
		bot: bot,
	}

	ctrl.configureRoutes()

	return ctrl, nil
}

func (ctrl *Controller) configureRoutes() {
	ctrl.bot.Handle("/start", ctrl.HandleStart)
}

func (ctrl *Controller) Start(context.Context) error {
	ctrl.stopChan = make(chan struct{})
	go ctrl.bot.Start()
	go ctrl.Poll()
	ctrl.log.Info("started bot")
	return nil
}

func (ctrl *Controller) Stop(context.Context) error {
	defer close(ctrl.stopChan)
	ctrl.bot.Stop()
	ctrl.log.Info("stopped bot")
	return nil
}
