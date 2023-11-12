package models

import "time"

type SensorDataSchema struct {
	ID           int
	SensorID     int
	Temperature  float64
	Transparency int
	FishSpecie   string
	FishAmount   int
	Timestamp    time.Time
}
