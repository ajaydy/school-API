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

func GetAllFaculty(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]FacultyModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`AND LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

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
		WHERE is_delete = false
		%s
		ORDER BY name  %s
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var faculties []FacultyModel
	for rows.Next() {
		var faculty FacultyModel

		rows.Scan(
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

		faculties = append(faculties, faculty)
	}

	return faculties, nil

}

func (s *FacultyModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO faculty(
			code,
			abbreviation,
			name,
			description,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,$5,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.Code, s.Abbreviation, s.Name, s.Description, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *FacultyModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE faculty
		SET
			code=$1,
			abbreviation=$2,
			name=$3,
			description=$4,
			updated_at=NOW(),
			updated_by=$5
		WHERE id=$6
		RETURNING id,created_at,updated_at,created_by,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.Code, s.Abbreviation, s.Name, s.Description, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *FacultyModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE faculty
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
