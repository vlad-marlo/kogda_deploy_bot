package storage

import (
	"context"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/model"
)

type ChatRepository interface {
	GetConsumersIDs(ctx context.Context) ([]int64, error)
	CreateChat(ctx context.Context, id int64, chatType model.ChatType) error
}

type Storage interface {
	Chat() ChatRepository
}
