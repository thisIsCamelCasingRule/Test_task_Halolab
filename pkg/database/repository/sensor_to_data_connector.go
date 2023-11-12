package repository

import (
	"Test_task_Halolab/pkg/database/models"
	"gorm.io/gorm"
)

type SensorToDataRepository struct {
	DB *gorm.DB
}

func NewSensorToDataRepository(db *gorm.DB) SensorToDataRepository {
	return SensorToDataRepository{
		DB: db,
	}
}

type ISensorToDataConnector interface {
	GetSensorDataID(sensorToData models.SensorToDataConnectorSchema) (int, error)
	CreateSensorToDataConnector(sensorToData models.SensorToDataConnectorSchema) error
	//UpdateSensorToDataConnector(sensorToData models.SensorToDataConnectorSchema) error
}

func (d *SensorToDataRepository) GetSensorDataID(sensorToData models.SensorToDataConnectorSchema) (int, error) {
	var id int
	err := d.DB.Table("sensor_to_sensor_data_connector").
		Where("sensor_id = ?", sensorToData.SensorID).
		Where("group_name = ?", sensorToData.GroupName).
		Select("id").
		Find(&id).
		Error
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *SensorToDataRepository) CreateSensorToDataConnector(sensorToData models.SensorToDataConnectorSchema) error {
	err := d.DB.Table("sensor_to_sensor_data_connector").
		Save(&sensorToData).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (d *SensorToDataRepository) GetAllSensorToDataGroupID(groupName string) ([]int, error) {
	var idList []int
	err := d.DB.Table("sensor_to_sensor_data_connector").
		Select("id").
		Where("group_name = ?", groupName).
		Find(&idList).
		Error

	if err != nil {
		return nil, err
	}

	return idList, nil
}

//func (d *SensorToDataRepository) UpdateSensorToDataConnector(sensorToData models.SensorToDataConnectorSchema) error {
//	err := d.DB.Table("sensor_to_sensor_data_connector").
//		Where("id = ?", sensorToData.ID).
//		Updates(map[string]interface{}{"sensor_id": sensorToData.SensorID, "group_name": sensor.Y, "z": sensor.Z}).
//		Error
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
