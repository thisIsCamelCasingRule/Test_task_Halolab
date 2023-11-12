package service

import (
	"Test_task_Halolab/pkg/database"
	"Test_task_Halolab/pkg/redisclient"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	DB *database.Database
	Redis *redis.Client
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) InitService(db *database.Database) error {
	s.DB = db
	cli, err := redisclient.ConnectRedis()
	if err != nil {
		return err
	}

	s.Redis = cli

	return nil
}
