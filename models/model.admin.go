package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"time"
)

type (
	AdminModel struct {
		ID        uuid.UUID
		Username  string
		Password  string
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}
	AdminResponse struct {
		ID        uuid.UUID `json:"id"`
		Username  string    `json:"username"`
		CreatedBy uuid.UUID `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedBy uuid.UUID `json:"updated_by"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func (s AdminModel) Response() AdminResponse {
	return AdminResponse{
		ID:        s.ID,
		Username:  s.Username,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}
}

func GetOneAdminByUsername(ctx context.Context, db *sql.DB, username string) (AdminModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			username,
			password,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM admin
		WHERE 
			username = $1
	`)

	var admin AdminModel
	err := db.QueryRowContext(ctx, query, username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Password,
		&admin.CreatedBy,
		&admin.CreatedAt,
		&admin.UpdatedBy,
		&admin.UpdatedAt,
	)

	if err != nil {
		return AdminModel{}, err
	}

	return admin, nil

}
