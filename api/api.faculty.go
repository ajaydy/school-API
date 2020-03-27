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
	FacultyModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	FacultyDetailParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewFacultyModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *FacultyModule {
	return &FacultyModule{
		db:     db,
		cache:  cache,
		name:   "module/faculty",
		logger: logger,
	}
}

func (s FacultyModule) Detail(ctx context.Context, param FacultyDetailParam) (interface{}, *helpers.Error) {
	faculty, err := models.GetOneFaculty(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneFaculty", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return faculty.Response(), nil
}
