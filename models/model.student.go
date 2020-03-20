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
	StudentModel struct {
		ID          uuid.UUID
		Name        string
		Address     string
		DateOfBirth time.Time
		Gender      int
		Email       string
		PhoneNo     string
		IsActive    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	StudentResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Address     string    `json:"address"`
		DateOfBirth time.Time `json:"date_of_birth"`
		Gender      int       `json:"gender"`
		Email       string    `json:"email"`
		PhoneNo     string    `json:"phone_no"`
		IsActive    bool      `json:"is_active"`
		CreatedBy   uuid.UUID `json:"created_by"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedBy   uuid.UUID `json:"updated_by"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
)

func (s StudentModel) Response() StudentResponse {
	return StudentResponse{
		ID:          s.ID,
		Name:        s.Name,
		Address:     s.Address,
		DateOfBirth: s.DateOfBirth,
		Gender:      s.Gender,
		Email:       s.Email,
		IsActive:    s.IsActive,
		PhoneNo:     s.PhoneNo,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedBy:   s.UpdatedBy.UUID,
		UpdatedAt:   s.UpdatedAt.Time,
	}
}

func GetAllStudent(ctx context.Context, db *sql.DB) ([]StudentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			date_of_birth,	
			gender,
			email,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM student`)

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []StudentModel
	for rows.Next() {
		var student StudentModel
		rows.Scan(
			&student.ID,
			&student.Name,
			&student.Address,
			&student.DateOfBirth,
			&student.Gender,
			&student.Email,
			&student.PhoneNo,
			&student.IsActive,
			&student.CreatedBy,
			&student.CreatedAt,
			&student.UpdatedBy,
			&student.UpdatedAt,
		)

		students = append(students, student)
	}

	return students, nil

}

func GetOneStudent(ctx context.Context, db *sql.DB, studentID uuid.UUID) (StudentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			address,
			date_of_birth,	
			gender,
			email,
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM student
		WHERE 
			id = $1
	`)

	var student StudentModel
	err := db.QueryRowContext(ctx, query, studentID).Scan(
		&student.ID,
		&student.Name,
		&student.Address,
		&student.DateOfBirth,
		&student.Gender,
		&student.Email,
		&student.PhoneNo,
		&student.IsActive,
		&student.CreatedBy,
		&student.CreatedAt,
		&student.UpdatedBy,
		&student.UpdatedAt,
	)

	if err != nil {
		return StudentModel{}, err
	}

	return student, nil

}
