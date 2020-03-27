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
		Faculty     FacultyResponse `json:"faculty_id"`
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
