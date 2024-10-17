package pgx

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/storage"
	"go.uber.org/zap"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	chatRepo *ChatRepository
	log      *zap.Logger
	pool     *pgxpool.Pool
}

func New(pool *pgxpool.Pool, log *zap.Logger) *Storage {
	s := &Storage{
		pool:     pool,
		log:      log.With(zap.String("layer", "Storage")),
		chatRepo: NewChatRepository(pool, log),
	}
	return s
}

func (s *Storage) Chat() storage.ChatRepository {
	return s.chatRepo
}
