package api

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"primeprice.com/dal"
	"primeprice.com/pkg/logger"

	"github.com/gin-gonic/gin"
)

// ListProduct get all searches
func ListProduct(c *gin.Context) {
	searchID, err := primitive.ObjectIDFromHex(c.Param("search_id"))
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ps, err := dal.ListProduct(searchID)
	if err != nil {
		logger.Println(err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, ps)
}
