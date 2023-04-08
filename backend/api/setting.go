package api

import (
	"net/http"
	"strconv"

	"primeprice.com/controller"
	"primeprice.com/dal"
	"primeprice.com/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateSetting update existing config
func UpdateSetting(c *gin.Context) {
	var s dal.AppConfig
	defer c.Request.Body.Close()

	err := c.BindJSON(&s)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = dal.UpdateSetting(s)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "successs")
}

// GetAppConfig gets existing config
func GetAppConfig(c *gin.Context) {
	s, err := dal.GetAppConfig()
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, s)
}

// RunScraping change scraping running status
func RunScraping(c *gin.Context) {
	rs, err := strconv.Atoi(c.Param("run_status"))
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = dal.UpdateRunStatus(rs)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if rs == 1 {
		controller.StartsAllScraping()
	}
	c.JSON(http.StatusOK, "successs")
}
