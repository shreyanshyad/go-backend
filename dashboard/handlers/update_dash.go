package handlers

import (
	"backend/dashboard/models"
	"backend/middlewares"
	"backend/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/**
 * @api {put} /dashboard/:id Update dashboard
 * @apiName Update dashboard
 * @apiParam {String} id Dashboard ID
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT token
 * @apiBody name Name of dashboard
 * @apiBody description Description of dashboard
 */

func (h *DashHandler) UpdateDash(w http.ResponseWriter, r *http.Request) {
	// get id from mux.Vars
	dashId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// get user id from mux.Vars
	userId, err := middlewares.GetUserIDFromVars(r)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	d := &createDashRequest{}
	err = d.fromJSON(r.Body)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// validating payload
	validate := validator.New()
	err = validate.Struct(d)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// updating dashboard
	newdash := &models.Dash{ID: dashId, Name: d.Name, Description: d.Description}
	err = h.s.UpdateDash(userId, dashId, newdash)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, newdash)
}
