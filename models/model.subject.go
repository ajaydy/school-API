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
	SubjectModel struct {
		ID          uuid.UUID
		Name        string
		Description string
		Duration    int
		IsDelete    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}
	SubjectResponse struct {
		ID          uuid.UUID
		Name        string
		Description string
		Duration    int
		IsDelete    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.UUID
		UpdatedAt   time.Time
	}
)

func (s SubjectModel) Response() SubjectResponse {
	return SubjectResponse{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Duration:    s.Duration,
		IsDelete:    s.IsDelete,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt,
		UpdatedBy:   s.UpdatedBy.UUID,
		UpdatedAt:   s.UpdatedAt.Time,
	}
}

func GetOneSubject(ctx context.Context, db *sql.DB, subjectID uuid.UUID) (SubjectModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			description,
			duration,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM subject
		WHERE 
			id = $1
	`)

	var subject SubjectModel
	err := db.QueryRowContext(ctx, query, subjectID).Scan(
		&subject.ID,
		&subject.Name,
		&subject.Description,
		&subject.Duration,
		&subject.IsDelete,
		&subject.CreatedBy,
		&subject.CreatedAt,
		&subject.UpdatedBy,
		&subject.UpdatedAt,
	)

	if err != nil {
		return SubjectModel{}, err
	}

	return subject, nil

}
