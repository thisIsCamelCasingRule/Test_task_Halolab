package repository

import (
	"Test_task_Halolab/pkg/database/models"
	"gorm.io/gorm"
)

type SensorDataRepository struct {
	DB *gorm.DB
}

func NewSensorDataRepository(db *gorm.DB) SensorDataRepository {
	return SensorDataRepository{
		DB: db,
	}
}

type ISensorData interface {
	CreateSensorData(sensor models.SensorDataSchema) error
	GetSensorData(sensorID int) (models.SensorDataSchema, error)
	GetAllSensorData(sensorID int) ([]models.SensorDataSchema, error)
}

func (d *SensorDataRepository) CreateSensorData(sensor models.SensorDataSchema) error {
	err := d.DB.Table("sensor_data").
		Save(&sensor).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (d *SensorDataRepository) GetSensorData(sensorID int) (models.SensorDataSchema, error) {
	sensor := models.SensorDataSchema{}
	err := d.DB.Table("sensor_data").
		Where("sensor_id = ?", sensorID).
		First(&sensor).
		Error
	if err != nil {
		return sensor, err
	}

	return sensor, nil
}

func (d *SensorDataRepository) GetAllSensorData(sensorID int) ([]models.SensorDataSchema, error) {
	var sensorDataList []models.SensorDataSchema
	err := d.DB.Table("sensor_data").
		Where("sensor_id = ?", sensorID).
		Find(&sensorDataList).
		Error
	if err != nil {
		return nil, err
	}

	return sensorDataList, nil
}

//func (d *SensorDataRepository) GetAllSensorGroupData(sensorID int) ([]models.SensorDataSchema, error) {
//	var sensorGroupDataList []models.SensorDataSchema
//	err := d.DB.Table("sensor").
//		Where("group_name = ?", groupName).
//		Find(&sensorGroupDataList).
//		Error
//	if err != nil {
//		return nil, err
//	}
//
//	return sensorGroupDataList, nil
//}
