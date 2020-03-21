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
	LecturerModel struct {
		ID        uuid.UUID
		Name      string
		PhoneNo   string
		Address   string
		Email     string
		IsActive  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}
	LecturerResponse struct {
		ID        uuid.UUID
		Name      string
		PhoneNo   string
		Address   string
		Email     string
		IsActive  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.UUID
		UpdatedAt time.Time
	}
)

func (s LecturerModel) Response() LecturerResponse {
	return LecturerResponse{
		ID:        s.ID,
		Name:      s.Name,
		PhoneNo:   s.PhoneNo,
		Address:   s.Address,
		Email:     s.Email,
		IsActive:  s.IsActive,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}
}

func GetOneLecturer(ctx context.Context, db *sql.DB, lecturerID uuid.UUID) (LecturerModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at,
			email
		FROM lecturer
		WHERE 
			id = $1
	`)

	var lecturer LecturerModel
	err := db.QueryRowContext(ctx, query, lecturerID).Scan(
		&lecturer.ID,
		&lecturer.Name,
		&lecturer.Address,
		&lecturer.PhoneNo,
		&lecturer.IsActive,
		&lecturer.CreatedBy,
		&lecturer.CreatedAt,
		&lecturer.UpdatedBy,
		&lecturer.UpdatedAt,
		&lecturer.Email,
	)

	if err != nil {
		return LecturerModel{}, err
	}

	return lecturer, nil

}

func GetAllLecturer(ctx context.Context, db *sql.DB) ([]LecturerModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at,
			email
		FROM lecturer`)

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var lecturers []LecturerModel
	for rows.Next() {
		var lecturer LecturerModel
		rows.Scan(
			&lecturer.ID,
			&lecturer.Name,
			&lecturer.Address,
			&lecturer.PhoneNo,
			&lecturer.IsActive,
			&lecturer.CreatedBy,
			&lecturer.CreatedAt,
			&lecturer.UpdatedBy,
			&lecturer.UpdatedAt,
			&lecturer.Email,
		)

		lecturers = append(lecturers, lecturer)
	}

	return lecturers, nil

}
