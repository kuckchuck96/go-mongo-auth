package service

import (
	"fmt"
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/jwt"
	"go-mongo-auth/internal/util"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Login struct {
	Email    string `form:"email" json:"email" binding:"email,required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" binding:"required"`
	Email     string             `json:"email" binding:"email,required"`
	Password  string             `json:"password,omitempty" binding:"required"`
	Verified  bool               `json:"verified" bson:"verified" default:"false"`
	CreatedAt time.Time          `json:"created" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated" bson:"updatedAt"`
}

const (
	userCollection = "user"
)

func (user *User) New() {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func Authenticate(login Login) (map[string]any, error) {
	login.Password = util.EncodeString(login.Password)

	res := database.FindOneDocument(userCollection, bson.M{"email": login.Email, "password": login.Password})

	var user User
	if err := res.Decode(&user); err != nil {
		log.Println(err)
		return nil, err
	}

	ss, err := jwt.CreateToken(user, config.GetChrono("jwt.auth.expiry"))
	if err != nil {
		log.Println("Error creating JWT.", err)
		return nil, err
	}

	user.Password = ""

	return map[string]any{
		"user":  user,
		"token": ss,
	}, nil
}

func Register(user User) (*mongo.InsertOneResult, error) {
	// Check for existing users by email
	exists := database.FindOneDocument(userCollection, bson.M{"email": user.Email})

	if exists.Err() == nil {
		return nil, fmt.Errorf("existing document found with email: '%v'", user.Email)
	}

	user.Password = util.EncodeString(user.Password)

	res, err := database.CreateOneDocument(userCollection, user)
	if err != nil {
		return nil, err
	}

	return res, nil
}
