package telebot

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"time"
)

var botCreatedDate = time.Date(2024, time.April, 23, 0, 0, 0, 0, time.Local)

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

	today := time.Now()
	days := int(today.Sub(botCreatedDate).Hours()/24) + 1

	return ctx.Send(fmt.Sprintf(startMsg, days))
}

func (ctrl *Controller) Abet() {}

func (ctrl *Controller) HandleTick() {
	now := time.Now().In(ctrl.cfg.App.Location)
	hour, minute, second := now.Clock()
	if hour == 13 && minute == 59 && second == 00 {
		go ctrl.Abet()
	}
}
