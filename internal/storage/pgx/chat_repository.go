package pgx

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/model"
	"github.com/vlad-marlo/kogda_deploy_bot/internal/storage"
	"go.uber.org/zap"
)

var _ storage.ChatRepository = (*ChatRepository)(nil)

type ChatRepository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func NewChatRepository(pool *pgxpool.Pool, log *zap.Logger) *ChatRepository {
	return &ChatRepository{
		pool: pool,
		log:  log.With(zap.String("layer", "ChatRepository")),
	}
}

func (c *ChatRepository) CreateChat(ctx context.Context, id int64, chatType model.ChatType) error {
	const query = `INSERT INTO chats(id, chat_type)
VALUES ($1, $2);`
	c.log.Debug("creating chat", zap.Int64("id", id), zap.String("type", string(chatType)))
	_, err := c.pool.Exec(ctx, query, id, chatType)
	if err != nil {
		return fmt.Errorf("create chat: %w", err)
	}
	c.log.Info("inserted chat", zap.Int64("id", id), zap.String("type", string(chatType)))
	return nil
}

func (c *ChatRepository) GetConsumersIDs(ctx context.Context) ([]int64, error) {
	const query = `SELECT id FROM chats;`
	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.Query: %w", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64

		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		c.log.Debug("got chat", zap.Int64("id", id))

		ids = append(ids, id)
		c.log.Debug("ids", zap.Int64s("ids", ids))
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return ids, nil
}
