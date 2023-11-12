package database

import (
	"Test_task_Halolab/pkg/database/repository"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var connectionStr string

type Database struct {
	SensorRepository       repository.SensorRepository
	SensorDataRepository   repository.SensorDataRepository
	SensorToDataRepository repository.SensorToDataRepository
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) InitDatabase() error {
	connectionStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	sqlDB, err := sql.Open("pgx", connectionStr)
	if err != nil {
		return errors.New(connectionStr)
	}

	//connectionStr = "user=postgres password=postgres port=5432 dbname=postgres sslmode=disable"
	db, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return errors.New(connectionStr)
	}

	d.SensorRepository = repository.NewSensorRepository(db)
	d.SensorToDataRepository = repository.NewSensorToDataRepository(db)
	d.SensorDataRepository = repository.NewSensorDataRepository(db)

	return nil
}
