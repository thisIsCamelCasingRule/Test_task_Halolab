package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

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

func (a *API) GetGroupAverageTransparency(c *gin.Context) {
	groupName := c.Param("groupName")
	transparency, err := a.Service.GroupAverageTransparency(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, transparency)
}

func (a *API) GetGroupSpecies(c *gin.Context) {
	groupName := c.Param("groupName")
	species, err := a.Service.GroupSpeciesList(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, species)
}

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
