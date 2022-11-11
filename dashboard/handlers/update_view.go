package handlers

import (
	"backend/utils"
	"net/http"
)

/**
 * @api {put} /view/:id Update view
 * @apiName Update view
 * @apiParam {String} id Dashboard ID
 * @apiGroup Dashboard
 * @apiHeader {String} Authorization JWT token
 * @apiBody name Name of view
 * @apiBody description Description of view
 */

func (h *DashHandler) UpdateView(w http.ResponseWriter, r *http.Request) {
	utils.WriteFailureResponse(w, http.StatusNotImplemented, "Not implemented")
}
