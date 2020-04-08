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

	AttendanceListByClassParam struct {
		ClassID uuid.UUID `json:"class_id"`
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

func (s AttendanceModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	attendances, err := models.GetAllAttendance(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllAttendance", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var attendancesResponse []models.AttendanceResponse
	for _, attendance := range attendances {
		response, err := attendance.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/AttendanceResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		attendancesResponse = append(attendancesResponse, response)
	}

	return attendancesResponse, nil
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
		return nil, helpers.ErrorWrap(err, s.name, "Add/AttendanceResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s AttendanceModule) Update(ctx context.Context, param AttendanceUpdateParam) (interface{}, *helpers.Error) {

	attendance := models.AttendanceModel{
		ID:       param.ID,
		IsAttend: param.IsAttend,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := attendance.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	attendances, err := attendance.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/AttendanceResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return attendances, nil

}

func (s AttendanceModule) ListByClass(ctx context.Context, filter helpers.Filter, param AttendanceListByClassParam) (
	interface{}, *helpers.Error) {

	attendances, err := models.GetAllAttendanceByClass(ctx, s.db, helpers.Filter{
		FilterOption: helpers.FilterOption{
			Limit:  999,
			Offset: 0,
		},

		ClassID: param.ClassID,
	})
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListByClass/GetAllAttendanceByClass", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var attendancesResponse []models.AttendanceResponse
	for _, attendance := range attendances {
		response, err := attendance.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListByClass/AttendanceResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		attendancesResponse = append(attendancesResponse, response)
	}

	return attendancesResponse, nil
}
