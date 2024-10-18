package telebot

import (
	"context"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/config"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/controller"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/storage"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

type Controller struct {
	bot      *telebot.Bot
	log      *zap.Logger
	cfg      *config.Config
	stopChan chan struct{}
	storage  storage.Storage
}

var _ controller.Controller = (*Controller)(nil)

func New(log *zap.Logger, cfg *config.Config, store storage.Storage) (*Controller, error) {
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
		log:     log,
		cfg:     cfg,
		bot:     bot,
		storage: store,
	}

	ctrl.configureRoutes()

	return ctrl, nil
}

func (ctrl *Controller) configureRoutes() {
	ctrl.log.Info("admins ids", zap.Int64s("ids", ctrl.cfg.Telegram.AdminIDs))
	ctrl.bot.Handle("/help", ctrl.HandleHelp)
	ctrl.bot.Handle("/start", ctrl.HandleStart)
	adminOnly := ctrl.bot.Group()
	adminOnly.Use(middleware.Whitelist(ctrl.cfg.Telegram.AdminIDs...))
	adminOnly.Handle("/abet", ctrl.HandleAbet)
	adminOnly.Handle("/spam", ctrl.HandleAbet)
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
