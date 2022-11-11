package handlers

import (
	"backend/dashboard/models"
	"backend/middlewares"
	"backend/utils"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/**
 * @api {post} /view Create view
 * @apiName Create view
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT Authorization token
 * @apiBody {String} name Name of the view
 * @apiBody {uuid} dashboardId of the dashboard to which the view belongs
 * @apiBody {String} [description] Description of the view
 */

type createViewRequest struct {
	DashboardId uuid.UUID `json:"dashboardId" validate:"required,uuid4"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
}

func (req *createViewRequest) fromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(req)
}

func (h *DashHandler) CreateView(w http.ResponseWriter, r *http.Request) {
	// decoding payload to createViewRequest
	v := &createViewRequest{}
	err := v.fromJSON(r.Body)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// extracting user id
	userId, err := uuid.Parse(mux.Vars(r)[middlewares.KeyUser])
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, "failed to parse user id from jwt.")
		return
	}

	// creating view
	newview := &models.View{DashID: v.DashboardId, Name: v.Name, Description: v.Description}
	err = h.s.Vs.AddView(newview, userId)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, newview)
}
