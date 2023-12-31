package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetSensorAverageTemperature godoc
// @Summary Retrieves sensor average temperature
// @Tags sensor
// @Produce json
// @Param codeName path string false "sensor code in format (group id)"
// @Param from query string false "from datetime"
// @Param till query string false "till datetime"
// @Success 200 {object} float64
// @Router /sensor/{codeName}/temperature/average [get]
func (a *API) GetSensorAverageTemperature(c *gin.Context) {
	codeName := c.Param("codeName")
	codeList := strings.Fields(codeName)
	if len(codeList) != 2 {
		c.JSON(http.StatusBadRequest, errors.New("invalid codename"))
		return
	}

	sensorID, err := strconv.Atoi(codeList[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	from := c.Query("from")
	var fromDate time.Time
	if from != "" {
		fromDate, err = time.Parse(time.DateTime, from)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	till := c.Query("till")
	var tillDate time.Time
	if till != "" {
		tillDate, err = time.Parse(time.DateTime, till)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	temperature, err := a.Service.GetSensorAverageTemperature(sensorID, codeList[0], fromDate, tillDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, temperature)
}
