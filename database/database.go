package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ctx context.Context = context.Background()

const (
	defaultMaxConns          = int32(5)
	defaultMinConns          = int32(0)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

func dbConfig() (*pgxpool.Config, error) {
	DATABASE_URL := os.Getenv("DB_DSN")

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig, nil
}

func CreatePool() (*pgxpool.Pool, error) {
	config, err := dbConfig()
	if err != nil {
		log.Println("Failed to create config!!")
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Println("Error while creating connection to the database!!")
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Println("Could not ping database!!")
		return nil, err
	}

	log.Println("Connected to the Database!!")

	return pool, nil
}
