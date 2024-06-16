package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"string,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var db *pgxpool.Pool
var redisClient *redis.Client
var ctx context.Context = context.Background()

func SetConnections(pool *pgxpool.Pool, client *redis.Client) {
	db = pool
	redisClient = client
}
