package services

import (
	er "backend/dashboard/errors"
	"backend/dashboard/models"
	"backend/dashboard/repository"
	"context"
	"log"

	"github.com/google/uuid"
)

type ViewService struct {
	vR *repository.ViewRepository
	Rs *RoleService
	l  *log.Logger
}

func NewViewService(r *repository.ViewRepository, rs *RoleService, l *log.Logger) *ViewService {
	return &ViewService{r, rs, l}
}

// Get all views attached to a particular dashboard
func (s *ViewService) GetViewsByDashIdForUser(dashId, userId uuid.UUID) ([]*models.View, error) {
	views, err := s.vR.GetViewsByDashIdForUser(dashId, userId)
	if err != nil {
		return nil, err
	}
	log.Println(dashId, userId)
	log.Println("Views: ", views)
	return views, nil
}

func (s *ViewService) AddView(v *models.View, userId uuid.UUID) error {
	err := s.vR.AddView(context.Background(), v, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *ViewService) GetView(viewId, userId uuid.UUID) (*models.View, error) {
	return s.vR.GetView(viewId, userId)
}

// func (s *ViewService) ExistsPermissionForViewForDashboard(viewId, dashId, userId uuid.UUID, perm string) (bool, error) {
// 	return s.Rs.ExistsPermissionForUserForDashboard(userId, dashId, perm)
// }

func (s *ViewService) DeleteView(viewId, userId uuid.UUID) error {
	can, err := s.Rs.ExistsPermissionForUserForView(userId, viewId, "delete")
	if err != nil {
		return err
	}
	if !can {
		return er.ErrNoPerm
	}
	return s.vR.DeleteView(viewId)
}
