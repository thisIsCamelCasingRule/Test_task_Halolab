package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetGroupAverageTemperature godoc
// @Summary Retrieves group average temperature from last record of each sensor
// @Tags group
// @Produce json
// @Param groupName path string true "Group name"
// @Success 200 {object} float64
// @Router /group/{groupName}/temperature/average [get]
func (a *API) GetGroupAverageTemperature(c *gin.Context) {
	groupName := c.Param("groupName")
	fmt.Println(groupName)
	temperature, err := a.Service.GroupAverageTemperature(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, temperature)
}

// GetGroupAverageTransparency godoc
// @Summary Retrieves group average transparency from last record of each sensor
// @Tags group
// @Produce json
// @Param groupName path string true "Group name"
// @Success 200 {object} int
// @Router /group/{groupName}/transparency/average [get]
func (a *API) GetGroupAverageTransparency(c *gin.Context) {
	groupName := c.Param("groupName")
	transparency, err := a.Service.GroupAverageTransparency(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, transparency)
}

// GetGroupSpecies godoc
// @Summary Retrieves group's fish species list with amount from last record of each sensor
// @Tags group
// @Produce json
// @Param groupName path string true "Group name"
// @Success 200 {array} models.FishMeasurement
// @Router /group/{groupName}/species [get]
func (a *API) GetGroupSpecies(c *gin.Context) {
	groupName := c.Param("groupName")
	species, err := a.Service.GroupSpeciesList(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, species)
}

// GetTopNSpecies godoc
// @Summary Retrieves group's list with amount of top N fish species detected from last record of each sensor
// @Tags group
// @Produce json
// @Param groupName path string true "Group name"
// @Param N path int true "N"
// @Param from query string false "from datetime"
// @Param till query string false "till datetime"
// @Success 200 {array} models.FishMeasurement
// @Router /group/{groupName}/species/top/{N} [get]
func (a *API) GetTopNSpecies(c *gin.Context) {
	groupName := c.Param("groupName")
	nString := c.Param("N")

	n, err := strconv.Atoi(nString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
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

	species, err := a.Service.GetTopNGroupSpeciesList(n, groupName, fromDate, tillDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, species)
}
