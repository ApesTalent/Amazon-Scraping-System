package api

import (
	"net/http"
	"os"
	"time"

	"primeprice.com/pkg/util"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"primeprice.com/dal"
)

// Login login
func Login(c *gin.Context) {
	var u dal.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid json provided")
		return
	}
	user, err := dal.Login(u)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

// Register register
func Register(c *gin.Context) {
	var u dal.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid json provided")
		return
	}
	err := dal.CreateUser(u)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Username already exists")
		return
	}
	c.JSON(http.StatusOK, "success")
}

// ProtectedUser return current user from token
func ProtectedUser(c *gin.Context) {
	id, err := util.ConvStrToObjID(c.MustGet("id").(string))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	u, err := dal.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, u)
}

// UpdateRole updates user's role
func UpdateRole(c *gin.Context) {
	var u dal.User
	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid json provided")
		return
	}
	err := dal.UpdateRole(u)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, "success")
}

// ChangePassword change user's password
func ChangePassword(c *gin.Context) {
	var p dal.PasswordPayload
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid json provided")
		return
	}
	id, err := util.ConvStrToObjID(c.MustGet("id").(string))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	err = dal.ChangePassword(p, id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, "success")
}

// CreateToken create token based on user data
func CreateToken(u dal.User) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["id"] = u.ID
	atClaims["role"] = u.Role
	atClaims["exp"] = time.Now().Add(time.Minute * 600).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}
