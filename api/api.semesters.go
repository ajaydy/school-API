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
	SemestersModule struct {
		db    *sql.DB
		cache *redis.Pool
		name  string
	}
)

func NewSemestersModule(db *sql.DB, cache *redis.Pool) *SemestersModule {
	return &SemestersModule{
		db:    db,
		cache: cache,
		name:  "module/students",
	}
}
func (s SemestersModule) List(ctx context.Context) (interface{}, *helpers.Error) {
	semesters, err := models.GetAllSemester(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSemester", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var semesterResponse []models.SemesterResponse
	for _, semester := range semesters {
		semesterResponse = append(semesterResponse, semester.Response())
	}

	return semesterResponse, nil
}

func (s SemestersModule) Detail(ctx context.Context, semesterID uuid.UUID) (interface{}, *helpers.Error) {
	semester, err := models.GetOneSemester(ctx, s.db, semesterID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneSemester", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return semester.Response(), nil
}
