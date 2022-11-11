package repository

import (
	"backend/dashboard/db"
	"backend/dashboard/models"
	"context"
	"log"

	"github.com/google/uuid"
)

type DashRepository struct {
	Conn *db.DashboardDb
	L    *log.Logger
}

func NewDashRepository(conn *db.DashboardDb, l *log.Logger) *DashRepository {
	return &DashRepository{conn, l}
}

// Add dashboard to database
func (repo *DashRepository) AddDash(ctx context.Context, dash *models.Dash, userId uuid.UUID) error {
	tx, err := repo.Conn.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var uuid uuid.UUID
	err = tx.QueryRow("INSERT INTO dashboard (name, description) VALUES ($1, $2) RETURNING id", dash.Name, dash.Description).Scan(&uuid)
	if err != nil {
		return err
	}
	dash.ID = uuid

	// assign the user as admin on new dashboard
	_, err = tx.Exec("INSERT INTO user_role_dashboard (dashboard_id, user_id, role_id) VALUES ($1, $2, $3)", dash.ID, userId, 1)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Get dashboard by id
func (repo *DashRepository) GetDash(id uuid.UUID) (*models.Dash, error) {
	dash := &models.Dash{}
	err := repo.Conn.Conn.QueryRow("SELECT id, name, description FROM dashboard WHERE id = $1", id).Scan(&dash.ID, &dash.Name, &dash.Description)
	if err != nil {
		return nil, err
	}
	return dash, nil
}

func (repo *DashRepository) GetAllDashboardsForUser(id uuid.UUID) ([]*models.Dash, error) {
	rows, err := repo.Conn.Conn.Query("SELECT dash_id, dash_name, dash_description FROM dashboard_perms WHERE perm_name='read' AND user_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dashboards := make([]*models.Dash, 0)
	for rows.Next() {
		dash := &models.Dash{}
		err := rows.Scan(&dash.ID, &dash.Name, &dash.Description)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, dash)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return dashboards, nil
}

// Delete dashboard by id
func (repo *DashRepository) DeleteDash(id uuid.UUID) error {
	_, err := repo.Conn.Conn.Exec("DELETE FROM dashboard WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// Update dashboard by id
func (repo *DashRepository) UpdateDash(dash *models.Dash) error {
	_, err := repo.Conn.Conn.Exec("UPDATE dashboard SET name = $1, description = $2 WHERE id = $3", dash.Name, dash.Description, dash.ID)
	if err != nil {
		return err
	}
	return nil
}
