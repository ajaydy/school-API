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
	StudentEnrollModel struct {
		ID        uuid.UUID
		SessionID uuid.UUID
		StudentID uuid.UUID
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	StudentEnrollResponse struct {
		ID        uuid.UUID       `json:"id"`
		Session   SessionResponse `json:"session"`
		Student   StudentResponse `json:"student"`
		IsDelete  bool            `json:"is_delete"`
		CreatedBy uuid.UUID       `json:"created_by"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedBy uuid.UUID       `json:"updated_by"`
		UpdatedAt time.Time       `json:"updated_at"`
	}
)

func (s StudentEnrollModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (StudentEnrollResponse, error) {

	session, err := GetOneSession(ctx, db, s.SessionID)
	if err != nil {
		logger.Err.Printf(`model.student.enroll.go/GetOneSession/%v`, err)
		return StudentEnrollResponse{}, nil
	}

	student, err := GetOneStudent(ctx, db, s.StudentID)
	if err != nil {
		logger.Err.Printf(`model.student.enroll.go/GetOneStudent/%v`, err)
		return StudentEnrollResponse{}, err
	}

	sessions, err := session.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.student.enroll.go/sessionResponse/%v`, err)
		return StudentEnrollResponse{}, err
	}

	students, err := student.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.student.enroll.go/studentResponse/%v`, err)
		return StudentEnrollResponse{}, err
	}

	return StudentEnrollResponse{
		ID:        s.ID,
		Session:   sessions,
		Student:   students,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil
}

func GetOneStudentEnroll(ctx context.Context, db *sql.DB, studentEnrollID uuid.UUID) (StudentEnrollModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			session_id,
			student_id,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM student_enroll
		WHERE 
			id = $1
	`)

	var student StudentEnrollModel
	err := db.QueryRowContext(ctx, query, studentEnrollID).Scan(
		&student.ID,
		&student.SessionID,
		&student.StudentID,
		&student.IsDelete,
		&student.CreatedBy,
		&student.CreatedAt,
		&student.UpdatedBy,
		&student.UpdatedAt,
	)

	if err != nil {
		return StudentEnrollModel{}, err
	}

	return student, nil

}

func GetOneStudentEnrollBySession(ctx context.Context, db *sql.DB, studentEnrollID uuid.UUID) (StudentEnrollModel, error) {

	query := fmt.Sprintf(`
	SELECT
			se.id,
			session_id,
			student_id,
			se.is_delete,
			se.created_by,
			se.created_at,
			se.updated_by,
			se.updated_at
		FROM student_enroll se
		INNER JOIN student s ON se.student_id = s.id
		WHERE session_id = $1
	`)

	var student StudentEnrollModel
	err := db.QueryRowContext(ctx, query, studentEnrollID).Scan(
		&student.ID,
		&student.SessionID,
		&student.StudentID,
		&student.IsDelete,
		&student.CreatedBy,
		&student.CreatedAt,
		&student.UpdatedBy,
		&student.UpdatedAt,
	)

	if err != nil {
		return StudentEnrollModel{}, err
	}

	return student, nil

}

//select  student_id from student_enroll inner join session where session_id = param.id
func (s *StudentEnrollModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO student_enroll(
			session_id,
			student_id,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.SessionID, s.StudentID, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func GetTimetableForStudent(ctx context.Context, db *sql.DB, filter helpers.Filter, studentID uuid.UUID) ([]StudentEnrollModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			se.id,
			session_id,
			student_id,
			se.is_delete,
			se.created_by,
			se.created_at,
			se.updated_by,
			se.updated_at
		FROM student_enroll se
		INNER JOIN session s ON se.session_id=s.id
		WHERE student_id = $1
		%s
		ORDER BY s.day  %s, s.start_time %s
		LIMIT $2 OFFSET $3`, searchQuery, filter.Dir, filter.Dir)

	rows, err := db.QueryContext(ctx, query, studentID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []StudentEnrollModel
	for rows.Next() {
		var student StudentEnrollModel

		rows.Scan(
			&student.ID,
			&student.SessionID,
			&student.StudentID,
			&student.IsDelete,
			&student.CreatedBy,
			&student.CreatedAt,
			&student.UpdatedBy,
			&student.UpdatedAt,
		)

		students = append(students, student)
	}

	return students, nil
}

func GetAllStudentEnrollBySession(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]StudentEnrollModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			se.id,
			session_id,
			student_id,
			se.is_delete,
			se.created_by,
			se.created_at,
			se.updated_by,
			se.updated_at
		FROM student_enroll se
		INNER JOIN student s ON se.student_id = s.id
		WHERE session_id = $1
		%s
		ORDER BY  s.name %s
		LIMIT $2 OFFSET $3`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.SessionID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []StudentEnrollModel
	for rows.Next() {
		var student StudentEnrollModel

		rows.Scan(
			&student.ID,
			&student.SessionID,
			&student.StudentID,
			&student.IsDelete,
			&student.CreatedBy,
			&student.CreatedAt,
			&student.UpdatedBy,
			&student.UpdatedAt,
		)

		students = append(students, student)
	}

	return students, nil

}

func (s *StudentEnrollModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE student_enroll
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
