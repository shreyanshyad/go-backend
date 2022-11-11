package handlers

import (
	db "backend/auth/db"
	utils "backend/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

type Auth struct {
	l  *log.Logger
	db *db.UsersDb
}

// NewAuth creates a new auth handler with the given logger
func NewAuth(l *log.Logger, db *db.UsersDb) *Auth {
	return &Auth{l, db}
}

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (req *registerRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(req)
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (req *loginRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(req)
}

/**
 * @api {post} /register Register User
 * @apiName Register
 * @apiGroup Auth
 * @apiBody {String} email Email of the user
 * @apiBody {String} password Password of the user
 * @apiBody {String} username Username of the user
 */

// Register a new user
func (auth *Auth) Register(rw http.ResponseWriter, r *http.Request) {
	auth.l.Println("Handle POST Login")

	req := registerRequest{}
	err := req.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	//validate the request
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		auth.l.Println("Error validating user", err)
		utils.WriteFailureResponse(rw, http.StatusBadRequest, err.Error())
		return
	}

	newUser := &db.User{Username: req.Username, Email: req.Email, Password: req.Password}
	err = auth.db.AddUser(newUser)
	if err != nil {
		auth.l.Println("Error adding user", err)
		utils.WriteFailureResponse(rw, http.StatusInternalServerError, err.Error())
		return
	}

	auth.l.Printf("New user with %s added successfully", newUser.Id.String())

	res := make(map[string]interface{})
	res["id"] = newUser.Id.String()
	res["token"], err = utils.GenerateJWT(newUser.Id.String(), "samudai-dash", "samudai-auth")

	if err != nil {
		auth.l.Println("Error generating token", err)
		utils.WriteFailureResponse(rw, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessResponse(rw, http.StatusOK, res)
}

/**
 * @api {post} /login Login User
 * @apiName Login
 * @apiGroup Auth
 * @apiBody {String} email Email of the user
 * @apiBody {String} password Password of the user
 */

// Login handles the login request
func (auth *Auth) Login(rw http.ResponseWriter, r *http.Request) {
	auth.l.Println("Handle POST Login")

	req := loginRequest{}
	err := req.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	//validate the request
	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		auth.l.Println("Error validating user", err)
		utils.WriteFailureResponse(rw, http.StatusBadRequest, err.Error())
		return
	}

	user, err := auth.db.GetUser(req.Email, req.Password)
	if err == db.ErrNoMatch {
		auth.l.Println("Error getting user", err)
		if err == db.ErrNoMatch {
			utils.WriteFailureResponse(rw, http.StatusUnauthorized, err.Error())
		} else {
			utils.WriteFailureResponse(rw, http.StatusInternalServerError, err.Error())

		}
		return
	}

	auth.l.Printf("User with %d found successfully", user.Id)

	res := make(map[string]interface{})
	res["id"] = user.Id.String()
	res["token"], err = utils.GenerateJWT(user.Id.String(), "samudai-dash", "samudai-auth")

	if err != nil {
		auth.l.Println("Error generating token", err)
		utils.WriteFailureResponse(rw, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessResponse(rw, http.StatusOK, res)
}

func (auth *Auth) Test(rw http.ResponseWriter, r *http.Request) {
	auth.l.Println("Handle Test")

	utils.WriteSuccessResponse(rw, http.StatusOK, "Test")
}
