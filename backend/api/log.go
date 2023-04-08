package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"primeprice.com/pkg/logger"
)

type Logs struct {
	Logs []string `json:"logs"`
}

// GetLogs gets application logs
func GetLogs(c *gin.Context) {
	content, err := ioutil.ReadFile(logger.LogFile)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	lines := strings.Split(string(content), "\n")

	// reverse lines
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	if len(lines) > 100 {
		lines = lines[:100]
	}

	c.JSON(http.StatusOK, Logs{lines})
}
