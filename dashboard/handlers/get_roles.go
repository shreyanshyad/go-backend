package handlers

import (
	"backend/utils"
	"net/http"
)

/**
 * @api {get} /roles Get all roles
 * @apiName Get all roles
 * @apiGroup Role
 * @apiHeader {String} Authorization JWT Authorization token
 */

func (h *DashHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.s.Rs.GetAllRoles()
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, "Failed to get roles")
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, roles)
}
