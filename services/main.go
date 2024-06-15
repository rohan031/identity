package services

import "github.com/jackc/pgx/v5/pgxpool"

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"string,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var db *pgxpool.Pool

func SetDbPool(pool *pgxpool.Pool) {
	db = pool
}
