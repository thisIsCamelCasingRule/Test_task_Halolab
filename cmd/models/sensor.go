package models

type Sensor struct {
	ID             int
	Coordinates    Coordinates
	DataOutputRate int
}

type SensorGroup struct {
	Name       string
	SensorList []Sensor
}

func NewGroup(groupName string, sensorList []Sensor) SensorGroup {
	return SensorGroup{
		Name:       groupName,
		SensorList: sensorList,
	}
}

func NewSensor(id int, coordinates Coordinates, rate int) Sensor {
	return Sensor{
		ID:             id,
		Coordinates:    coordinates,
		DataOutputRate: rate,
	}
}

type Coordinates struct {
	X float64
	Y float64
	Z float64
}

type SensorData struct {
	SensorID            int
	GroupName           string
	Temperature         float64
	Transparency        int
	FishMeasurementList []FishMeasurement
}

func NewSensorData(sensorID int, groupName string, tmprtre float64, trsprncy int, measurementList []FishMeasurement) SensorData {
	return SensorData{
		SensorID:            sensorID,
		GroupName:           groupName,
		Temperature:         tmprtre,
		Transparency:        trsprncy,
		FishMeasurementList: measurementList,
	}
}

type FishMeasurement struct {
	Specie string
	Count  int
}

type FishMeasurementList []FishMeasurement

func (a FishMeasurementList) Len() int           { return len(a) }
func (a FishMeasurementList) Less(i, j int) bool { return a[i].Count < a[j].Count }
func (a FishMeasurementList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
