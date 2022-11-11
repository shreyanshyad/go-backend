package services

import (
	er "backend/dashboard/errors"
	"backend/dashboard/models"
	"backend/dashboard/perms"
	"backend/dashboard/repository"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// service to manage dashboard, its views and access to them
type DashService struct {
	ds  *repository.DashRepository
	Vs  *ViewService
	Rs  *RoleService
	rdb *redis.Client
	L   *log.Logger
}

// Creates a new instance of DashService
func NewDashService(r *repository.DashRepository, v *ViewService, rr *RoleService, rdb *redis.Client, l *log.Logger) *DashService {
	return &DashService{r, v, rr, rdb, l}
}

// Create a new dashboard witb new ID updated in the model
func (s *DashService) AddDash(dash *models.Dash, userId uuid.UUID) error {
	return s.ds.AddDash(context.Background(), dash, userId)
}

// Get a dashboard with given id
func (s *DashService) GetDashByIdForUser(userId, dashId uuid.UUID) (*models.Dash, error) {
	can, err := s.Rs.ExistsPermissionForUserForDashboard(userId, dashId, perms.READ_PERM)
	s.L.Println("Permission: ", can, "for user: ", userId, "for dashboard: ", dashId)
	if err != nil {
		return nil, err
	}
	if !can {
		return nil, er.ErrNoPerm
	}
	dash, err := s.ds.GetDash(dashId)
	if err != nil {
		return nil, err
	}
	dash.Views, err = s.Vs.GetViewsByDashIdForUser(dashId, userId)
	if err != nil {
		return nil, err
	}
	return dash, err
}

// Get all dashboards for a user
func (s *DashService) GetAllDashboardsForUser(userId uuid.UUID) ([]*models.Dash, error) {
	dashs, err := s.ds.GetAllDashboardsForUser(userId)
	if err != nil {
		return nil, err
	}
	for _, dash := range dashs {
		dash.Views, err = s.Vs.vR.GetViewsByDashIdForUser(dash.ID, userId)
		if err != nil {
			return nil, err
		}
	}
	return dashs, nil
}

// Delete a dashboard with given id
func (s *DashService) DeleteDashById(userId, dashId uuid.UUID) error {
	can, err := s.Rs.ExistsPermissionForUserForDashboard(userId, dashId, perms.DELETE_PERM)
	if err != nil {
		return err
	}
	if !can {
		return er.ErrNoPerm
	}

	return s.ds.DeleteDash(dashId)
}

func (s *DashService) UpdateDash(userId, dashId uuid.UUID, dash *models.Dash) error {
	can, err := s.Rs.ExistsPermissionForUserForDashboard(userId, dashId, perms.WRITE_PERM)
	if err != nil {
		return err
	}
	if !can {
		return er.ErrNoPerm
	}

	dash.ID = dashId
	return s.ds.UpdateDash(dash)
}
