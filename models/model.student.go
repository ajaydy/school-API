package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"school/helpers"
	"time"
)

type (
	StudentModel struct {
		ID          uuid.UUID
		ProgramID   uuid.UUID
		Name        string
		Address     string
		DateOfBirth time.Time
		Gender      int
		Email       string
		PhoneNo     string
		StudentCode string
		Password    string
		IsActive    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	StudentResponse struct {
		ID          uuid.UUID       `json:"id"`
		Program     ProgramResponse `json:"program"`
		Name        string          `json:"name"`
		Address     string          `json:"address"`
		DateOfBirth time.Time       `json:"date_of_birth"`
		Gender      int             `json:"gender"`
		Email       string          `json:"email"`
		PhoneNo     string          `json:"phone_no"`
		StudentCode string          `json:"student_code"`
		IsActive    bool            `json:"is_active"`
		CreatedBy   uuid.UUID       `json:"created_by"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedBy   uuid.UUID       `json:"updated_by"`
		UpdatedAt   time.Time       `json:"updated_at"`
	}
)

func (s StudentModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (StudentResponse, error) {

	program, err := GetOneProgram(ctx, db, s.ProgramID)
	if err != nil {
		logger.Err.Printf(`model.student.go/GetOneProgram/%v`, err)
		return StudentResponse{}, nil
	}

	programs, err := program.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.student.go/programResponse/%v`, err)
		return StudentResponse{}, nil
	}

	return StudentResponse{
		ID:          s.ID,
		Program:     programs,
		Name:        s.Name,
		Address:     s.Address,
		DateOfBirth: s.DateOfBirth,
		Gender:      s.Gender,
		Email:       s.Email,
		PhoneNo:     s.PhoneNo,
		StudentCode: s.StudentCode,
		IsActive:    s.IsActive,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedBy:   s.UpdatedBy.UUID,
		UpdatedAt:   s.UpdatedAt.Time,
	}, nil

}

func GetOneStudent(ctx context.Context, db *sql.DB, studentID uuid.UUID) (StudentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			program_id,
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

//func GetAllStudent(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]StudentModel, error) {
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
//			date_of_birth,
//			gender,
//			email,
//			phone_no,
//			is_active,
//			created_by,
//			created_at,
//			updated_by,
//			updated_at
//		FROM student
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
//	var students []StudentModel
//	for rows.Next() {
//		var student StudentModel
//
//		rows.Scan(
//			&student.ID,
//			&student.Name,
//			&student.Address,
//			&student.DateOfBirth,
//			&student.Gender,
//			&student.Email,
//			&student.PhoneNo,
//			&student.IsActive,
//			&student.CreatedBy,
//			&student.CreatedAt,
//			&student.UpdatedBy,
//			&student.UpdatedAt,
//		)
//
//		students = append(students, student)
//	}
//
//	return students, nil
//
//}

//func GetOneStudentByEmail(ctx context.Context, db *sql.DB, email string) (StudentModel, error) {
//
//	query := fmt.Sprintf(`
//		SELECT
//			id,
//			name,
//			address,
//			date_of_birth,
//			gender,
//			email,
//			phone_no,
//			password,
//			is_active,
//			created_by,
//			created_at,
//			updated_by,
//			updated_at
//		FROM student
//		WHERE
//			email = $1
//	`)
//
//	var student StudentModel
//	err := db.QueryRowContext(ctx, query, email).Scan(
//		&student.ID,
//		&student.Name,
//		&student.Address,
//		&student.DateOfBirth,
//		&student.Gender,
//		&student.Email,
//		&student.PhoneNo,
//		&student.Password,
//		&student.IsActive,
//		&student.CreatedBy,
//		&student.CreatedAt,
//		&student.UpdatedBy,
//		&student.UpdatedAt,
//	)
//
//	if err != nil {
//		return StudentModel{}, err
//	}
//
//	return student, nil
//
//}
//
//func (s *StudentModel) Insert(ctx context.Context, db *sql.DB) error {
//
//	query := fmt.Sprintf(`
//		INSERT INTO student(
//			name,
//			address,
//			date_of_birth,
//			gender,
//			email,
//			phone_no,
//			password,
//			created_by,
//			created_at)
//		VALUES(
//		$1,$2,$3,$4,$5,$6,$7,$8,now())
//		RETURNING id, created_at`)
//
//	err := db.QueryRowContext(ctx, query,
//		s.Name, s.Address, s.DateOfBirth, s.Gender, s.Email, s.PhoneNo, s.Password, s.CreatedBy).Scan(
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
//func (s *StudentModel) Update(ctx context.Context, db *sql.DB) error {
//
//	query := fmt.Sprintf(`
//		UPDATE student
//		SET
//			"name"=$1,
//			address=$2,
//			date_of_birth=$3,
//			gender=$4,
//			email=$5,
//			phone_no=$6,
//			updated_at=NOW(),
//			updated_by=$7
//		WHERE id=$8
//		RETURNING id,created_at`)
//
//	err := db.QueryRowContext(ctx, query,
//		s.Name, s.Address, s.DateOfBirth, s.Gender, s.Email, s.PhoneNo, s.UpdatedBy, s.ID).Scan(
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
