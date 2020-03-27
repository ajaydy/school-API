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
	ProgramModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	ProgramDetailParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewProgramModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *ProgramModule {
	return &ProgramModule{
		db:     db,
		cache:  cache,
		name:   "module/program",
		logger: logger,
	}
}

func (s ProgramModule) Detail(ctx context.Context, param ProgramDetailParam) (interface{}, *helpers.Error) {
	program, err := models.GetOneProgram(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneProgram", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := program.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/ProgramResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}
