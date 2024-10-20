package telebot

import (
	"context"
	"fmt"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/model"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/pkg/retry"
	"go.uber.org/zap"
	"gopkg.in/telebot.v4"
	"sync"
	"time"
)

var botCreatedDate = time.Date(2024, time.April, 23, 0, 0, 0, 0, time.Local)

func (ctrl *Controller) getDays() int {
	today := time.Now().In(ctrl.cfg.App.Location)
	days := int(today.Sub(botCreatedDate).Hours()/24) + 1
	return days
}

func (ctrl *Controller) HandleStop(ctx telebot.Context) error {
	chat := ctx.Chat()
	ctrl.log.Debug("deleting chat")
	err := ctrl.storage.Chat().DeleteChat(context.Background(), chat.ID)
	if err != nil {
		ctrl.log.Error(
			"Failed to delete chat",
			zap.Error(err),
			zap.Int64("chat_id", chat.ID),
		)
	}

	return ctx.Send(fmt.Sprintf(startMsg, ctrl.getDays()))
}

func (ctrl *Controller) HandleStart(ctx telebot.Context) error {
	chat := ctx.Chat()
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

func (ctrl *Controller) HandleHelp(ctx telebot.Context) error {
	ctrl.log.Info("Handling help")
	return ctx.Send(infoMsg)
}

func (ctrl *Controller) spam(msg string) error {
	now := time.Now().In(ctrl.cfg.App.Location)
	ctrl.log.Info("now", zap.String("now", now.String()))

	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	defer cancel()

	ids, err := ctrl.storage.Chat().GetChatsWithNoMessageInIt(ctx, now, msg)
	if err != nil {
		ctrl.log.Error("error while getting users", zap.Error(err))
		return err
	}

	ctrl.log.Info("got ids", zap.Int64s("chats", ids))
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)

		go func() {
			defer wg.Done()

			defer func() {
				if err := retry.TryWithAttempts(
					func() error {
						return ctrl.storage.Chat().SafeThatMessageHadBeenSend(
							context.Background(),
							id,
							now,
							msg,
						)
					},
					5,
					200*time.Millisecond,
				); err != nil {
					ctrl.log.Error("failed to spam", zap.Error(err), zap.Int64("chat_id", id))
				}
			}()

			if _, err := ctrl.bot.Send(telebot.ChatID(id), fmt.Sprintf(
				msg,
				ctrl.getDays(),
			)); err != nil {
				ctrl.log.Error("error while sending abet message", zap.Error(err), zap.Int64("chat_id", id))
				return
			}
			ctrl.log.Info("sent abet message", zap.Int64("chat_id", id))
		}()
	}
	wg.Wait()
	return nil
}

func (ctrl *Controller) HandleAbet(tbContext telebot.Context) error {
	ctrl.log.Info("Handling spam")
	err := ctrl.spam(newDayMsg)
	if err != nil {
		ctrl.log.Info("error while spamming chat", zap.Error(err))
		return tbContext.Send("Что-то не получилось получить пользователей, кринж какой то")
	}
	return tbContext.Send("Да да, сообщения отправлены)))")
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
			for range 3 {
				msg, err := ctrl.bot.Send(
					telebot.ChatID(id),
					fmt.Sprintf(
						abetMsg,
						ctrl.getDays(),
					),
				)
				if err != nil {
					ctrl.log.Error(
						"error while sending abet message",
						zap.Error(err),
						zap.Int64("chat_id", id),
					)
					return
				}
				ctrl.log.Info(
					"sent abet message",
					zap.Int64("chat_id", id),
					zap.Any("msg", msg),
				)
			}
		}()
	}
}

func (ctrl *Controller) HandleTick() {
	now := time.Now().In(ctrl.cfg.App.Location)
	date := now.Weekday()
	hour, minute, second := now.Clock()
	abetPredicate := hour == 13 && minute == 59 && second == 0
	hourTillAbetPredicate := hour == 12 && minute == 59 && second == 0
	newDayPredicate := hour == 0 && minute == 0 && second == 0
	newMorningPredicate := hour == 8 && minute == 0 && second == 0
	switch date {
	case time.Tuesday, time.Friday:
	default:
		abetPredicate = false
		hourTillAbetPredicate = false
	}
	switch {
	case abetPredicate:
		ctrl.log.Info("go ctrl.Abet()")
		ctrl.Abet()
	case newDayPredicate:
		ctrl.log.Info("go ctrl.spam()")
		err := ctrl.spam(newDayMsg)
		if err != nil {
			ctrl.log.Error("error while spamming abet", zap.Error(err))
		}
	case hourTillAbetPredicate:
		ctrl.log.Info("go ctrl.spam(hourTillAbetMsg)")
		err := ctrl.spam(hourTillAbetMsg)
		if err != nil {
			ctrl.log.Error("error while spamming abet", zap.Error(err))
		}
	case newMorningPredicate:
		ctrl.log.Info("go ctrl.newDayMsg()")
		err := ctrl.spam(newMorningMsg)
		if err != nil {
			ctrl.log.Error("error while spamming abet", zap.Error(err))
		}
	}
}
