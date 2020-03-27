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
