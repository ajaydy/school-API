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
	ClassRoomModel struct {
		ID        uuid.UUID
		FacultyID uuid.UUID
		Floor     int
		RoomNo    int
		Code      string
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	ClassRoomResponse struct {
		ID        uuid.UUID       `json:"id"`
		Faculty   FacultyResponse `json:"faculty"`
		Floor     int             `json:"floor"`
		RoomNo    int             `json:"room_no"`
		Code      string          `json:"code"`
		IsDelete  bool            `json:"is_delete"`
		CreatedBy uuid.UUID       `json:"created_by"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedBy uuid.UUID       `json:"updated_by"`
		UpdatedAt time.Time       `json:"updated_at"`
	}
)

func (s ClassRoomModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (ClassRoomResponse, error) {

	faculty, err := GetOneFaculty(ctx, db, s.FacultyID)
	if err != nil {
		logger.Err.Printf(`model.classroom.go/GetOneFaculty/%v`, err)
		return ClassRoomResponse{}, nil
	}

	return ClassRoomResponse{
		ID:        s.ID,
		Faculty:   faculty.Response(),
		Floor:     s.Floor,
		RoomNo:    s.RoomNo,
		Code:      s.Code,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil
}

func GetOneClassroom(ctx context.Context, db *sql.DB, classroomID uuid.UUID) (ClassRoomModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			faculty_id,
			floor,
			room_no,
			code,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM classroom
		WHERE 
			id = $1
	`)

	var classroom ClassRoomModel
	err := db.QueryRowContext(ctx, query, classroomID).Scan(
		&classroom.ID,
		&classroom.FacultyID,
		&classroom.Floor,
		&classroom.RoomNo,
		&classroom.Code,
		&classroom.IsDelete,
		&classroom.CreatedBy,
		&classroom.CreatedAt,
		&classroom.UpdatedBy,
		&classroom.UpdatedAt,
	)

	if err != nil {
		return ClassRoomModel{}, err
	}

	return classroom, nil

}
