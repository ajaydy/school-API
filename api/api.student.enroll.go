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

	StudentEnrollParamAdd struct {
		SessionID uuid.UUID `json:"session_id" value:"required"`
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

func (s StudentEnrollModule) Add(ctx context.Context, param StudentEnrollParamAdd) (interface{}, *helpers.Error) {

	student := models.StudentEnrollModel{
		SessionID: param.SessionID,
		StudentID: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := student.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := student.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s StudentEnrollModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	studentID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	student, err := models.GetTimetableForStudent(ctx, s.db, filter, studentID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllTimetable", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var studentResponse []models.StudentEnrollResponse
	for _, students := range student {
		response, err := students.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/TimetableResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		studentResponse = append(studentResponse, response)
	}

	return studentResponse, nil
}
