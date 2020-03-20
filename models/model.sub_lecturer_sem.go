package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"time"
)

type (
	SubLecturerSem struct {
		ID         uuid.UUID
		SubjectID  uuid.UUID
		LecturerID uuid.UUID
		SemesterID uuid.UUID
		IsDelete   bool
		CreatedBy  uuid.UUID
		CreatedAt  time.Time
		UpdatedBy  uuid.NullUUID
		UpdatedAt  pq.NullTime
	}

	SubLecturerSemResponse struct {
		ID        uuid.UUID
		Lecturer  LecturerResponse
		Semester  SemesterResponse
		Subject   SubjectResponse
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.UUID
		UpdatedAt time.Time
	}
)

func (s SubLecturerSem) Response() SubLecturerSemResponse {
	return SubLecturerSemResponse{
		ID:        s.ID,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}
}

func GetAllLecturerSemester(ctx context.Context, db *sql.DB) ([]SubLecturerSem, error) {
	query := fmt.Sprintf(`
	SELECT
			id,
			subject_id,
			lecturer_id,
			semester_id,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM sub_lecturer_sem
	`)
	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subLecturerSems []SubLecturerSem
	for rows.Next() {
		var subLecturerSem SubLecturerSem
		err := rows.Scan(
			&subLecturerSem.ID,
			&subLecturerSem.SubjectID,
			&subLecturerSem.LecturerID,
			&subLecturerSem.SemesterID,
			&subLecturerSem.IsDelete,
			&subLecturerSem.CreatedBy,
			&subLecturerSem.CreatedAt,
			&subLecturerSem.UpdatedBy,
			&subLecturerSem.UpdatedAt,
		)
		if err != nil {
			return nil, err

		}
		subLecturerSems = append(subLecturerSems, subLecturerSem)
	}

	return subLecturerSems, nil

}
