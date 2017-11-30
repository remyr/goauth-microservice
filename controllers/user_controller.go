package controllers

import (
	"github.com/remyr/goauth-microservice/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/remyr/goauth-microservice/models"
)

type UserControllerImpl struct {
	database utils.DatabaseAccessor
}

// Create a new controller
func NewUserController(database utils.DatabaseAccessor) *UserControllerImpl {
	return &UserControllerImpl{database}
}

// Register routes for UserController
func (uc *UserControllerImpl) Register(router *gin.Engine) {
	router.POST("/login")
	router.POST("/register", uc.signup)
}

func (uc *UserControllerImpl) signup (ctx *gin.Context) {
	var u models.User

	if err := ctx.BindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := u.HashPassword()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.Password = hash
	u.ID = bson.NewObjectId()

	saveErr := u.Save(uc.database.GetDB())
	if saveErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": saveErr.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user created"})
}