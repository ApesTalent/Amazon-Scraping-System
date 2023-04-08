package api

import (
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"primeprice.com/dal"
	"primeprice.com/pkg/logger"

	"github.com/gin-gonic/gin"
)

// CreateProxy create new Proxy
func CreateProxy(c *gin.Context) {
	var p dal.Proxy
	defer c.Request.Body.Close()

	err := c.BindJSON(&p)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = dal.CreateProxy(p)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "successs")
}

func CreateProxies(c *gin.Context) {
	fmt.Println("Create proxies...")
	var plist dal.Proxy
	var p dal.Proxy
	defer c.Request.Body.Close()

	err := c.BindJSON(&plist)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	list := strings.Split(plist.Proxy, "\n")

	fmt.Println("proxy lists...")
	fmt.Println(list)

	for _, v := range list {
		p.Proxy = v
		err := dal.CreateProxy(p)
		if err != nil {
			logger.Println(err)
		}
	}

	c.JSON(http.StatusOK, "successs")
}

// DeleteProxy delete a Proxy
func DeleteProxy(c *gin.Context) {
	proxyID, err := primitive.ObjectIDFromHex(c.Param("proxy_id"))
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = dal.DeleteProxy(proxyID)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "successs")
}

// DeleteProxy delete a Proxy
func DeleteProxies(c *gin.Context) {

	err := dal.DeleteProxies()
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "successs")
}

// ListProxy get all Proxyes
func ListProxy(c *gin.Context) {
	ps, err := dal.ListProxy()
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, ps)
}
