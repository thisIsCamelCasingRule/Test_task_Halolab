package repository

import (
	"Test_task_Halolab/pkg/database/models"
	"gorm.io/gorm"
)

type SensorRepository struct {
	DB *gorm.DB
}

func NewSensorRepository(db *gorm.DB) SensorRepository {
	return SensorRepository{
		DB: db,
	}
}

type ISensor interface {
	CreateSensor(sensor models.SensorSchema) error
	GetSensor(id int, groupName string) (models.SensorSchema, error)
	UpdateSensor(sensor models.SensorSchema) error
	GetSensorInCoordinatesRange(xMin, xMax, yMin, yMax, zMin, zMax float64) ([]models.SensorSchema, error)
}

func (d *SensorRepository) CreateSensor(sensor models.SensorSchema) error {
	err := d.DB.Table("sensor").
		//Exec("INSERT INTO sensor (id, group_name, x, y, z) VALUES (?,?,?,?,?)", sensor.ID, sensor.GroupName, sensor.X, sensor.Y, sensor.Z).
		Save(&sensor).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (d *SensorRepository) GetSensor(id int, groupName string) (models.SensorSchema, error) {
	sensor := models.SensorSchema{}
	err := d.DB.Table("sensor").
		Where("id = ?", id).
		Where("group_name = ?", groupName).
		First(&sensor).
		Error
	if err != nil {
		return sensor, err
	}

	return sensor, nil
}

func (d *SensorRepository) GetSensorInCoordinatesRange(xMin, xMax, yMin, yMax, zMin, zMax float64) ([]models.SensorSchema, error) {
	var sensorList []models.SensorSchema
	err := d.DB.Table("sensor").
		Where("x BETWEEN ? AND ?", xMin, xMax).
		Where("y BETWEEN ? AND ?", yMin, yMax).
		Where("z BETWEEN ? AND ?", zMin, zMax).
		Find(&sensorList).
		Error
	if err != nil {
		return nil, err
	}

	return sensorList, err
}

func (d *SensorRepository) UpdateSensor(sensor models.SensorSchema) error {
	err := d.DB.Table("sensor").
		Where("id = ?", sensor.ID).
		Where("group_name = ?", sensor.GroupName).
		//Exec("INSERT INTO sensor (id, group_name, x, y, z) VALUES (?,?,?,?,?)", sensor.ID, sensor.GroupName, sensor.X, sensor.Y, sensor.Z).
		Updates(map[string]interface{}{"x": sensor.X, "y": sensor.Y, "z": sensor.Z}).
		Error
	if err != nil {
		return err
	}

	return nil
}
