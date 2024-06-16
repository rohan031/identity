package database

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient() (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_DSN"))
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the redis!!")
	client := redis.NewClient(opt)

	return client, nil
}
