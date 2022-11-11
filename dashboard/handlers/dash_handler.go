package handlers

import (
	"backend/dashboard/services"
	"log"
)

type DashHandler struct {
	l *log.Logger
	s *services.DashService
}

// NewDash creates a new dashboard handler with the given logger and repository
func NewDash(l *log.Logger, s *services.DashService) *DashHandler {
	return &DashHandler{l, s}
}
