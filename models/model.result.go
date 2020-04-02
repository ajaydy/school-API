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
	ResultModel struct {
		ID              uuid.UUID
		StudentEnrollID uuid.UUID
		Grade           string
		Marks           int
		IsDelete        bool
		CreatedBy       uuid.UUID
		CreatedAt       time.Time
		UpdatedBy       uuid.NullUUID
		UpdatedAt       pq.NullTime
	}

	ResultResponse struct {
		ID            uuid.UUID             `json:"id"`
		StudentEnroll StudentEnrollResponse `json:"student_enroll"`
		Grade         string                `json:"grade"`
		Marks         int                   `json:"marks"`
		IsDelete      bool                  `json:"is_delete"`
		CreatedBy     uuid.UUID             `json:"created_by"`
		CreatedAt     time.Time             `json:"created_at"`
		UpdatedBy     uuid.UUID             `json:"updated_by"`
		UpdatedAt     time.Time             `json:"updated_at"`
	}
)

func (s ResultModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (ResultResponse, error) {

	studentEnroll, err := GetOneStudentEnroll(ctx, db, s.StudentEnrollID)
	if err != nil {
		logger.Err.Printf(`model.result.go/GetOneStudentEnroll/%v`, err)
		return ResultResponse{}, nil
	}

	studentEnrolls, err := studentEnroll.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.result.go/studentEnrollResponse/%v`, err)
		return ResultResponse{}, nil
	}

	return ResultResponse{
		ID:            s.ID,
		StudentEnroll: studentEnrolls,
		Grade:         s.Grade,
		Marks:         s.Marks,
		IsDelete:      s.IsDelete,
		CreatedBy:     s.CreatedBy,
		CreatedAt:     s.CreatedAt,
		UpdatedBy:     s.UpdatedBy.UUID,
		UpdatedAt:     s.UpdatedAt.Time,
	}, nil
}

func GetOneResult(ctx context.Context, db *sql.DB, resultID uuid.UUID) (ResultModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			student_enroll_id,
			grade,
			marks,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM result
		WHERE 
			id = $1
	`)

	var result ResultModel
	err := db.QueryRowContext(ctx, query, resultID).Scan(
		&result.ID,
		&result.StudentEnrollID,
		&result.Grade,
		&result.Marks,
		&result.IsDelete,
		&result.CreatedBy,
		&result.CreatedAt,
		&result.UpdatedBy,
		&result.UpdatedAt,
	)

	if err != nil {
		return ResultModel{}, err
	}

	return result, nil

}

func GetAllResultForOneStudent(ctx context.Context, db *sql.DB, filter helpers.Filter, studentID uuid.UUID) ([]ResultModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	fmt.Println(studentID)

	query := fmt.Sprintf(`
		SELECT
			r.id,
			student_enroll_id,
			grade,
			marks,
			r.is_delete,
			r.created_by,
			r.created_at,
			r.updated_by,
			r.updated_at
		FROM result r
		INNER JOIN student_enroll se ON r.student_enroll_id = se.id
		INNER JOIN session s ON se.session_id = s.id
		INNER JOIN subject su ON s.subject_id = su.id
		WHERE 
		se.student_id = $1
		%s
		ORDER BY  su.name %s
		LIMIT $2 OFFSET $3`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, studentID, filter.Limit, filter.Offset)
	fmt.Println(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []ResultModel
	for rows.Next() {
		var result ResultModel
		rows.Scan(
			&result.ID,
			&result.StudentEnrollID,
			&result.Grade,
			&result.Marks,
			&result.IsDelete,
			&result.CreatedBy,
			&result.CreatedAt,
			&result.UpdatedBy,
			&result.UpdatedAt,
		)

		results = append(results, result)
	}

	return results, nil

}

func (s *ResultModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO result(
			student_enroll_id,
			grade,
			marks,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.StudentEnrollID, s.Grade, s.Marks, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ResultModel) UpdateByStudentEnroll(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE result
		SET
			grade=$1,
			marks=$2,
			updated_at=NOW(),
			updated_by=$3
		WHERE student_enroll_id = $4
		RETURNING id,updated_at,created_at,created_by`)

	err := db.QueryRowContext(ctx, query,
		s.Grade, s.Marks, s.UpdatedBy, s.StudentEnrollID).Scan(
		&s.ID, &s.UpdatedAt, &s.CreatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ResultModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE result
		SET
			grade=$1,
			marks=$2,
			updated_at=NOW(),
			updated_by=$3
		WHERE id=$4
		RETURNING id,updated_at,created_at,created_by`)

	err := db.QueryRowContext(ctx, query,
		s.Grade, s.Marks, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.UpdatedAt, &s.CreatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}
