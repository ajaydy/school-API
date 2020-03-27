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
	ClassroomModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	ClassroomDetailParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewClassroomModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *ClassroomModule {
	return &ClassroomModule{
		db:     db,
		cache:  cache,
		name:   "module/classroom",
		logger: logger,
	}
}

func (s ClassroomModule) Detail(ctx context.Context, param ClassroomDetailParam) (interface{}, *helpers.Error) {
	classroom, err := models.GetOneClassroom(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneClassroom", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := classroom.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/ClassroomResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}
