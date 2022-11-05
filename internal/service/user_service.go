package service

import (
	"fmt"
	"go-mongo-auth/internal/config"
	"go-mongo-auth/internal/database"
	"go-mongo-auth/internal/jwt"
	"go-mongo-auth/internal/utils"
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

	userService struct {
		config      config.Config
		jwt         jwt.IJwtToken
		mongoClient database.IMongoClient
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
	return &userService{
		config:      config,
		jwt:         jwt,
		mongoClient: mongoClient,
	}
}

func (s *userService) Authenticate(login Login) (AuthenticatedResponse, error) {
	login.Password = utils.EncodeString(login.Password)

	var user User
	if err := s.mongoClient.FindOneDocument(_userCollection, bson.M{"email": login.Email, "password": login.Password}).Decode(&user); err != nil {
		log.Println(err)
		return AuthenticatedResponse{}, err
	}

	ss, err := s.jwt.CreateToken(user, s.config.Jwt.Auth.Expiry)
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

func (s *userService) Register(user User) (RegisteredResponse, error) {
	// Default attribute values
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Check for existing users by email
	if s.mongoClient.FindOneDocument(_userCollection, bson.M{"email": user.Email}).Err() == nil {
		return RegisteredResponse{}, fmt.Errorf("existing document found with email: '%v'", user.Email)
	}

	user.Password = utils.EncodeString(user.Password)

	res, err := s.mongoClient.CreateOneDocument(_userCollection, user)
	if err != nil {
		return RegisteredResponse{}, err
	}

	return RegisteredResponse{res.InsertedID}, nil
}
