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
	AttendanceModel struct {
		ID              uuid.UUID
		StudentID       uuid.UUID
		StudentEnrollID uuid.UUID
		IsAttend        bool
		CreatedBy       uuid.UUID
		CreatedAt       time.Time
		UpdatedBy       uuid.NullUUID
		UpdatedAt       pq.NullTime
	}

	AttendanceResponse struct {
		ID            uuid.UUID             `json:"id"`
		Student       StudentResponse       `json:"student"`
		StudentEnroll StudentEnrollResponse `json:"student_enroll"`
		IsAttend      bool                  `json:"is_attend"`
		CreatedBy     uuid.UUID             `json:"created_by"`
		CreatedAt     time.Time             `json:"created_at"`
		UpdatedBy     uuid.UUID             `json:"updated_by"`
		UpdatedAt     time.Time             `json:"updated_at"`
	}
)

func (s AttendanceModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (AttendanceResponse, error) {

	student, err := GetOneStudent(ctx, db, s.StudentID)
	if err != nil {
		logger.Err.Printf(`model.attendance.go/GetOneStudent/%v`, err)
		return AttendanceResponse{}, nil
	}

	studentEnroll, err := GetOneStudentEnroll(ctx, db, s.StudentEnrollID)
	if err != nil {
		logger.Err.Printf(`model.attendance.go/GetOneStudentEnroll/%v`, err)
		return AttendanceResponse{}, nil
	}

	students, err := student.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.attendance.go/studentResponse/%v`, err)
		return AttendanceResponse{}, nil
	}

	studentEnrolls, err := studentEnroll.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.attendance.go/studentEnrollResponse/%v`, err)
		return AttendanceResponse{}, nil
	}

	return AttendanceResponse{
		ID:            s.ID,
		Student:       students,
		StudentEnroll: studentEnrolls,
		IsAttend:      s.IsAttend,
		CreatedBy:     s.CreatedBy,
		CreatedAt:     s.CreatedAt,
		UpdatedBy:     s.UpdatedBy.UUID,
		UpdatedAt:     s.UpdatedAt.Time,
	}, nil
}

func GetOneAttendance(ctx context.Context, db *sql.DB, attendanceID uuid.UUID) (AttendanceModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			student_id,
			student_enroll_id,
			is_attend,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM attendance
		WHERE 
			id = $1
	`)

	var attendance AttendanceModel
	err := db.QueryRowContext(ctx, query, attendanceID).Scan(
		&attendance.ID,
		&attendance.StudentID,
		&attendance.StudentEnrollID,
		&attendance.IsAttend,
		&attendance.CreatedBy,
		&attendance.CreatedAt,
		&attendance.UpdatedBy,
		&attendance.UpdatedAt,
	)

	if err != nil {
		return AttendanceModel{}, err
	}

	return attendance, nil

}
