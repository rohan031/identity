package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"string,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var db *pgxpool.Pool
var ctx context.Context = context.Background()

func SetDbPool(pool *pgxpool.Pool) {
	db = pool
}
