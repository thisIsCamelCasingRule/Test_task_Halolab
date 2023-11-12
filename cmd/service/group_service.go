package service

import (
	"Test_task_Halolab/cmd/models"
	schema "Test_task_Halolab/pkg/database/models"
	"context"
	"fmt"
	"sort"
	"time"
)

func (s *Service) GroupAverageTemperature(groupName string) (float64, error) {
	dataList, err := s.GetSensorGroupData(groupName)
	if err != nil {
		return 0, err
	}

	average := 0.0
	for _, data := range dataList {
		average += data.Temperature
	}

	average = average / float64(len(dataList))

	err = s.CacheAverageTemperature(average, groupName)
	if err != nil {
		fmt.Println("Redis error: ", err)
	}

	return average, nil
}

func (s *Service) GroupAverageTransparency(groupName string) (int, error) {
	dataList, err := s.GetSensorGroupData(groupName)
	if err != nil {
		return 0, err
	}

	average := 0
	for _, data := range dataList {
		average += data.Transparency
	}

	average = average / len(dataList)

	err = s.CacheAverageTransparency(average, groupName)
	if err != nil {
		fmt.Println("Redis error: ", err)
	}

	return average, nil
}

func (s *Service) GroupSpeciesList(groupName string) ([]models.FishMeasurement, error) {
	dataList, err := s.GetSensorGroupData(groupName)
	if err != nil {
		return nil, err
	}

	fishList := s.GetFishMeasurementFromSensorData(dataList)

	return fishList, nil
}

func (s *Service) GetTopNGroupSpeciesList(n int, groupName string, fromDate time.Time, toDate time.Time) ([]models.FishMeasurement, error) {
	dataList, err := s.GetSensorGroupData(groupName)
	if err != nil {
		return nil, err
	}

	timeSortedList := s.SortSensorDataByTime(dataList, fromDate, toDate)

	fishList := s.GetFishMeasurementFromSensorData(timeSortedList)

	sort.Sort(models.FishMeasurementList(fishList))

	if n > len(fishList) {
		return fishList, nil
	}

	return fishList[:n], nil
}

func (s *Service) GetRegionMinTemperature(xMin, xMax, yMin, yMax, zMin, zMax float64) (float64, error) {
	sensorList, err := s.DB.SensorRepository.GetSensorInCoordinatesRange(xMin, xMax, yMin, yMax, zMin, zMax)
	if err != nil {
		return 0, err
	}

	sensorDataIDList := make([]int, len(sensorList))
	for i, sensor := range sensorList {
		id, err := s.DB.SensorToDataRepository.GetSensorDataID(schema.SensorToDataConnectorSchema{
			SensorID:  sensor.ID,
			GroupName: sensor.GroupName,
		})
		if err != nil {
			return 0, err
		}

		sensorDataIDList[i] = id
	}

	var dataList []schema.SensorDataSchema
	for _, id := range sensorDataIDList {
		sensorDataList, err := s.DB.SensorDataRepository.GetAllSensorData(id)
		if err != nil {
			return 0, err
		}

		dataList = append(dataList, sensorDataList...)
	}

	if len(dataList) == 0 {
		return 0, nil
	}

	minTemperature := dataList[0].Temperature
	for _, data := range dataList[1:] {
		if data.Temperature < minTemperature {
			minTemperature = data.Temperature
		}
	}

	return minTemperature, nil

}

func (s *Service) GetRegionMaxTemperature(xMin, xMax, yMin, yMax, zMin, zMax float64) (float64, error) {
	sensorList, err := s.DB.SensorRepository.GetSensorInCoordinatesRange(xMin, xMax, yMin, yMax, zMin, zMax)
	if err != nil {
		return 0, err
	}

	sensorDataIDList := make([]int, len(sensorList))
	for i, sensor := range sensorList {
		id, err := s.DB.SensorToDataRepository.GetSensorDataID(schema.SensorToDataConnectorSchema{
			SensorID:  sensor.ID,
			GroupName: sensor.GroupName,
		})
		if err != nil {
			return 0, err
		}

		sensorDataIDList[i] = id
	}

	var dataList []schema.SensorDataSchema
	for _, id := range sensorDataIDList {
		sensorDataList, err := s.DB.SensorDataRepository.GetAllSensorData(id)
		if err != nil {
			return 0, err
		}

		dataList = append(dataList, sensorDataList...)
	}

	if len(dataList) == 0 {
		return 0, nil
	}

	maxTemperature := dataList[0].Temperature
	for _, data := range dataList[1:] {
		if data.Temperature < maxTemperature {
			maxTemperature = data.Temperature
		}
	}

	return maxTemperature, nil

}

