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

type addUserToViewRequest struct {
	UserId uuid.UUID `json:"userId" validate:"required,uuid4"`
	Role   string    `json:"role" validate:"required"`
}

func (r *addUserToViewRequest) fromJSON(body io.Reader) error {
	e := json.NewDecoder(body)
	return e.Decode(r)
}

/**
 * @api {post} /view/:viewId/users Add user to view
 * @apiName Upser a user with a role to a view
 * @apiGroup Role
 * @apiHeader {String} Authorization JWT Authorization token
 * @apiParam {String} viewId View ID
 * @apiBody {String} userId User ID for which the role is to be added
 * @apiBody {String} role Role to be added
 */

func (h *DashHandler) AddUserToView(w http.ResponseWriter, r *http.Request) {
	//get viewId from mux.Vars
	viewId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusBadRequest, "Invalid view id")
		return
	}

	// extracting user id
	userId, err := uuid.Parse(mux.Vars(r)[middlewares.KeyUser])
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, "failed to parse user id from jwt.")
		return
	}

	//extracting request body
	req := &addUserToViewRequest{}
	err = req.fromJSON(r.Body)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	can, err := h.s.Rs.ExistsPermissionForUserForView(userId, viewId, perms.ACCESS_MOD)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !can {
		utils.WriteFailureResponse(w, http.StatusUnauthorized, "You are not authorized to add users to this view")
		return
	}

	// adding user to view
	h.l.Println("Adding user ", req.UserId, " to view ", viewId, " with role ", req.Role)
	err = h.s.Rs.AddUserToView(viewId, req.UserId, req.Role)
	if err != nil {
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessResponseMsg(w, http.StatusOK, "User added to view successfully")
}
