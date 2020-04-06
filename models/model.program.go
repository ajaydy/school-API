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
	ProgramModel struct {
		ID          uuid.UUID
		FacultyID   uuid.UUID
		Name        string
		Code        int
		Description string
		IsDelete    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	ProgramResponse struct {
		ID          uuid.UUID       `json:"id"`
		Faculty     FacultyResponse `json:"faculty"`
		Name        string          `json:"name"`
		Code        int             `json:"code"`
		Description string          `json:"description""`
		IsDelete    bool            `json:"is_delete"`
		CreatedBy   uuid.UUID       `json:"created_by"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedBy   uuid.UUID       `json:"updated_by"`
		UpdatedAt   time.Time       `json:"updated_at"`
	}
)

func (s ProgramModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (ProgramResponse, error) {

	faculty, err := GetOneFaculty(ctx, db, s.FacultyID)
	if err != nil {
		logger.Err.Printf(`model.program.go/GetOneFaculty/%v`, err)
		return ProgramResponse{}, nil
	}

	return ProgramResponse{
		ID:          s.ID,
		Faculty:     faculty.Response(),
		Name:        s.Name,
		Code:        s.Code,
		Description: s.Description,
		IsDelete:    s.IsDelete,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedBy:   s.UpdatedBy.UUID,
		UpdatedAt:   s.UpdatedAt.Time,
	}, nil

}

func GetOneProgram(ctx context.Context, db *sql.DB, programID uuid.UUID) (ProgramModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			faculty_id,
			name,
			code,
			description,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM program
		WHERE 
			id = $1
	`)

	var program ProgramModel
	err := db.QueryRowContext(ctx, query, programID).Scan(
		&program.ID,
		&program.FacultyID,
		&program.Name,
		&program.Code,
		&program.Description,
		&program.IsDelete,
		&program.CreatedBy,
		&program.CreatedAt,
		&program.UpdatedBy,
		&program.UpdatedAt,
	)

	if err != nil {
		return ProgramModel{}, err
	}

	return program, nil

}

func GetAllProgram(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ProgramModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`AND LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			faculty_id,
			name,
			code,
			description,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM program
		WHERE is_delete = false 
		%s
		ORDER BY name %s 
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var programs []ProgramModel
	for rows.Next() {
		var program ProgramModel

		rows.Scan(
			&program.ID,
			&program.FacultyID,
			&program.Name,
			&program.Code,
			&program.Description,
			&program.IsDelete,
			&program.CreatedBy,
			&program.CreatedAt,
			&program.UpdatedBy,
			&program.UpdatedAt,
		)

		programs = append(programs, program)
	}

	return programs, nil

}

func (s *ProgramModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO program(
			faculty_id,
			name,
			code,
			description,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,$5,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.FacultyID, s.Name, s.Code, s.Description, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ProgramModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE program
		SET
			faculty_id=$1,
			name=$2,
			code=$3,
			description=$4,
			updated_at=NOW(),
			updated_by=$5
		WHERE id=$6
		RETURNING id,created_at,updated_at,created_by,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.FacultyID, s.Name, s.Code, s.Description, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ProgramModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE program
		SET
			is_delete=true,
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
