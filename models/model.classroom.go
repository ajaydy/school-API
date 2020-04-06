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

func GetAllClassroom(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ClassRoomModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`AND LOWER(code) LIKE LOWER('%%%s%%')`, filter.Search)
	}

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
		WHERE is_delete = false 
		%s
		ORDER BY floor %s,room_no %s
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var classrooms []ClassRoomModel
	for rows.Next() {
		var classroom ClassRoomModel

		rows.Scan(
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

		classrooms = append(classrooms, classroom)
	}

	return classrooms, nil

}

func (s *ClassRoomModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO classroom(
			faculty_id,
			floor,
			room_no,
			code,
			created_by,
			created_at)
		VALUES(
		$1,$2,$3,$4,$5,now())
		RETURNING id, created_at,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.FacultyID, s.Floor, s.RoomNo, s.Code, s.CreatedBy).Scan(
		&s.ID, &s.CreatedAt, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ClassRoomModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE classroom
		SET
			faculty_id=$1,
			floor=$2,
			room_no=$3,
			code=$4,
			updated_at=NOW(),
			updated_by=$5
		WHERE id=$6
		RETURNING id,created_at,updated_at,created_by,is_delete`)

	err := db.QueryRowContext(ctx, query,
		s.FacultyID, s.Floor, s.RoomNo, s.Code, s.UpdatedBy, s.ID).Scan(
		&s.ID, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.IsDelete,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *ClassRoomModel) Delete(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE classroom
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
