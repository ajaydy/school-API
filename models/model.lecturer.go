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
		Password  string
		IsActive  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	LecturerResponse struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		PhoneNo   string    `json:"phone_no"`
		Address   string    `json:"address"`
		Email     string    `json:"email"`
		IsActive  bool      `json:"is_active"`
		CreatedBy uuid.UUID `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedBy uuid.UUID `json:"updated_by"`
		UpdatedAt time.Time `json:"updated_at"`
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
			email,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at,
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
		&lecturer.Email,
		&lecturer.Password,
		&lecturer.IsActive,
		&lecturer.CreatedBy,
		&lecturer.CreatedAt,
		&lecturer.UpdatedBy,
		&lecturer.UpdatedAt,
	)

	if err != nil {
		return LecturerModel{}, err
	}

	return lecturer, nil

}

//func GetAllLecturer(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]LecturerModel, error) {
//
//	var searchQuery string
//
//	if filter.Search != "" {
//		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
//	}
//
//	query := fmt.Sprintf(`
//		SELECT
//			id,
//			name,
//			address,
//			phone_no,
//			is_active,
//			created_by,
//			created_at,
//			updated_by,
//			updated_at,
//			email
//		FROM lecturer
//		%s
//		ORDER BY name  %s
//		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)
//
//	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)
//
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	var lecturers []LecturerModel
//	for rows.Next() {
//		var lecturer LecturerModel
//		rows.Scan(
//			&lecturer.ID,
//			&lecturer.Name,
//			&lecturer.Address,
//			&lecturer.PhoneNo,
//			&lecturer.IsActive,
//			&lecturer.CreatedBy,
//			&lecturer.CreatedAt,
//			&lecturer.UpdatedBy,
//			&lecturer.UpdatedAt,
//			&lecturer.Email,
//		)
//
//		lecturers = append(lecturers, lecturer)
//	}
//
//	return lecturers, nil
//
//}
//
//func (s *LecturerModel) Insert(ctx context.Context, db *sql.DB) error {
//
//	query := fmt.Sprintf(`
//		INSERT INTO lecturer(
//			name,
//			address,
//			email,
//			phone_no,
//			created_by,
//			created_at)
//		VALUES(
//		$1,$2,$3,$4,$5,now())
//		RETURNING id, created_at`)
//
//	err := db.QueryRowContext(ctx, query,
//		s.Name, s.Address, s.Email, s.PhoneNo, s.CreatedBy).Scan(
//		&s.ID, &s.CreatedAt,
//	)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}
//
//func (s *LecturerModel) Update(ctx context.Context, db *sql.DB) error {
//
//	query := fmt.Sprintf(`
//		UPDATE lecturer
//		SET
// 			"name"=$1,
//			address=$2,
//			email=$3,
//			phone_no=$4,
//			updated_at=NOW(),
//			updated_by=$5
//		WHERE id=$6
//		RETURNING id,created_at`)
//
//	err := db.QueryRowContext(ctx, query,
//		s.Name, s.Address, s.Email, s.PhoneNo, s.UpdatedBy, s.ID).Scan(
//		&s.ID, &s.CreatedAt,
//	)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}
