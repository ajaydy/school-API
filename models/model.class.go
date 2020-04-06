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
	ClassModel struct {
		ID        uuid.UUID
		SessionID uuid.UUID
		Date      time.Time
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	ClassResponse struct {
		ID        uuid.UUID       `json:"id"`
		Session   SessionResponse `json:"session"`
		Date      time.Time       `json:"date"`
		IsDelete  bool            `json:"is_delete"`
		CreatedBy uuid.UUID       `json:"created_by"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedBy uuid.UUID       `json:"updated_by"`
		UpdatedAt time.Time       `json:"updated_at"`
	}
)

func (s ClassModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (ClassResponse, error) {

	session, err := GetOneSession(ctx, db, s.SessionID)
	if err != nil {
		logger.Err.Printf(`model.class.go/GetOneSession/%v`, err)
		return ClassResponse{}, nil
	}

	sessions, err := session.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.class.go/SessionResponse/%v`, err)
		return ClassResponse{}, nil
	}

	return ClassResponse{
		ID:        s.ID,
		Session:   sessions,
		Date:      s.Date,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil
}

func GetOneClass(ctx context.Context, db *sql.DB, classID uuid.UUID) (ClassModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			session_id,
			date,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM class
		WHERE 
			id= $1
	`)

	var class ClassModel
	err := db.QueryRowContext(ctx, query, classID).Scan(
		&class.ID,
		&class.SessionID,
		&class.Date,
		&class.IsDelete,
		&class.CreatedBy,
		&class.CreatedAt,
		&class.UpdatedBy,
		&class.UpdatedAt,
	)

	if err != nil {
		return ClassModel{}, err
	}

	return class, nil

}

func GetAllClassBySession(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ClassModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`AND LOWER(name) LIKE LOWER('%s')`, filter.Search)
	}

	query := fmt.Sprintf(`
			SELECT
			id,
			session_id,
			date,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
			FROM class
			WHERE 
			session_id = $1
			%s
			ORDER BY date %s
			LIMIT $2 OFFSET $3`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.SessionID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var classes []ClassModel
	for rows.Next() {
		var class ClassModel
		rows.Scan(
			&class.ID,
			&class.SessionID,
			&class.Date,
			&class.IsDelete,
			&class.CreatedBy,
			&class.CreatedAt,
			&class.UpdatedBy,
			&class.UpdatedAt,
		)

		classes = append(classes, class)
	}

	return classes, nil

}

func (s *ClassModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO class(
			session_id,
			date,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.SessionID, s.Date, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}
