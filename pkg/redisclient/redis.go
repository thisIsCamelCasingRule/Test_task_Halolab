package redisclient

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

func ConnectRedis() (*redis.Client, error) {
	db, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		return nil, err
	}

	cli := redis.NewClient(&redis.Options{
		Addr:	  fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")), // todo: move to env vars
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:		 db ,  // use default DB
	})

	_, err = cli.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return cli, nil
}