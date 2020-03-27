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
	FacultyModel struct {
		ID           uuid.UUID
		Code         int
		Abbreviation string
		Name         string
		Description  string
		IsDelete     bool
		CreatedBy    uuid.UUID
		CreatedAt    time.Time
		UpdatedBy    uuid.NullUUID
		UpdatedAt    pq.NullTime
	}

	FacultyResponse struct {
		ID           uuid.UUID `json:"id"`
		Code         int       `json:"code"`
		Abbreviation string    `json:"abbreviation"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		IsDelete     bool      `json:"is_delete"`
		CreatedBy    uuid.UUID `json:"created_by"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedBy    uuid.UUID `json:"updated_by"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
)

func (s FacultyModel) Response() FacultyResponse {
	return FacultyResponse{
		ID:           s.ID,
		Code:         s.Code,
		Abbreviation: s.Abbreviation,
		Name:         s.Name,
		Description:  s.Description,
		IsDelete:     s.IsDelete,
		CreatedBy:    s.CreatedBy,
		CreatedAt:    s.CreatedAt,
		UpdatedBy:    s.UpdatedBy.UUID,
		UpdatedAt:    s.UpdatedAt.Time,
	}
}

func GetOneFaculty(ctx context.Context, db *sql.DB, facultyID uuid.UUID) (FacultyModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			code,
			abbreviation,
			name,
			description,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM faculty
		WHERE 
			id = $1
	`)

	var faculty FacultyModel
	err := db.QueryRowContext(ctx, query, facultyID).Scan(
		&faculty.ID,
		&faculty.Code,
		&faculty.Abbreviation,
		&faculty.Name,
		&faculty.Description,
		&faculty.IsDelete,
		&faculty.CreatedBy,
		&faculty.CreatedAt,
		&faculty.UpdatedBy,
		&faculty.UpdatedAt,
	)

	if err != nil {
		return FacultyModel{}, err
	}

	return faculty, nil

}
