package telebot

import (
	"context"
	"fmt"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/model"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
	"time"
)

var botCreatedDate = time.Date(2024, time.April, 23, 0, 0, 0, 0, time.Local)

func (ctrl *Controller) getDays() int {
	today := time.Now().In(ctrl.cfg.App.Location)
	days := int(today.Sub(botCreatedDate).Hours()/24) + 1
	return days
}

func (ctrl *Controller) HandleStart(ctx telebot.Context) error {
	chat := ctx.Chat()
	switch chat.Type {
	case telebot.ChatPrivate:
	case telebot.ChatGroup:
	case telebot.ChatChannel:
	case telebot.ChatChannelPrivate:
	case telebot.ChatSuperGroup:
	}
	ctrl.log.Debug("Got request")
	err := ctrl.storage.Chat().CreateChat(
		context.Background(),
		chat.ID,
		model.ChatType(chat.Type),
	)
	if err != nil {
		ctrl.log.Error(
			"Failed to create chat",
			zap.Error(err),
			zap.Int64("chat_id", chat.ID),
		)
	}

	return ctx.Send(fmt.Sprintf(startMsg, ctrl.getDays()))
}

func (ctrl *Controller) Abet() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	ids, err := ctrl.storage.Chat().GetConsumersIDs(ctx)
	if err != nil {
		ctrl.log.Error("error while getting users", zap.Error(err))
		return
	}
	ctrl.log.Info("got ids", zap.Int64s("chats", ids))
	for _, id := range ids {
		go func() {
			msg, err := ctrl.bot.Send(telebot.ChatID(id), fmt.Sprintf(
				abetMsg,
				ctrl.getDays(),
			))
			if err != nil {
				ctrl.log.Error("error while sending abet message", zap.Error(err), zap.Int64("chat_id", id))
				return
			}
			ctrl.log.Info("sent abet message", zap.Int64("chat_id", id), zap.Any("msg", msg))
		}()
	}
}

func (ctrl *Controller) HandleTick() {
	now := time.Now().In(ctrl.cfg.App.Location)
	hour, minute, second := now.Clock()
	abetPredicate := hour == 13 && minute == 59 && second == 0
	if abetPredicate || second%10 == 0 {
		ctrl.log.Info("go ctrl.Abet()")
		go ctrl.Abet()
	}
}
