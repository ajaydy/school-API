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
	AttendanceModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	AttendanceDetailParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewAttendanceModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *AttendanceModule {
	return &AttendanceModule{
		db:     db,
		cache:  cache,
		name:   "module/attendance",
		logger: logger,
	}
}

func (s AttendanceModule) Detail(ctx context.Context, param AttendanceDetailParam) (interface{}, *helpers.Error) {
	attendance, err := models.GetOneAttendance(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneAttendance", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := attendance.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/AttendanceResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}
