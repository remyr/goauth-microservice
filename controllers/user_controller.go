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

type Signin struct {
	Email			string		`json:"email" binding:"required"`
	Password		string		`json:"password" binding:"required"`
}

// Create a new controller
func NewUserController(database utils.DatabaseAccessor) *UserControllerImpl {
	return &UserControllerImpl{database}
}

// Register routes for UserController
func (uc *UserControllerImpl) Register(router *gin.Engine) {
	router.POST("/login", uc.signin)
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

func (uc *UserControllerImpl) signin (ctx *gin.Context) {
	var data Signin
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var u = new(models.User)
	if u.FindByEmail(data.Email, uc.database.GetDB()); !u.ID.Valid() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if !u.CheckPasswordHash(data.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}
	token, _ := u.GenerateToken()
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}