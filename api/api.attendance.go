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

	AttendanceAddParam struct {
		ClassID   uuid.UUID `json:"class_id"`
		StudentID uuid.UUID `json:"student_id"`
	}

	AttendanceUpdateParam struct {
		ID       uuid.UUID `json:"id"`
		IsAttend bool      `json:"is_attend"`
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

func (s AttendanceModule) ListByClass(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	attendance, err := models.GetAllAttendanceByClass(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListByClass/GetAllAttendanceBySession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var attendanceResponse []models.AttendanceResponse
	for _, attendances := range attendance {
		response, err := attendances.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListByClass/ClassResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		attendanceResponse = append(attendanceResponse, response)
	}

	return attendanceResponse, nil
}

func (s AttendanceModule) Add(ctx context.Context, param AttendanceAddParam) (interface{}, *helpers.Error) {

	attendance := models.AttendanceModel{
		StudentID: param.StudentID,
		ClassID:   param.ClassID,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := attendance.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := attendance.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s AttendanceModule) UpdateIsAttend(ctx context.Context, param AttendanceUpdateParam) (interface{}, *helpers.Error) {

	attendance := models.AttendanceModel{
		ID:       param.ID,
		IsAttend: param.IsAttend,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := attendance.UpdateIsAttend(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "UpdateIsAttend/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	attendances, err := attendance.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "UpdateIsAttend/AttendanceResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return attendances, nil

}
