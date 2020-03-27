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
	SessionModel struct {
		ID          uuid.UUID
		SubjectID   uuid.UUID
		LecturerID  uuid.UUID
		IntakeID    uuid.UUID
		ClassroomID uuid.UUID
		ProgramID   uuid.UUID
		Day         int
		StartTime   time.Time
		EndTime     time.Time
		IsDelete    bool
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	SessionResponse struct {
		ID        uuid.UUID         `json:"id"`
		Subject   SubjectResponse   `json:"subject"`
		Lecturer  LecturerResponse  `json:"lecturer"`
		Intake    IntakeResponse    `json:"intake"`
		Classroom ClassRoomResponse `json:"classroom"`
		Program   ProgramResponse   `json:"program"`
		Day       int               `json:"day"`
		StartTime time.Time         `json:"start_time"`
		EndTime   time.Time         `json:"end_time"`
		IsDelete  bool              `json:"is_delete"`
		CreatedBy uuid.UUID         `json:"created_by"`
		CreatedAt time.Time         `json:"created_at"`
		UpdatedBy uuid.UUID         `json:"updated_by"`
		UpdatedAt time.Time         `json:"updated_at"`
	}
)

func (s SessionModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (SessionResponse, error) {

	lecturer, err := GetOneLecturer(ctx, db, s.LecturerID)
	if err != nil {
		logger.Err.Printf(`model.session.go/GetOneLecturer/%v`, err)
		return SessionResponse{}, nil
	}

	intake, err := GetOneIntake(ctx, db, s.IntakeID)
	if err != nil {
		logger.Err.Printf(`model.session.go/GetOneIntake/%v`, err)
		return SessionResponse{}, nil
	}

	subject, err := GetOneSubject(ctx, db, s.IntakeID)
	if err != nil {
		logger.Err.Printf(`model.session.go/GetOneSubject/%v`, err)
		return SessionResponse{}, nil
	}

	classroom, err := GetOneClassroom(ctx, db, s.IntakeID)
	if err != nil {
		logger.Err.Printf(`model.session.go/GetOneClassroom/%v`, err)
		return SessionResponse{}, nil
	}

	classrooms, err := classroom.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.session.go/classroomResponse/%v`, err)
		return SessionResponse{}, nil
	}

	return SessionResponse{
		ID:        s.ID,
		Subject:   subject.Response(),
		Lecturer:  lecturer.Response(),
		Intake:    intake.Response(),
		Classroom: classrooms,

		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil
}

func GetOneSession(ctx context.Context, db *sql.DB, sessionID uuid.UUID) (SessionModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			subject_id,
			lecturer_id,
			intake_id,
			classroom_id,
			program_id,
			day,
			start_time,
			end_time,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM session
		WHERE 
			id = $1
	`)

	var session SessionModel
	err := db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID,
		&session.SubjectID,
		&session.LecturerID,
		&session.IntakeID,
		&session.ClassroomID,
		&session.ProgramID,
		&session.Day,
		&session.StartTime,
		&session.EndTime,
		&session.IsDelete,
		&session.CreatedBy,
		&session.CreatedAt,
		&session.UpdatedBy,
		&session.UpdatedAt,
	)

	if err != nil {
		return SessionModel{}, err
	}

	return session, nil

}
