package api

import (
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
	"school/models"
	"time"
)

type (
	ClassModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	ClassDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	ClassAddParam struct {
		SessionID uuid.UUID `json:"session_id"`
		Date      time.Time `json:"date"`
	}
)

func NewClassModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *ClassModule {
	return &ClassModule{
		db:     db,
		cache:  cache,
		name:   "module/class",
		logger: logger,
	}
}

func (s ClassModule) Detail(ctx context.Context, param ClassDetailParam) (interface{}, *helpers.Error) {
	class, err := models.GetOneClass(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneClass", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := class.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/ClassResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s ClassModule) ListBySession(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	class, err := models.GetAllClassBySession(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListBySession/GetAllClassBySession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var classResponse []models.ClassResponse
	for _, classes := range class {
		response, err := classes.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListBySession/ClassResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		classResponse = append(classResponse, response)
	}

	return classResponse, nil
}

func (s ClassModule) Add(ctx context.Context, param ClassAddParam) (interface{}, *helpers.Error) {

	class := models.ClassModel{
		SessionID: param.SessionID,
		Date:      param.Date,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := class.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "AddClass/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	students, err := models.GetAllStudentEnrollBySession(ctx, s.db, helpers.Filter{
		FilterOption: helpers.FilterOption{
			Limit:  999,
			Offset: 0,
			Dir:    "asc",
		},
		SessionID: param.SessionID,
	})

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetAllStudentEnrollBySession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	for _, student := range students {
		attendance := models.AttendanceModel{
			StudentID: student.StudentID,
			ClassID:   class.ID,
			CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		}

		err = attendance.Insert(ctx, s.db)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/AttendanceInsert", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
	}

	response, err := class.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}
