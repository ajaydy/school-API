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
