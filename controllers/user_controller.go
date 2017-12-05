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

type token_response struct {
	Token	string	`json:"token"`
}

type error_response struct {
	Error 	string	`json:"error"`
}

// Create a new controller
func NewUserController(database utils.DatabaseAccessor) *UserControllerImpl {
	return &UserControllerImpl{database}
}

// Register routes for UserController
func (uc *UserControllerImpl) Register(router *gin.Engine) {
	router.POST("/login", uc.signin)
	router.POST("/register", uc.signup)
	router.GET("/verify-account", uc.verify)
	g := router.Group("/user")
	{
		g.DELETE("/:id", uc.del)
	}
}

func (uc *UserControllerImpl) signup (ctx *gin.Context) {
	u := models.NewUser()

	if err := ctx.BindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, error_response{err.Error()})
		return
	}

	hash, err := u.HashPassword()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, error_response{err.Error()})
		return
	}

	u.Password = hash
	saveErr := u.Save(uc.database.GetDB())

	if saveErr != nil {
		ctx.JSON(http.StatusBadRequest, error_response{saveErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user created"})
}

func (uc *UserControllerImpl) signin (ctx *gin.Context) {
	var data Signin
	var u = new(models.User)

	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, error_response{err.Error()})
		return
	}

	if u.FindByEmail(data.Email, uc.database.GetDB()); !u.ID.Valid() {
		ctx.JSON(http.StatusNotFound, error_response{"user not found"})
		return
	}

	if !u.CheckPasswordHash(data.Password) {
		ctx.JSON(http.StatusBadRequest, error_response{"invalid credentials"})
		return
	}

	if !u.IsActive {
		ctx.JSON(http.StatusBadRequest, error_response{"user is not active"})
		return
	}

	token, _ := u.GenerateToken()
	ctx.JSON(http.StatusOK, token_response{token})
}

func (uc *UserControllerImpl) del (ctx *gin.Context) {
	id := ctx.Param("id")
	u := new(models.User)

	if !bson.IsObjectIdHex(id) {
		ctx.JSON(http.StatusBadRequest, error_response{"invalid user id"})
		return
	}

	if u.FindByID(id, uc.database.GetDB()); !u.ID.Valid() {
		ctx.JSON(http.StatusNotFound, error_response{"user not found"})
		return
	}

	if err := u.RemoveById(uc.database.GetDB()); err != nil {
		ctx.JSON(http.StatusBadRequest, error_response{err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (uc *UserControllerImpl) verify (ctx *gin.Context) {
	token := ctx.Query("token")
	u := new(models.User)

	if len(token) < 1 {
		ctx.String(http.StatusBadRequest, "unable to verify account")
		return
	}

	if u.FindByToken(token, uc.database.GetDB()); !u.ID.Valid() {
		ctx.String(http.StatusBadRequest, "unable to verify account")
		return
	}

	u.IsActive = true
	u.ValidationToken = ""
	if err:= u.Save(uc.database.GetDB()); err != nil {
		ctx.String(http.StatusBadRequest, "error")
		return
	}

	ctx.String(http.StatusOK, "account verified ")
}