package handlers

import (
	er "backend/dashboard/errors"
	"backend/dashboard/perms"
	"backend/middlewares"
	"backend/utils"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type deleteUserFromViewRequest struct {
	UserID uuid.UUID `json:"userId" validate:"required"`
}

func (d *deleteUserFromViewRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

type DeleteUserFromDash struct {
	UserID uuid.UUID `json:"userId" validate:"required"`
}

func (d *DeleteUserFromDash) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

/**
 * @api {delete} /view/:id/users Remove user from view
 * @apiName Remove user from virew
 * @apiParam {String} id View ID
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT token
 * @apiBody userId UUID of user to remove
 */

func (h *DashHandler) DeleteUserFromView(w http.ResponseWriter, r *http.Request) {
	viewId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := middlewares.GetUserIDFromVars(r)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	d := &deleteUserFromViewRequest{}
	err = d.FromJSON(r.Body)
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

	can, err := h.s.Rs.ExistsPermissionForUserForView(userId, viewId, perms.ACCESS_MOD)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !can {
		utils.WriteFailureResponse(w, http.StatusForbidden, "You don't have permission to remove users from this view")
		return
	}

	err = h.s.Rs.RevokeViewLevelRoleFromUser(viewId, d.UserID)
	if err != nil {
		h.l.Println(err)
		if err == er.ErrCannotRevokeLastAdmin {
			utils.WriteFailureResponse(w, http.StatusBadRequest, err.Error())
		} else {
			utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteSuccessResponseMsg(w, http.StatusOK, "User removed from view")
}

/**
 * @api {delete} /dashboard/:id/users Remove user from dashboard
 * @apiName Remove user from dashboard
 * @apiParam {String} id Dashboard ID
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT token
 * @apiBody userId UUID of user to remove
 */

func (h *DashHandler) DeleteUserFromDash(w http.ResponseWriter, r *http.Request) {
	dashId, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := middlewares.GetUserIDFromVars(r)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	d := &DeleteUserFromDash{}
	err = d.FromJSON(r.Body)
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

	can, err := h.s.Rs.ExistsPermissionForUserForDashboard(userId, dashId, perms.ACCESS_MOD)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !can {
		utils.WriteFailureResponse(w, http.StatusForbidden, "You don't have permission to remove users from this dashboard")
		return
	}

	err = h.s.Rs.RevokeDashLevelRoleFromUser(dashId, d.UserID)
	if err != nil {
		h.l.Println(err)
		if err == er.ErrCannotRevokeLastAdmin {
			utils.WriteFailureResponse(w, http.StatusBadRequest, err.Error())
		} else {
			utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteSuccessResponseMsg(w, http.StatusOK, "User removed from dashboard")
}
