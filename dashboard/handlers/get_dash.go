package handlers

import (
	er "backend/dashboard/errors"
	"backend/middlewares"
	"backend/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/**
 * @api {get} /dashboard/:dashboardId Get dashboard by id
 * @apiName Get dashboard by its id you have access to
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT Authorization token
 * @apiParam {String} dashboardId Dashboard ID
 */

func (h *DashHandler) GetDash(w http.ResponseWriter, r *http.Request) {
	//getting dashboard id from url
	vars := mux.Vars(r)
	dashId, err := uuid.Parse(vars["id"])
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	//getting user id from mux.Vars
	userId, err := middlewares.GetUserIDFromVars(r)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	//getting dashboard from database
	dash, err := h.s.GetDashByIdForUser(userId, dashId)
	if err != nil {
		h.l.Printf("Could not get dashboard: %v", err)
		if err == er.ErrNoPerm {
			utils.WriteFailureResponse(w, http.StatusForbidden, err.Error())
		} else {
			utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	//sending dashboard to client
	utils.WriteSuccessResponse(w, http.StatusOK, dash)
	if err != nil {
		h.l.Printf("Could not encode dashboard: %v", err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

/**
 * @api {get} /dashboard Get all dashboards
 * @apiName Get all dashboards you have access to
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT Authorization token
 */

func (h *DashHandler) GetDashs(w http.ResponseWriter, r *http.Request) {
	//getting user id from mux.Vars
	userId, err := middlewares.GetUserIDFromVars(r)
	if err != nil {
		h.l.Println(err)
		utils.WriteFailureResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	//getting dashboards from database
	dashs, err := h.s.GetAllDashboardsForUser(userId)
	if err != nil {
		h.l.Printf("Could not get dashboards: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//sending dashboards to client
	utils.WriteSuccessResponse(w, http.StatusOK, dashs)
}
