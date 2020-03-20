package models

import (
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"time"
)

type (
	StudSubLecturerSem struct {
		ID             uuid.UUID
		SubLecturerSem uuid.UUID
		StudentID      uuid.UUID
		Score          string
		IsDelete       bool
		CreatedBy      uuid.UUID
		CreatedAt      time.Time
		UpdatedBy      uuid.NullUUID
		UpdatedAt      pq.NullTime
	}

	StudSubLecturerSemResponse struct {
		ID             uuid.UUID
		SubLecturerSem SubLecturerSemResponse
		StudentID      StudentResponse
		Score          string
		IsDelete       bool
		CreatedBy      uuid.UUID
		CreatedAt      time.Time
		UpdatedBy      uuid.UUID
		UpdatedAt      time.Time
	}
)

func (s StudSubLecturerSem) Response() StudSubLecturerSemResponse {
	return StudSubLecturerSemResponse{
		ID:        s.ID,
		Score:     s.Score,
		IsDelete:  s.IsDelete,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}
}
