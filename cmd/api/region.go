package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetRegionTemperatureMin godoc
// @Summary Retrieves region min temperature
// @Tags region
// @Produce json
// @Param xMin query string false "min X value"
// @Param xMax query string false "max X value"
// @Param yMin query string false "min Y value"
// @Param yMax query string false "max Y value"
// @Param zMin query string false "min Z value"
// @Param zMax query string false "max Z value"
// @Success 200 {object} float64
// @Router /region/temperature/min [get]
func (a *API) GetRegionTemperatureMin(c *gin.Context) {
	xMinString := c.Query("xMin")
	xMin, err := strconv.ParseFloat(xMinString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	yMinString := c.Query("yMin")
	yMin, err := strconv.ParseFloat(yMinString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	zMinString := c.Query("zMin")
	zMin, err := strconv.ParseFloat(zMinString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	xMaxString := c.Query("xMax")
	xMax, err := strconv.ParseFloat(xMaxString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	yMaxString := c.Query("yMax")
	yMax, err := strconv.ParseFloat(yMaxString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	zMaxString := c.Query("zMax")
	zMax, err := strconv.ParseFloat(zMaxString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	temperature, err := a.Service.GetRegionMinTemperature(xMin, xMax, yMin, yMax, zMin, zMax)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, temperature)
}

// GetRegionTemperatureMax godoc
// @Summary Retrieves region max temperature
// @Tags region
// @Produce json
// @Param xMin query string false "min X value"
// @Param xMax query string false "max X value"
// @Param yMin query string false "min Y value"
// @Param yMax query string false "max Y value"
// @Param zMin query string false "min Z value"
// @Param zMax query string false "max Z value"
// @Success 200 {object} float64
// @Router /region/temperature/max [get]
func (a *API) GetRegionTemperatureMax(c *gin.Context) {
	xMinString := c.Query("xMin")
	xMin, err := strconv.ParseFloat(xMinString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	yMinString := c.Query("yMin")
	yMin, err := strconv.ParseFloat(yMinString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	zMinString := c.Query("zMin")
	zMin, err := strconv.ParseFloat(zMinString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	xMaxString := c.Query("xMax")
	xMax, err := strconv.ParseFloat(xMaxString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	yMaxString := c.Query("yMax")
	yMax, err := strconv.ParseFloat(yMaxString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	zMaxString := c.Query("zMax")
	zMax, err := strconv.ParseFloat(zMaxString, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	temperature, err := a.Service.GetRegionMaxTemperature(xMin, xMax, yMin, yMax, zMin, zMax)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, temperature)
}
