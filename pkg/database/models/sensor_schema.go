package models

type SensorSchema struct {
	ID        int    `gorm:"primaryKey;autoIncrement:false"`
	GroupName string `gorm:"primaryKey"`
	X         float64
	Y         float64
	Z         float64
}
