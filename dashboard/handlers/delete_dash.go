package handlers

import (
	"backend/middlewares"
	"backend/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *DashHandler) DeleteDash(w http.ResponseWriter, r *http.Request) {
	// Get the dashboard id from the request
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusBadRequest, "Invalid dashboard id")
		return
	}

	// Get the user id from the request
	userId, err := middlewares.GetUserIDFromVars(r)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, "Failed to decode id")
		return
	}

	// Delete the dashboard from the database
	err = h.s.DeleteDashById(userId, id)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusForbidden, "Failed to delete dashboard")
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, nil)
}
