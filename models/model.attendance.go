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
		ID        uuid.UUID
		StudentID uuid.UUID
		ClassID   uuid.UUID
		IsAttend  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	AttendanceResponse struct {
		ID        uuid.UUID       `json:"id"`
		Student   StudentResponse `json:"student"`
		Class     ClassResponse   `json:"class"`
		IsAttend  bool            `json:"is_attend"`
		CreatedBy uuid.UUID       `json:"created_by"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedBy uuid.UUID       `json:"updated_by"`
		UpdatedAt time.Time       `json:"updated_at"`
	}
)

func (s AttendanceModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (AttendanceResponse, error) {

	student, err := GetOneStudent(ctx, db, s.StudentID)
	if err != nil {
		logger.Err.Printf(`model.attendance.go/GetOneStudent/%v`, err)
		return AttendanceResponse{}, nil
	}

	students, err := student.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.attendance.go/studentResponse/%v`, err)
		return AttendanceResponse{}, nil
	}

	class, err := GetOneClass(ctx, db, s.ClassID)
	if err != nil {
		logger.Err.Printf(`model.attendance.go/GetOneClass/%v`, err)
		return AttendanceResponse{}, nil
	}

	classes, err := class.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.attendance.go/classResponse/%v`, err)
		return AttendanceResponse{}, nil
	}

	return AttendanceResponse{
		ID:        s.ID,
		Student:   students,
		Class:     classes,
		IsAttend:  s.IsAttend,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil
}

func GetOneAttendance(ctx context.Context, db *sql.DB, attendanceID uuid.UUID) (AttendanceModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			student_id,
			class_id,
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
		&attendance.ClassID,
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

func GetAllAttendanceByClass(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]AttendanceModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
			SELECT
			a.id,
			student_id,
			class_id,
			a.is_attend,
			a.created_by,
			a.created_at,
			a.updated_by,
			a.updated_at
			FROM attendance	a
			INNER JOIN student s ON a.student_id = s.id
			WHERE 
			class_id = $1
			
			%s
			ORDER BY s.name %s
			LIMIT $2 OFFSET $3`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.ClassID, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var attendances []AttendanceModel
	for rows.Next() {
		var attendance AttendanceModel
		rows.Scan(
			&attendance.ID,
			&attendance.StudentID,
			&attendance.ClassID,
			&attendance.IsAttend,
			&attendance.CreatedBy,
			&attendance.CreatedAt,
			&attendance.UpdatedBy,
			&attendance.UpdatedAt,
		)

		attendances = append(attendances, attendance)
	}

	return attendances, nil

}

func (s *AttendanceModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO attendance(
			student_id,
			class_id,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,now())
		RETURNING id, created_at,is_attend`)

	err := db.QueryRowContext(ctx, query,
		s.StudentID, s.ClassID, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsAttend,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *AttendanceModel) UpdateIsAttend(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE attendance
		SET
			is_attend=$1,
			updated_at=NOW(),
			updated_by=$2
		WHERE id=$3
		RETURNING id,student_id,class_id,created_at,updated_at,created_by`)

	err := db.QueryRowContext(ctx, query,
		s.IsAttend, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.StudentID, &s.ClassID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}