func (s *Service) GetFishMeasurementFromSensorData(dataList []schema.SensorDataSchema) []models.FishMeasurement {
	fishMap := make(map[string]int)
	for _, data := range dataList {
		_, ok := fishMap[data.FishSpecie]
		if ok {
			fishMap[data.FishSpecie] += data.FishAmount
		} else {
			fishMap[data.FishSpecie] = data.FishAmount
		}
	}

	fishList := make([]models.FishMeasurement, len(fishMap))
	iter := 0
	for key, value := range fishMap {
		fishList[iter] = models.FishMeasurement{
			Specie: key,
			Count:  value,
		}

		iter += 1
	}

	return fishList
}

func (s *Service) GetSensorAverageTemperature(sensorID int, groupName string, fromDate, toDate time.Time) (float64, error) {
	id, err := s.DB.SensorToDataRepository.GetSensorDataID(schema.SensorToDataConnectorSchema{
		SensorID:  sensorID,
		GroupName: groupName,
	})
	if err != nil {
		return 0, err
	}

	dataList, err := s.DB.SensorDataRepository.GetAllSensorData(id)
	if err != nil {
		return 0, err
	}

	sortedList := s.SortSensorDataByTime(dataList, fromDate, toDate)

	averageTemperature := 0.0
	for _, data := range sortedList {
		averageTemperature += data.Temperature
	}

	averageTemperature = averageTemperature / float64(len(dataList))

	return averageTemperature, nil
}

func (s *Service) GetSensorGroupData(groupName string) ([]schema.SensorDataSchema, error) {
	idList, err := s.DB.SensorToDataRepository.GetAllSensorToDataGroupID(groupName)
	if err != nil {
		return nil, err
	}

	var dataList []schema.SensorDataSchema
	for _, id := range idList {
		sensorDataList, err := s.DB.SensorDataRepository.GetAllSensorData(id)
		if err != nil {
			return nil, err
		}

		dataList = append(dataList, sensorDataList...)
	}

	return dataList, nil
}

func (s *Service) SortSensorDataByTime(dataList []schema.SensorDataSchema, fromDate, toDate time.Time) []schema.SensorDataSchema {
	var timeSortedList []schema.SensorDataSchema
	if !fromDate.IsZero() {
		if !toDate.IsZero() {
			for _, data := range dataList {
				if data.Timestamp.After(fromDate) && data.Timestamp.Before(toDate) {
					timeSortedList = append(timeSortedList, data)
				}
			}
		} else {
			for _, data := range dataList {
				if data.Timestamp.After(fromDate) {
					timeSortedList = append(timeSortedList, data)
				}
			}
		}
	} else if !toDate.IsZero() {
		for _, data := range dataList {
			if data.Timestamp.Before(toDate) {
				timeSortedList = append(timeSortedList, data)
			}
		}
	} else {
		return dataList
	}

	return timeSortedList
}

func (s *Service) CacheAverageTemperature(temperature float64, groupName string) error {
	key := fmt.Sprintf("%s_%s", "temperature", groupName)
	err := s.Redis.Set(context.Background(), key, temperature, time.Second * 10).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CacheAverageTransparency(temperature int, groupName string) error {
	key := fmt.Sprintf("%s_%s", "transparency", groupName)
	err := s.Redis.Set(context.Background(), key, temperature, time.Second * 10).Err()
	if err != nil {
		return err
	}

	return nil
}

