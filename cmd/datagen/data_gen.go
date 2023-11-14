package datagen

import (
	"Test_task_Halolab/cmd/models"
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	minOceanTemperature float64 = 0.0
	maxOceanTemperature float64 = 30.0

	maxZ float64 = 11.0
	minZ float64 = 0.0

	maxY int = 1000
	minY int = -1000

	maxX int = 2000.0
	minX int = 0.0

	minRadius int = 50
	maxRadius int = 200

	maxRate int = 10
	minRate int = 1

	maxFishAmount int = 1000

	maxTransparency int = 200
	minTransparency int = 10
)

var s float64
var b float64
var d float64
var groupNameList = []string{"alpha", "beta", "gamma", "theta"} // move to yaml
var fishSpecies = []string{"Atlantic Blue Tuna", "Atlantic Cod", "Blue Marlin", "Coelacanth", "Yellow Tuna"}
var maxSensorsEachGroup = 5

func init() {
	rand.Seed(time.Now().UnixNano())
	s = maxOceanTemperature / float64(maxY)
	b = maxOceanTemperature
	d = 0.9
}

type DataGen struct {
	WaitGroup *sync.WaitGroup
	Groups    []models.SensorGroup
	Ctx       context.Context
	Pool      chan models.SensorData
}

func NewDataGen(ctx context.Context) DataGen {
	return DataGen{
		WaitGroup: &sync.WaitGroup{},
		Groups:    generateGroups(),
		Ctx:       ctx,
		Pool:      startGroupPool(),
	}
}

func generateGroups() []models.SensorGroup {
	groupsList := make([]models.SensorGroup, len(groupNameList))

	for j, groupName := range groupNameList {
		sensorsAmount := rand.Intn(maxSensorsEachGroup-1) + 1
		sensorList := make([]models.Sensor, sensorsAmount)
		x := rand.Intn(maxX-maxRadius) + maxRadius // generate x in range [maxRadius; maxX - maxRadius), represents group placement on x axis
		y := rand.Intn(maxY-maxRadius) + maxRadius // generate y in range [maxRadius; maxY - maxRadius), represents group placement on y axis
		sensorYSign := 1.0                         // defines whether group placed upper or lower on y coordinate
		if rand.Float64() < 0.5 {
			sensorYSign = -1.0
		}
		radius := rand.Intn(maxRadius-minRadius) + minRadius // generate radius for group
		for i := 0; i < sensorsAmount; i++ {
			sensorCoordinates := generateSensorCoordinates(x, y, radius, sensorYSign)
			sensor := models.NewSensor(i, sensorCoordinates, rand.Intn(maxRate-minRate)+minRate)
			sensorList[i] = sensor
		}

		group := models.NewGroup(groupName, sensorList)

		groupsList[j] = group
	}

	return groupsList
}

func startSensor(ctx context.Context, sensor models.Sensor, groupName string, poolChan chan models.SensorData, wg *sync.WaitGroup) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("sensor stopped!")
				wg.Done()
				return
			default:
				// generate Sensor data here and push it to result channel
				tmprtre := generateTemperature(sensor.Coordinates)
				fishList := generateFishData()
				trsprncy := rand.Intn(maxTransparency-minTransparency) + minTransparency

				poolChan <- models.NewSensorData(sensor.ID, groupName, tmprtre, trsprncy, fishList) //push created data to result channel

				time.Sleep(time.Duration(sensor.DataOutputRate * int(time.Second)))
			}
		}
	}()
}

func (d *DataGen) StartSensorGroups() {
	for _, sensorGroup := range d.Groups {
		d.WaitGroup.Add(len(sensorGroup.SensorList))
		for _, sensor := range sensorGroup.SensorList {
			startSensor(d.Ctx, sensor, sensorGroup.Name, d.Pool, d.WaitGroup)
		}
	}
}

func startGroupPool() chan models.SensorData {
	return make(chan models.SensorData, 10) // todo: reassign capacity
}

func generateFishData() []models.FishMeasurement {
	arrayLen := len(fishSpecies)
	fishSpeciesAmount := rand.Intn(arrayLen-1) + 1
	fishMeasurementList := make([]models.FishMeasurement, fishSpeciesAmount)
	generated := make(map[int]bool, fishSpeciesAmount)
	for i := 0; i < fishSpeciesAmount; i++ {
		j := rand.Intn(arrayLen)
		fishMeasurement := models.FishMeasurement{}
		if !generated[j] {
			fishMeasurement.Specie = fishSpecies[j]
			fishMeasurement.Count = rand.Intn(maxFishAmount)
			generated[j] = true
		} else {
			i -= 1 // if already generated specie than
			continue
		}
		fishMeasurementList[i] = fishMeasurement
	}

	return fishMeasurementList
}

func generateTemperature(coordinates models.Coordinates) float64 {
	surfaceTemperature := -1*math.Abs(coordinates.Y*s) + b
	k := math.Log2(surfaceTemperature/(1-math.Pow(d, maxZ))) / math.Log2(d)
	g := -1.0 * (surfaceTemperature * math.Pow(d, maxZ)) / (1 - math.Pow(d, maxZ))
	actualTemperature := math.Pow(d, coordinates.Z+k) + g

	sign := 1.0
	if rand.Float64() < 0.5 {
		sign = -1.0
	}

	deviation := rand.Float64() * 5.0

	if deviation > actualTemperature && sign == -1.0 {
		return actualTemperature
	}

	return actualTemperature + sign*deviation
}

func generateSensorCoordinates(x int, y int, radius int, sign float64) models.Coordinates {
	sensorX := float64(rand.Intn(2*radius) + x - radius)
	sensorY := sign * float64(rand.Intn(2*radius)+y-radius)
	sensorZ := float64(rand.Intn(int(maxZ)))

	return models.Coordinates{
		X: sensorX,
		Y: sensorY,
		Z: sensorZ,
	}
}
