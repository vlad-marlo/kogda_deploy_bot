package storage

import (
	"context"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/model"
	"time"
)

type ChatRepository interface {
	GetConsumersIDs(ctx context.Context) ([]int64, error)
	GetChatsWithNoMessageInIt(ctx context.Context, date time.Time, msg string) ([]int64, error)
	CreateChat(ctx context.Context, id int64, chatType model.ChatType) error
	SafeThatMessageHadBeenSend(ctx context.Context, id int64, date time.Time, msg string) error
}

type Storage interface {
	Chat() ChatRepository
}
