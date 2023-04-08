package api

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"primeprice.com/controller"

	"primeprice.com/dal"
	"primeprice.com/pkg/logger"
	"primeprice.com/pkg/util"

	"github.com/gin-gonic/gin"
)

// CreateSearch create new search
func CreateSearch(c *gin.Context) {
	var s dal.Search
	defer c.Request.Body.Close()

	err := c.BindJSON(&s)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id, err := util.ConvStrToObjID(c.MustGet("id").(string))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	s.UserID = id
	s.Status = 1

	ns, err := dal.CreateSearch(s)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// start scraping for new created search
	go func(search dal.Search) {
		controller.StartRecursiveScraping(search.ID)
	}(ns)

	c.JSON(http.StatusOK, "successs")
}

// UpdateSearch update existing search
func UpdateSearch(c *gin.Context) {
	var s dal.Search
	defer c.Request.Body.Close()

	err := c.BindJSON(&s)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = dal.UpdateSearch(s)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "successs")
}

// DeleteSearch delete a search
func DeleteSearch(c *gin.Context) {
	searchID, err := primitive.ObjectIDFromHex(c.Param("search_id"))
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = dal.DeleteSearch(searchID)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, "successs")
}

// ListSearch get all searches
func ListSearch(c *gin.Context) {
	ss, err := dal.ListSearch()
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, ss)
}

// GetSearch get a search
func GetSearch(c *gin.Context) {
	searchID, err := primitive.ObjectIDFromHex(c.Param("search_id"))
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	s, err := dal.GetSearchByID(searchID)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, s)
}
