package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

func NewHandler(db *pgxpool.Pool, rdb *redis.Client) *Handler {
	return &Handler{
		DB:  db,
		RDB: rdb,
	}
}

func (h *Handler) getCache(ctx context.Context, key string, dest any) bool {
	val, err := h.RDB.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return false
	}
	return json.Unmarshal([]byte(val), dest) == nil
}

func (h *Handler) setCache(ctx context.Context, key string, data any, ttl time.Duration) {
	bytes, _ := json.Marshal(data)
	h.RDB.Set(ctx, key, bytes, ttl)
}

func (h *Handler) invalidateCache(ctx context.Context, keys ...string) {
	h.RDB.Del(ctx, keys...)
}
