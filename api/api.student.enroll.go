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
	StudentEnrollModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	StudentEnrollDetailParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewStudentEnrollModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *StudentEnrollModule {
	return &StudentEnrollModule{
		db:     db,
		cache:  cache,
		name:   "module/studentEnroll",
		logger: logger,
	}

}

func (s StudentEnrollModule) Detail(ctx context.Context, param StudentEnrollDetailParam) (interface{}, *helpers.Error) {
	student, err := models.GetOneStudentEnroll(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneStudentEnroll", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := student.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/StudentEnrollResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}
