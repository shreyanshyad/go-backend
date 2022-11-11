package middlewares

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	KeyUser = "user-key"
)

func GetUserIDFromVars(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(mux.Vars(r)[KeyUser])
}
