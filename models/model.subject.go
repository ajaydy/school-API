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
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Duration    int       `json:"duration"`
		IsDelete    bool      `json:"is_delete"`
		CreatedBy   uuid.UUID `json:"created_by"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedBy   uuid.UUID `json:"updated_by"`
		UpdatedAt   time.Time `json:"updated_at"`
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

func GetAllSubject(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]SubjectModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

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
		%s
		ORDER BY name  %s
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subjects []SubjectModel
	for rows.Next() {
		var subject SubjectModel
		rows.Scan(
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

		subjects = append(subjects, subject)
	}

	return subjects, nil

}

func (s *SubjectModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO subject(
			name,
			description,
			duration,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.Description, s.Duration, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

//
//func (s *SubjectModel) Update(ctx context.Context, db *sql.DB) error {
//
//	query := fmt.Sprintf(`
//		UPDATE subject
//		SET
//			"name"=$1,
//			"description"=$2,
//			duration=$3,
//			updated_at=NOW(),
//			updated_by=$4
//		WHERE id=$5
//		RETURNING id,created_at`)
//
//	err := db.QueryRowContext(ctx, query,
//		s.Name, s.Description, s.Duration, s.UpdatedBy, s.ID).Scan(
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
