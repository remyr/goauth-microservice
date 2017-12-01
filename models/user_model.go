package models

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type User struct {
	ID			bson.ObjectId	`json:"id" bson:"_id"`
	Username	string 			`json:"username omitempty"`
	Email    	string 			`json:"email" binding:"required"`
	Password 	string 			`json:"password" binding:"required"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (u User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u User) HashPassword() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u User) GenerateToken() (string, error) {
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	claims := Claims{
		u.Email,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}

func (u User) Save(db *mgo.Database) error {
	err := u.coll(db).Insert(u)
	return err
}

func (u *User) FindByEmail(email string, db *mgo.Database) error {
	return u.coll(db).Find(bson.M{"email": email}).One(u)
}

func (u User) coll(db *mgo.Database) *mgo.Collection {
	return db.C("users")
}
