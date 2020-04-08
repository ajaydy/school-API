package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"school/helpers"
	"strings"
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

func GetAllResult(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ResultModel, error) {

	var searchQuery string

	var studentIDQuery string

	var subjectIDQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`AND LOWER(st.name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	if filter.StudentID != uuid.Nil {
		studentIDQuery = fmt.Sprintf(`AND st.id = '%s'`, filter.StudentID)
	}

	if filter.SubjectID != uuid.Nil {
		subjectIDQuery = fmt.Sprintf(`AND s.subject_id = '%s'`, filter.SubjectID)
	}

	query := fmt.Sprintf(`
		SELECT
			r.id,
			r.student_enroll_id,
			r.grade,
			r.marks,
			r.is_delete,
			r.created_by,
			r.created_at,
			r.updated_by,
			r.updated_at
		FROM result r
		INNER JOIN student_enroll se ON r.student_enroll_id = se.id
		INNER JOIN student st ON se.student_id = st.id
		INNER JOIN session s ON se.session_id = s.id
		WHERE r.is_delete = false
		%s %s %s
		ORDER BY  session_id %s ,st.name  %s
		LIMIT $1 OFFSET $2`, searchQuery, studentIDQuery, subjectIDQuery, filter.Dir, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)
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

func GetAllResultByStudentEnroll(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ResultModel, error) {

	var filters []string

	if filter.StudentEnrollID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			student_enroll_id = '%s'`,
			filter.StudentEnrollID))
	}
	filterJoin := strings.Join(filters, " AND ")
	if filterJoin != "" {
		filterJoin = fmt.Sprintf("WHERE %s", filterJoin)
	}

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
		FROM result r
		%s
		LIMIT $1 OFFSET $2`, filterJoin)
	fmt.Println(query)
	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

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

func GetAllResultForOneStudent(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ResultModel, error) {

	var filters []string

	if filter.Search != "" {
		filters = append(filters, fmt.Sprintf(`
		LOWER(su.name) LIKE LOWER('%%%s%%')`,
			filter.Search))
	}

	if filter.StudentID != uuid.Nil {
		filters = append(filters, fmt.Sprintf(`
			se.student_id = '%s'`,
			filter.StudentID))
	}
	filterJoin := strings.Join(filters, " AND ")
	if filterJoin != "" {
		filterJoin = fmt.Sprintf("WHERE %s", filterJoin)
	}
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
		%s
		ORDER BY  su.name %s
		LIMIT $1 OFFSET $2`, filterJoin, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)
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
		RETURNING id,updated_at,created_at,created_by,student_enroll_id`)

	err := db.QueryRowContext(ctx, query,
		s.Grade, s.Marks, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.UpdatedAt, &s.CreatedAt, &s.CreatedBy, &s.StudentEnrollID,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ResultModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE result
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
