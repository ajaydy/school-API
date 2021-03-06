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
	LecturerModel struct {
		ID        uuid.UUID
		ProgramID uuid.UUID
		Name      string
		PhoneNo   string
		Address   string
		Email     string
		Gender    int
		Password  string
		IsActive  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	LecturerResponse struct {
		ID        uuid.UUID       `json:"id"`
		Program   ProgramResponse `json:"program"`
		Name      string          `json:"name"`
		PhoneNo   string          `json:"phone_no"`
		Address   string          `json:"address"`
		Email     string          `json:"email"`
		Gender    string          `json:"gender"`
		IsActive  bool            `json:"is_active"`
		CreatedBy uuid.UUID       `json:"created_by"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedBy uuid.UUID       `json:"updated_by"`
		UpdatedAt time.Time       `json:"updated_at"`
	}

	LecturerUpdatePasswordResponse struct {
		Message string `json:"message"`
	}
)

func (s LecturerModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (LecturerResponse, error) {

	program, err := GetOneProgram(ctx, db, s.ProgramID)
	if err != nil {
		logger.Err.Printf(`model.lecturer.go/GetOneProgram/%v`, err)
		return LecturerResponse{}, nil
	}

	programResponse, err := program.Response(ctx, db, logger)

	if err != nil {
		logger.Err.Printf(`model.lecturer.go/ProgramResponse/%v`, err)
		return LecturerResponse{}, nil
	}

	gender, err := util.GetGender(s.Gender)

	if err != nil {
		logger.Err.Printf(`model.lecturer.go/GetGender/%v`, err)
		return LecturerResponse{}, nil
	}

	return LecturerResponse{
		ID:        s.ID,
		Program:   programResponse,
		Name:      s.Name,
		PhoneNo:   s.PhoneNo,
		Address:   s.Address,
		Email:     s.Email,
		Gender:    gender,
		IsActive:  s.IsActive,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil
}

func GetOneLecturer(ctx context.Context, db *sql.DB, lecturerID uuid.UUID) (LecturerModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			program_id,
			name,
			address,
			phone_no,
			email,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM lecturer
		WHERE 
			id = $1
	`)

	var lecturer LecturerModel
	err := db.QueryRowContext(ctx, query, lecturerID).Scan(
		&lecturer.ID,
		&lecturer.ProgramID,
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

func GetOneLecturerByEmail(ctx context.Context, db *sql.DB, email string) (LecturerModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			program_id,
			name,
			address,
			phone_no,
			email,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM lecturer
		WHERE 
			email = $1 
	`)

	var lecturer LecturerModel
	err := db.QueryRowContext(ctx, query, email).Scan(
		&lecturer.ID,
		&lecturer.ProgramID,
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

func GetAllLecturer(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]LecturerModel, error) {

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
			phone_no,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at,
			email
		FROM lecturer
		WHERE is_active=true
		%s
		ORDER BY name  %s
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var lecturers []LecturerModel
	for rows.Next() {
		var lecturer LecturerModel
		rows.Scan(
			&lecturer.ID,
			&lecturer.ProgramID,
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

func (s *LecturerModel) Insert(ctx context.Context, db *sql.DB) error {

	password, err := bcrypt.GenerateFromPassword([]byte(s.Password), 12)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		INSERT INTO lecturer(
			name,
			program_id,
			address,
			email,
			phone_no,
			gender,
			password,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,$5,$6,$7,$8,now())
		RETURNING id, created_at,is_active`)

	err = db.QueryRowContext(ctx, query,
		s.Name, s.ProgramID, s.Address, s.Email, s.PhoneNo, s.Gender, password, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *LecturerModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE lecturer
		SET
			name=$1,
			program_id=$2,
			address=$3,
			email=$4,
			phone_no=$5,
			updated_at=NOW(),
			updated_by=$6
		WHERE id=$7
		RETURNING id,created_at,updated_at,created_by,is_active,gender`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.ProgramID, s.Address, s.Email, s.PhoneNo, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsActive, &s.Gender,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *LecturerModel) PasswordUpdate(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE lecturer
		SET
			password = $1,
			updated_at=NOW(),
			updated_by=$2
		WHERE id=$3
		RETURNING id,created_at,updated_at,created_by,is_active`)

	err := db.QueryRowContext(ctx, query,
		s.Password, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *LecturerModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE lecturer
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
