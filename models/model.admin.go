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

	AdminUpdatePasswordResponse struct {
		Message string `json:"message"`
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

func GetOneAdmin(ctx context.Context, db *sql.DB, adminID uuid.UUID) (AdminModel, error) {

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
			id = $1
	`)

	var admin AdminModel
	err := db.QueryRowContext(ctx, query, adminID).Scan(
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

func (s *AdminModel) PasswordUpdate(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE admin
		SET
			password = $1,
			updated_at=NOW(),
			updated_by=$2
		WHERE id=$3
		RETURNING id,created_at,updated_at,created_by`)

	err := db.QueryRowContext(ctx, query,
		s.Password, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}
