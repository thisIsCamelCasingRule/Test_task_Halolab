package service

import (
	"Test_task_Halolab/cmd/models"
	schema "Test_task_Halolab/pkg/database/models"
	"context"
	"fmt"
	"gorm.io/gorm"
	"runtime"
	"time"
)

func (s *Service) SaveSensor(sensor models.Sensor, groupName string) error {
	sensorSchema := schema.SensorSchema{
		ID:        sensor.ID,
		GroupName: groupName,
		X:         sensor.Coordinates.X,
		Y:         sensor.Coordinates.Y,
		Z:         sensor.Coordinates.Z,
	}

	dbSensor, err := s.DB.SensorRepository.GetSensor(sensor.ID, groupName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if dbSensor.GroupName == "" {
		err = s.DB.SensorRepository.CreateSensor(sensorSchema)
		if err != nil {
			return err
		}
	} else {
		err = s.DB.SensorRepository.UpdateSensor(sensorSchema)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) SaveSensorToData(sensorID int, groupName string) error {
	sensorToData := schema.SensorToDataConnectorSchema{
		SensorID:  sensorID,
		GroupName: groupName,
	}

	id, err := s.DB.SensorToDataRepository.GetSensorDataID(sensorToData)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if id == 0 {
		err = s.DB.SensorToDataRepository.CreateSensorToDataConnector(sensorToData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) SaveSensorGroup(groupList []models.SensorGroup) error {
	for _, group := range groupList {
		for _, sensor := range group.SensorList {
			err := s.SaveSensor(sensor, group.Name)
			if err != nil {
				return err
			}

			err = s.SaveSensorToData(sensor.ID, group.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) SaveSensorData(data models.SensorData) error {
	// find sensor by id and group name and retrieve sensor id from sensor to data connection table
	sensorDataID, err := s.DB.SensorToDataRepository.GetSensorDataID(schema.SensorToDataConnectorSchema{
		SensorID:  data.SensorID,
		GroupName: data.GroupName,
	})
	if err != nil {
		return err
	}

	// save sensor data
	timestamp := time.Now()
	for _, specie := range data.FishMeasurementList {
		preparedData := schema.SensorDataSchema{
			SensorID:     sensorDataID,
			Temperature:  data.Temperature,
			Transparency: data.Transparency,
			FishSpecie:   specie.Specie,
			FishAmount:   specie.Count,
			Timestamp:    timestamp,
		}
		err = s.DB.SensorDataRepository.CreateSensorData(preparedData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) StartDataAggregator(ctx context.Context, pool <-chan models.SensorData) chan error {
	g := runtime.GOMAXPROCS(0)
	errChan := make(chan error)
	for i := 0; i < g; i++ {
		// create some aggregators to receive and save data from pool
		go func(id int) {
			for {
				select {
				case data, open := <-pool:
					if !open {
						return
					}
					// Save data to db
					fmt.Printf("aggregator %d received data: %v \n", id, data)
					err := s.SaveSensorData(data)
					if err != nil {
						errChan <- err
					}
				case <-ctx.Done():
					for data := range pool {
						err := s.SaveSensorData(data)
						if err != nil {
							errChan <- err
						}
					}
					fmt.Printf("aggregator %d stoped! \n", id)
					return
				default:
					continue
				}
			}
		}(i)
	}

	return errChan
}
