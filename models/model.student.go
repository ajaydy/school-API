package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"school/helpers"
	"school/util"
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
		Gender      string          `json:"gender"`
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

	programResponse, err := program.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.student.go/programResponse/%v`, err)
		return StudentResponse{}, nil
	}

	gender, err := util.GetGender(s.Gender)
	if err != nil {
		logger.Err.Printf(`model.student.go/Gender/%v`, err)
		return StudentResponse{}, nil
	}

	return StudentResponse{
		ID:          s.ID,
		Program:     programResponse,
		Name:        s.Name,
		Address:     s.Address,
		DateOfBirth: s.DateOfBirth,
		Gender:      gender,
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
			student_code,
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
		&student.ProgramID,
		&student.Name,
		&student.Address,
		&student.DateOfBirth,
		&student.Gender,
		&student.Email,
		&student.PhoneNo,
		&student.StudentCode,
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

func GetAllStudent(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]StudentModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`AND LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

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
			student_code,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM student
		WHERE is_active = true
		%s
		ORDER BY name  %s
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []StudentModel
	for rows.Next() {
		var student StudentModel

		rows.Scan(
			&student.ID,
			&student.ProgramID,
			&student.Name,
			&student.Address,
			&student.DateOfBirth,
			&student.Gender,
			&student.Email,
			&student.PhoneNo,
			&student.StudentCode,
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

func GetOneStudentByCode(ctx context.Context, db *sql.DB, code string) (StudentModel, error) {

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
			student_code,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM student
		WHERE
			student_code = $1 AND is_active=true
	`)

	var student StudentModel
	err := db.QueryRowContext(ctx, query, code).Scan(
		&student.ID,
		&student.ProgramID,
		&student.Name,
		&student.Address,
		&student.DateOfBirth,
		&student.Gender,
		&student.Email,
		&student.PhoneNo,
		&student.StudentCode,
		&student.Password,
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

func (s *StudentModel) Insert(ctx context.Context, db *sql.DB) error {

	password, err := bcrypt.GenerateFromPassword([]byte(s.Password), 12)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		INSERT INTO student(
			name,
			program_id,
			address,
			date_of_birth,
			gender,
			email,
			student_code,
			password,
			phone_no,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,now())
		RETURNING id, created_at,is_active`)

	err = db.QueryRowContext(ctx, query,
		s.Name, s.ProgramID, s.Address, s.DateOfBirth, s.Gender, s.Email, s.StudentCode, password, s.PhoneNo, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *StudentModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE student
		SET
			name=$1,
			program_id=$2,
			address=$3,
			date_of_birth=$4,
			gender=$5,
			email=$6,
			phone_no=$7,
			updated_at=NOW(),
			updated_by=$8
		WHERE id=$9
		RETURNING id,created_at,updated_at,created_by,student_code,is_active`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.ProgramID, s.Address, s.DateOfBirth, s.Gender, s.Email, s.PhoneNo, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.StudentCode, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *StudentModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE student
		SET
			is_active=false,
			updated_by=$1,
			updated_at=NOW()
		WHERE id=$2`)

	_, err := db.ExecContext(ctx, query,
		s.UpdatedBy, s.ID)

	if err != nil {
		return err
	}

	return nil
}
