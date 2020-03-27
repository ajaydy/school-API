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
	IntakeModel struct {
		ID        uuid.UUID
		Year      string
		Month     int
		Trimester int
		StartDate time.Time
		EndDate   time.Time
		IsDelete  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	IntakeResponse struct {
		ID        uuid.UUID `json:"id"`
		Year      string    `json:"year"`
		Month     int       `json:"month"`
		Trimester int       `json:"trimester"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		IsDelete  bool      `json:"is_delete"`
		CreatedBy uuid.UUID `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedBy uuid.UUID `json:"updated_by"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func (s IntakeModel) Response() IntakeResponse {
	return IntakeResponse{
		ID:        s.ID,
		Year:      s.Year,
		Month:     s.Month,
		Trimester: s.Trimester,
		StartDate: s.StartDate,
		EndDate:   s.EndDate,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}
}

func GetOneIntake(ctx context.Context, db *sql.DB, intakeID uuid.UUID) (IntakeModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			year,
			month,
			trimester,
			start_date,
			end_date,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM intake
		WHERE 
			id = $1
	`)

	var intake IntakeModel
	err := db.QueryRowContext(ctx, query, intakeID).Scan(
		&intake.ID,
		&intake.Year,
		&intake.Month,
		&intake.Trimester,
		&intake.StartDate,
		&intake.EndDate,
		&intake.IsDelete,
		&intake.CreatedBy,
		&intake.CreatedAt,
		&intake.UpdatedBy,
		&intake.UpdatedAt,
	)

	if err != nil {
		return IntakeModel{}, err
	}

	return intake, nil

}

//func GetAllIntake(ctx context.Context, db *sql.DB) ([]IntakeModel, error) {
//
//	query := fmt.Sprintf(`
//		SELECT
//			id,
//			year,
//			month,
//			is_delete,
//			created_by,
//			created_at,
//			updated_by,
//			updated_at
//		FROM intake`)
//
//	rows, err := db.QueryContext(ctx, query)
//
//	if err != nil {
//		return nil, err
//	}
//
//	defer rows.Close()
//
//	var intakes []IntakeModel
//	for rows.Next() {
//		var intake IntakeModel
//		rows.Scan(
//			&intake.ID,
//			&intake.Year,
//			&intake.Month,
//			&intake.IsDelete,
//			&intake.CreatedBy,
//			&intake.CreatedAt,
//			&intake.UpdatedBy,
//			&intake.UpdatedAt,
//		)
//
//		intakes = append(intakes, intake)
//	}
//
//	return intakes, nil
//
//}
//
//func (s *IntakeModel) Insert(ctx context.Context, db *sql.DB) error {
//
//	query := fmt.Sprintf(`
//		INSERT INTO intake(
//			year,
//			month,
//			created_by,
//			created_at)
//		VALUES(
//		$1,$2,$3,now())
//		RETURNING id, created_at`)
//
//	err := db.QueryRowContext(ctx, query,
//		s.Year, s.Month, s.CreatedBy).Scan(
//		&s.ID, &s.CreatedAt,
//	)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}
//
//func (s *IntakeModel) Update(ctx context.Context, db *sql.DB) error {
//
//	query := fmt.Sprintf(`
//		UPDATE intake
//		SET
//			year=$1,
//			month=$2,
//			updated_at=NOW(),
//			updated_by=$3
//		WHERE id=$4
//		RETURNING id,created_at`)
//
//	err := db.QueryRowContext(ctx, query,
//		s.Year, s.Month, s.UpdatedBy, s.ID).Scan(
//		&s.ID, &s.CreatedAt,
//	)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//
//}
