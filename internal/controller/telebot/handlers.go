package telebot

import (
	"fmt"
	"gopkg.in/telebot.v4"
	"time"
)

var botCreatedDate = time.Date(2024, time.April, 23, 0, 0, 0, 0, time.Local)

func (ctrl *Controller) HandleStart(ctx telebot.Context) error {
	today := time.Now()
	days := int(today.Sub(botCreatedDate).Hours()/24) + 1

	return ctx.Send(fmt.Sprintf(startMsg, days))
}
