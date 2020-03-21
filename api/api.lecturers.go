package api

import (
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
	"school/models"
)

type (
	LecturersModule struct {
		db    *sql.DB
		cache *redis.Pool
		name  string
	}
)

func NewLecturersModule(db *sql.DB, cache *redis.Pool) *LecturersModule {
	return &LecturersModule{
		db:    db,
		cache: cache,
		name:  "module/lecturers",
	}
}

func (s LecturersModule) List(ctx context.Context) (interface{}, *helpers.Error) {
	lecturers, err := models.GetAllLecturer(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllLecturer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var lecturerResponse []models.LecturerResponse
	for _, lecturer := range lecturers {
		lecturerResponse = append(lecturerResponse, lecturer.Response())
	}

	return lecturerResponse, nil
}

func (s LecturersModule) Detail(ctx context.Context, studentID uuid.UUID) (interface{}, *helpers.Error) {
	lecturer, err := models.GetOneLecturer(ctx, s.db, studentID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneLecturer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return lecturer.Response(), nil
}
