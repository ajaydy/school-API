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
	SemesterModel struct {
		ID        uuid.UUID
		Year      string
		Month     int
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	SemesterResponse struct {
		ID        uuid.UUID
		Year      string
		Month     int
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.UUID
		UpdatedAt time.Time
	}
)

func (s SemesterModel) Response() SemesterResponse {
	return SemesterResponse{
		ID:        s.ID,
		Year:      s.Year,
		Month:     s.Month,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}
}

func GetOneSemester(ctx context.Context, db *sql.DB, semesterID uuid.UUID) (SemesterModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			year,
			month,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM semester
		WHERE 
			id = $1
	`)

	var semester SemesterModel
	err := db.QueryRowContext(ctx, query, semesterID).Scan(
		&semester.ID,
		&semester.Year,
		&semester.Month,
		&semester.IsDelete,
		&semester.CreatedBy,
		&semester.CreatedAt,
		&semester.UpdatedBy,
		&semester.UpdatedAt,
	)

	if err != nil {
		return SemesterModel{}, err
	}

	return semester, nil

}
