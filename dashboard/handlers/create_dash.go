package handlers

import (
	"backend/dashboard/models"
	"backend/middlewares"
	"backend/utils"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type createDashRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (req *createDashRequest) fromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(req)
}

/**
 * @api {post} /dashboard Create dashboard
 * @apiName Create dashboard
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT Authorization token
 * @apiBody {String} name Name of the dashboard
 * @apiBody {String} [description] Description of the dashboard
 */

// Handle request to create a new dashboard
func (dash *DashHandler) CreateDash(rw http.ResponseWriter, r *http.Request) {
	// decoding payload to createDashRequest
	d := &createDashRequest{}
	err := d.fromJSON(r.Body)
	if err != nil {
		utils.WriteFailureResponse(rw, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// validating payload
	validate := validator.New()
	err = validate.Struct(d)
	if err != nil {
		utils.WriteFailureResponse(rw, http.StatusBadRequest, err.Error())
		return
	}

	// extracting user id
	userId, err := uuid.Parse(mux.Vars(r)[middlewares.KeyUser])
	if err != nil {
		utils.WriteFailureResponse(rw, http.StatusInternalServerError, "failed to parse user id from jwt.")
		return
	}

	// creating dashboard
	newdash := &models.Dash{Name: d.Name, Description: d.Description}
	dash.l.Println("Creating new dashboard with name: ", newdash.Name, " and description: ", newdash.Description)
	err = dash.s.AddDash(newdash, userId)
	if err != nil {
		utils.WriteFailureResponse(rw, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessResponse(rw, http.StatusCreated, newdash)
}
