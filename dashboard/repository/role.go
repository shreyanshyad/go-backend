package repository

import (
	"backend/dashboard/db"
	"backend/dashboard/models"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
)

// Contains methods for access management for dashboards and its views
type RoleRepository struct {
	Conn *db.DashboardDb
	L    *log.Logger
}

// Returns a new RoleRepository instance
func NewRoleRepository(conn *db.DashboardDb, l *log.Logger) *RoleRepository {
	return &RoleRepository{conn, l}
}

// Add a new kind of permission to database
func (r *RoleRepository) AddPermission(perm *models.Permission) error {
	var id int
	err := r.Conn.Conn.QueryRow("INSERT INTO permissions (permission) VALUES ($1) RETURNING id", perm.Name).Scan(&id)
	if err != nil {
		return err
	}
	perm.ID = id
	return nil
}

// Delete a permission with given id from database
func (r *RoleRepository) DeletePermission(id int) error {
	_, err := r.Conn.Conn.Exec("DELETE FROM permissions WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// Get a permission with given id from database
func (r *RoleRepository) GetPermission(id int) (*models.Permission, error) {
	perm := &models.Permission{}
	err := r.Conn.Conn.QueryRow("SELECT id, permission FROM permissions WHERE id = $1", id).Scan(&perm.ID, &perm.Name)
	if err != nil {
		return nil, err
	}
	return perm, nil
}

// Get a permission with given name from database. Names are unique to each permission.
func (r *RoleRepository) GetPermissionByName(name string) (*models.Permission, error) {
	perm := &models.Permission{}
	err := r.Conn.Conn.QueryRow("SELECT id, permission FROM permissions WHERE permission = $1", name).Scan(&perm.ID, &perm.Name)
	if err != nil {
		return nil, err
	}
	return perm, nil
}

// Update permissions details for an id
func (r *RoleRepository) UpdatePermission(perm *models.Permission) error {
	_, err := r.Conn.Conn.Exec("UPDATE permissions SET permission = $1 WHERE id = $2", perm.Name, perm.ID)
	if err != nil {
		return err
	}
	return nil
}

// Get list of all types of permissions from database
func (r *RoleRepository) GetAllPermissions() ([]*models.Permission, error) {
	rows, err := r.Conn.Conn.Query("SELECT id, permission FROM permissions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	perms := []*models.Permission{}
	for rows.Next() {
		perm := &models.Permission{}
		err := rows.Scan(&perm.ID, &perm.Name)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, nil
}

// Add a new role to database. Created role will have no permissions by default.
func (r *RoleRepository) AddRole(role *models.Role) error {
	var id int
	err := r.Conn.Conn.QueryRow("INSERT INTO roles (name) VALUES ($1) RETURNING id", role.Name).Scan(&id)
	if err != nil {
		return err
	}
	role.ID = id
	return nil
}

// Delete a role with given id from database
func (r *RoleRepository) DeleteRole(id int) error {
	_, err := r.Conn.Conn.Exec("DELETE FROM roles WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// Get a role with given id from database
func (r *RoleRepository) GetRole(id int) (*models.Role, error) {
	role := &models.Role{}
	err := r.Conn.Conn.QueryRow("SELECT id, name FROM roles WHERE id = $1", id).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// Get a role with given name from database. Names are unique to each role.
func (r *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	role := &models.Role{}
	err := r.Conn.Conn.QueryRow("SELECT id, name FROM roles WHERE name = $1", name).Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// Update role details for an id
func (r *RoleRepository) UpdateRole(role *models.Role) error {
	_, err := r.Conn.Conn.Exec("UPDATE roles SET name = $1 WHERE id = $2", role.Name, role.ID)
	if err != nil {
		return err
	}
	return nil
}

// Get list of all roles from database
func (r *RoleRepository) GetAllRoles() ([]*models.Role, error) {
	rows, err := r.Conn.Conn.Query("SELECT id, name FROM roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []*models.Role{}
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// Get all permissions for a role with given id
func (r *RoleRepository) GetPermissionsForRoleId(roleId int) ([]*models.Permission, error) {
	rows, err := r.Conn.Conn.Query("SELECT p.id, p.name FROM permissions p, role_has_permissions rhp WHERE rhp.role_id = $1 AND rhp.permission_id = p.id", roleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	perms := []*models.Permission{}
	for rows.Next() {
		perm := &models.Permission{}
		err := rows.Scan(&perm.ID, &perm.Name)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, nil
}

// Add a permissions based on their ids to a role with given id
func (r *RoleRepository) GivePermissionsToRole(ctx context.Context, roleId int, perms []*int) error {
	tx, err := r.Conn.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, perm := range perms {
		_, err := tx.Exec("INSERT INTO role_has_permissions (role_id, permission_id) VALUES ($1, $2)", roleId, perm)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Remove permissions with given ids from a role with given id
func (r *RoleRepository) RevokePermissionsFromRole(ctx context.Context, roleId int, perms []*int) error {
	tx, err := r.Conn.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, perm := range perms {
		_, err := tx.Exec("DELETE FROM role_has_permissions WHERE role_id = $1 AND permission_id = $2", roleId, perm)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Grant role to a user with given id
func (r *RoleRepository) GrantDashLevelRoleToUser(userId uuid.UUID, roleName string, dashId uuid.UUID) error {
	var exists bool
	err := r.Conn.Conn.QueryRow("SELECT EXISTS (SELECT * FROM user_role_dashboard WHERE user_id = $1 AND dashboard_id = $2)", userId, dashId).Scan(&exists)
	if err != nil {
		return err
	}

	var res sql.Result
	if !exists {
		res, err = r.Conn.Conn.Exec("INSERT INTO user_role_dashboard (user_id, dashboard_id, role_id) SELECT $1, $2, roles.id FROM roles WHERE roles.name=$3", userId, dashId, roleName)
		if err != nil {
			return err
		}
	} else {
		res, err = r.Conn.Conn.Exec("UPDATE user_role_dashboard SET role_id = roles.id FROM roles WHERE roles.name=$1 AND user_role_dashboard.user_id=$2 AND user_role_dashboard.dashboard_id=$3;", roleName, userId, dashId)
		if err != nil {
			return err
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("invalid role name")
	}

	return nil
}

// Revoke role from a user with given id
func (r *RoleRepository) RevokeDashLevelRoleFromUser(context context.Context, userId uuid.UUID, dashId uuid.UUID) error {
	tx, err := r.Conn.Conn.BeginTx(context, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec("DELETE FROM user_role_dashboard WHERE user_id = $1 AND dashboard_id = $2", userId, dashId)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.New("user does not have this role")
	}

	//remove user from all views of this dashboard
	_, err = tx.Exec("DELETE FROM user_role_view urv USING view v WHERE urv.user_id = $1 AND urv.view_id = v.id AND v.dashboard_id=$2", userId, dashId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Grant role for a view to a user
func (r *RoleRepository) GrantViewLevelRoleToUser(ctx context.Context, userId uuid.UUID, roleName string, viewId uuid.UUID) error {
	var exists bool
	err := r.Conn.Conn.QueryRow("SELECT EXISTS (SELECT * FROM user_role_view WHERE user_id = $1 AND view_id = $2)", userId, viewId).Scan(&exists)
	if err != nil {
		return err
	}

	var res sql.Result
	if !exists {
		r.L.Println("Inserting role " + roleName + " for user " + userId.String() + " for view " + viewId.String())
		res, err = r.Conn.Conn.Exec("INSERT INTO user_role_view (user_id, view_id, role_id) SELECT $1, $2, roles.id FROM roles WHERE roles.name=$3", userId, viewId, roleName)
		if err != nil {
			return err
		}
	} else {
		r.L.Println("Updating role " + roleName + " for user " + userId.String() + " for view " + viewId.String())
		res, err = r.Conn.Conn.Exec("UPDATE user_role_view SET role_id = roles.id FROM roles WHERE roles.name=$1 AND user_role_view.user_id=$2 AND user_role_view.view_id=$3;", roleName, userId, viewId)
		if err != nil {
			return err
		}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("invalid role name")
	}

	return nil
}

// Revoke role for a view from a user
func (r *RoleRepository) RevokeViewLevelRoleFromUser(userId uuid.UUID, viewId uuid.UUID) error {
	res, err := r.Conn.Conn.Exec("DELETE FROM user_role_view WHERE user_id = $1 AND view_id = $2", userId, viewId)
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.New("user does not have this role")
	}
	return nil
}

// Get all roles for a user for a dashboard
func (r *RoleRepository) GetRolesForUserForDashboard(userId uuid.UUID, dashId uuid.UUID) ([]*models.Role, error) {
	rows, err := r.Conn.Conn.Query("SELECT r.id, r.name FROM roles r JOIN user_role_dashboard urd ON r.id = urd.role_id WHERE urd.user_id = $1 AND urd.dashboard_id = $2", userId, dashId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []*models.Role{}
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) GetRolesForUsersForDashboard(dashId uuid.UUID) ([]*models.Role, error) {
	rows, err := r.Conn.Conn.Query("SELECT r.user_id, r.role_id, r.role_name FROM dashboard_perms r WHERE dash_id=$1 AND perm_name='read'", dashId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []*models.Role{}
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(&role.UserId, &role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) ExistsPermissionForUserForDashboard(userId uuid.UUID, dashID uuid.UUID, permName string) (bool, error) {
	var count bool
	err := r.Conn.Conn.QueryRow("SELECT EXISTS(SELECT * FROM dashboard_perms WHERE user_id=$1 AND dash_id=$2 AND perm_name=$3)", userId, dashID, permName).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *RoleRepository) ExistsPermissionForUserForView(userId uuid.UUID, viewID uuid.UUID, permName string) (bool, error) {
	var count bool
	err := r.Conn.Conn.QueryRow("SELECT EXISTS(SELECT * FROM view_perms WHERE user_id=$1 AND view_Id=$2 AND perm_name=$3)", userId, viewID, permName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count, nil
}

func (r *RoleRepository) IsOnlyAdminForDashboard(userId uuid.UUID, dashID uuid.UUID) (bool, error) {
	var count int
	err := r.Conn.Conn.QueryRow("SELECT COUNT(*) FROM user_role_dashboard WHERE dashboard_id=$1 AND role_id=1", dashID).Scan(&count)
	if err != nil {
		return false, err
	}
	r.L.Println("IsOnlyAdminForDashboard: count: ", count)
	if count == 0 {
		return false, errors.New("no admin role found for dashboard!!! this should not happen")
	}
	if count == 1 {
		var adminID uuid.UUID
		err = r.Conn.Conn.QueryRow("SELECT user_id FROM dashboard_perms WHERE dash_id=$1 AND role_name='admin'", dashID).Scan(&adminID)
		if err != nil {
			return false, err
		}
		if adminID == userId {
			return true, nil
		}
	}
	r.L.Println("User " + userId.String() + " is not the only admin for dashboard " + dashID.String())
	return false, nil
}

func (r *RoleRepository) IsOnlyAdminForView(userId uuid.UUID, viewID uuid.UUID) (bool, error) {
	var count int
	err := r.Conn.Conn.QueryRow("SELECT COUNT(*) FROM user_role_view WHERE view_id=$1 AND role_id=1", viewID).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, errors.New("no admin role found for view!!! this should not happen")
	}
	if count == 1 {
		var adminID uuid.UUID
		err = r.Conn.Conn.QueryRow("SELECT user_id FROM view_perms WHERE view_id=$1 AND role_name='admin'", viewID).Scan(&adminID)
		if err != nil {
			return false, err
		}
		if adminID == userId {
			return true, nil
		}
	}
	return false, nil
}
