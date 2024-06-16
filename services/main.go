package services

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type JSONResponse struct {
	Error   bool        `json:"error,omitempty"`
	Message string      `json:"string,omitempty"`
	Data    interface{} `json:"contact,omitempty"`
}

var db *pgxpool.Pool
var redisClient *redis.Client
var ctx context.Context = context.Background()

func SetConnections(pool *pgxpool.Pool, client *redis.Client) {
	db = pool
	redisClient = client
}

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
