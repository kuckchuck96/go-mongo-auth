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
)

type (
	IUserService interface {
		Authenticate(Login) (AuthenticatedResponse, error)
		Register(User) (RegisteredResponse, error)
	}

	UserService struct {
		Config      config.Config
		Jwt         jwt.IJwtToken
		MongoClient database.IMongoClient
	}

	Login struct {
		Email    string `form:"email" json:"email" binding:"email,required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	User struct {
		Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name      string             `json:"name" binding:"required"`
		Email     string             `json:"email" binding:"email,required"`
		Password  string             `json:"password,omitempty" binding:"required"`
		Verified  bool               `json:"verified" bson:"verified" default:"false"`
		CreatedAt time.Time          `json:"created" bson:"createdAt"`
		UpdatedAt time.Time          `json:"updated" bson:"updatedAt"`
	}

	AuthenticatedResponse struct {
		User  User   `json:"user"`
		Token string `json:"token"`
	}

	RegisteredResponse struct {
		Id any `json:"id"`
	}

	UserErrResponse struct {
		Error string `json:"error"`
	}
)

const _userCollection = "user"

func NewUserService(config config.Config, jwt jwt.IJwtToken, mongoClient database.IMongoClient) IUserService {
	return &UserService{
		Config:      config,
		Jwt:         jwt,
		MongoClient: mongoClient,
	}
}

func (s *UserService) Authenticate(login Login) (AuthenticatedResponse, error) {
	login.Password = util.EncodeString(login.Password)

	res := s.MongoClient.FindOneDocument(_userCollection, bson.M{"email": login.Email, "password": login.Password})

	var user User
	if err := res.Decode(&user); err != nil {
		log.Println(err)
		return AuthenticatedResponse{}, err
	}

	ss, err := s.Jwt.CreateToken(user, s.Config.Jwt.Auth.Expiry)
	if err != nil {
		log.Println("Error creating JWT.", err)
		return AuthenticatedResponse{}, err
	}

	user.Password = ""

	return AuthenticatedResponse{
		user,
		ss,
	}, nil
}

func (s *UserService) Register(user User) (RegisteredResponse, error) {
	// Default attribute values
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Check for existing users by email
	exists := s.MongoClient.FindOneDocument(_userCollection, bson.M{"email": user.Email})

	if exists.Err() == nil {
		return RegisteredResponse{}, fmt.Errorf("existing document found with email: '%v'", user.Email)
	}

	user.Password = util.EncodeString(user.Password)

	res, err := s.MongoClient.CreateOneDocument(_userCollection, user)
	if err != nil {
		return RegisteredResponse{}, err
	}

	return RegisteredResponse{res.InsertedID}, nil
}
