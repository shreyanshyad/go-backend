package services

import (
	er "backend/dashboard/errors"
	"backend/dashboard/models"
	"backend/dashboard/repository"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Role service to manage roles and permissions to dashboards and views
type RoleService struct {
	r   *repository.RoleRepository
	l   *log.Logger
	rdb *redis.Client
}

// Creates a new instance of RoleService
func NewRoleService(r *repository.RoleRepository, rdb *redis.Client, l *log.Logger) *RoleService {
	return &RoleService{r, l, rdb}
}

// Returns all roles for a user for a dashboard
func (s *RoleService) GetRolesForUserForDashboard(userId, dashId uuid.UUID) ([]*models.Role, error) {
	roles, err := s.r.GetRolesForUserForDashboard(userId, dashId)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		role.Permissions, err = s.r.GetPermissionsForRoleId(role.ID)
		if err != nil {
			return nil, err
		}
	}

	return roles, nil
}

// Returns true if the user has the specified permission for the dashboard
func (s *RoleService) ExistsPermissionForUserForDashboard(userId, dashId uuid.UUID, permName string) (bool, error) {
	return s.r.ExistsPermissionForUserForDashboard(userId, dashId, permName)
}

// Add a user to dashboard with given role
func (s *RoleService) AddUserToDash(dashId, userId uuid.UUID, roleName string) error {
	return s.r.GrantDashLevelRoleToUser(userId, roleName, dashId)
}

func (s *RoleService) AddUserToView(viewId, userId uuid.UUID, roleName string) error {
	return s.r.GrantViewLevelRoleToUser(context.Background(), userId, roleName, viewId)
}

func (s *RoleService) GetAllRoles() ([]*models.Role, error) {
	//can use redis here as roles and their permissions are not changing while running
	//check if redis has this
	res := s.rdb.Get(context.Background(), "all-roles")
	if res.Err() == nil {
		//if it does, return it
		var roles []*models.Role
		err := res.Scan(&roles)
		if err != nil {
			return nil, err
		}
		s.l.Println("Roles fetched from redis")
		return roles, nil
	}

	roles, err := s.r.GetAllRoles()
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		role.Permissions, err = s.r.GetPermissionsForRoleId(role.ID)
		if err != nil {
			return nil, err
		}
	}

	//set it in redis
	err = s.rdb.Set(context.Background(), "all-roles", roles, 0).Err()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *RoleService) ExistsPermissionForUserForView(userId, viewId uuid.UUID, permName string) (bool, error) {
	return s.r.ExistsPermissionForUserForView(userId, viewId, permName)
}

func (s *RoleService) GetRolesForUsersForDashboard(dashId uuid.UUID) ([]*models.Role, error) {
	roles, err := s.r.GetRolesForUsersForDashboard(dashId)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		role.Permissions, err = s.r.GetPermissionsForRoleId(role.ID)
		if err != nil {
			return nil, err
		}
	}
	return roles, nil
}

func (r *RoleService) RevokeDashLevelRoleFromUser(dashId, userId uuid.UUID) error {
	isonly, err := r.r.IsOnlyAdminForDashboard(userId, dashId)
	if err != nil {
		return err
	}
	if isonly {
		return er.ErrCannotRevokeLastAdmin
	}
	return r.r.RevokeDashLevelRoleFromUser(context.Background(), userId, dashId)
}

func (r *RoleService) RevokeViewLevelRoleFromUser(viewId, userId uuid.UUID) error {
	isonly, err := r.r.IsOnlyAdminForView(userId, viewId)
	if err != nil {
		return err
	}
	if isonly {
		return er.ErrCannotRevokeLastAdmin
	}

	return r.r.RevokeViewLevelRoleFromUser(userId, viewId)
}
