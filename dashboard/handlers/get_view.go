package handlers

import (
	"backend/middlewares"
	"backend/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/**
 * @api {get} /view/:id Get view by id
 * @apiName Get data for a view by its id
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT Authorization token
 * @apiParam {String} id view ID
 */

func (h *DashHandler) GetView(w http.ResponseWriter, r *http.Request) {
	// Get the view id from the request
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusBadRequest, "Invalid view id")
		return
	}

	// Get the user id from the request
	userId, err := middlewares.GetUserIDFromVars(r)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, "Failed to decode id")
		return
	}

	// Get the view from the database
	view, err := h.s.Vs.GetView(id, userId)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusForbidden, "Failed to get view")
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, view)
}
