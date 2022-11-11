package handlers

import (
	"backend/dashboard/perms"
	"backend/middlewares"
	"backend/utils"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type addUserToDashRequest struct {
	UserId uuid.UUID `json:"userId" validate:"required,uuid4"`
	Role   string    `json:"role" validate:"required"`
}

func (r *addUserToDashRequest) fromJSON(body io.Reader) error {
	e := json.NewDecoder(body)
	return e.Decode(r)
}

/**
 * @api {post} /dashboard/:dashboardId/users Add user to dashboard
 * @apiName Upsert a user with a role to a dashboard
 * @apiGroup Role
 * @apiHeader {String} Authorization JWT Authorization token
 * @apiParam {String} dashboardId Dashboard ID
 * @apiBody {String} userId User ID for which the role is to be added
 * @apiBody {String} role Role to be added
 */

func (h *DashHandler) AddUserToDash(w http.ResponseWriter, r *http.Request) {
	//get dashId from mux.Vars
	dashId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusBadRequest, "Invalid dashboard id")
		return
	}

	// decoding payload to addUserToDashRequest
	req := &addUserToDashRequest{}
	err = req.fromJSON(r.Body)
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

	can, err := h.s.Rs.ExistsPermissionForUserForDashboard(userId, dashId, perms.ACCESS_MOD)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !can {
		utils.WriteFailureResponse(w, http.StatusUnauthorized, "You are not authorized to add users to this dashboard")
		return
	}

	// adding user to dashboard
	err = h.s.Rs.AddUserToDash(dashId, req.UserId, req.Role)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, nil)
}

/**
 * @api {get} /dashboard/:id/users Get users for dashboard
 * @apiName Get data of all memebers of a dashboard. This includes their roles and permissions.
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT Authorization token
 * @apiParam {String} id Dashboard ID
 */

func (h *DashHandler) GetUsersFromDash(w http.ResponseWriter, r *http.Request) {
	//get dashId from mux.Vars
	dashId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusBadRequest, "Invalid dashboard id")
		return
	}

	// extracting user id
	userId, err := uuid.Parse(mux.Vars(r)[middlewares.KeyUser])
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, "failed to parse user id from jwt.")
		return
	}

	can, err := h.s.Rs.ExistsPermissionForUserForDashboard(userId, dashId, perms.READ_PERM)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !can {
		utils.WriteFailureResponse(w, http.StatusUnauthorized, "You are not authorized to view this dashboard")
		return
	}

	// getting users from dashboard
	users, err := h.s.Rs.GetRolesForUsersForDashboard(dashId)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, users)
}
