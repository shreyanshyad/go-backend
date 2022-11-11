package repository

import (
	"backend/dashboard/db"
	"backend/dashboard/models"
	"context"
	"log"

	"github.com/google/uuid"
)

// Contains methods for interacting with the dashaboard in database
type ViewRepository struct {
	Conn *db.DashboardDb
	L    *log.Logger
}

// Returns a new instance of ViewRepository
func NewViewRepository(conn *db.DashboardDb, l *log.Logger) *ViewRepository {
	return &ViewRepository{conn, l}
}

// Add view attached to a particular dashboard to database
func (repo *ViewRepository) AddView(ctx context.Context, view *models.View, userId uuid.UUID) error {
	tx, err := repo.Conn.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var uuid uuid.UUID
	err = tx.QueryRow("INSERT INTO view (dashboard_id, name, description) VALUES ($1, $2, $3) RETURNING id", view.DashID, view.Name, view.Description).Scan(&uuid)
	if err != nil {
		return err
	}
	view.ID = uuid

	// assign the user as admin on new view
	_, err = tx.Exec("INSERT INTO user_role_view (view_id, user_id, role_id) VALUES ($1, $2, $3)", view.ID, userId, 1)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Get view by id
func (repo *ViewRepository) GetView(viewId, userId uuid.UUID) (*models.View, error) {
	view := &models.View{}
	err := repo.Conn.Conn.QueryRow("SELECT view_id, view_name, view_desc FROM view_pemrs WHERE view_id = $1 AND user_id=$2 AND perm_name='read'", viewId, userId).Scan(&view.ID, &view.Name, &view.Description)
	if err != nil {
		return nil, err
	}
	return view, nil
}

// Get all views attached to a particular dashboard
func (repo *ViewRepository) GetViewsByDashId(dashId uuid.UUID) ([]*models.View, error) {
	rows, err := repo.Conn.Conn.Query("SELECT id, name, description FROM views WHERE dashboard_id = $1", dashId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	views := []*models.View{}
	for rows.Next() {
		view := &models.View{}
		err := rows.Scan(&view.ID, &view.Name, &view.Description)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

// Get all views attached to a particular dashboard for a particular user
func (repo *ViewRepository) GetViewsByDashIdForUser(dashId, userId uuid.UUID) ([]*models.View, error) {
	rows, err := repo.Conn.Conn.Query("SELECT view_id, view_name, view_desc FROM view_perms WHERE dash_id = $1 AND user_id = $2 AND perm_name='read'", dashId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	views := []*models.View{}
	for rows.Next() {
		view := &models.View{}
		err := rows.Scan(&view.ID, &view.Name, &view.Description)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

// Update view content by id
func (repo *ViewRepository) UpdateView(view *models.View) error {
	_, err := repo.Conn.Conn.Exec("UPDATE views SET name = $1, description = $2 WHERE id = $3", view.Name, view.Description, view.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete view by id
func (repo *ViewRepository) DeleteView(id uuid.UUID) error {
	_, err := repo.Conn.Conn.Exec("DELETE FROM views WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
