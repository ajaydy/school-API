package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
	"school/models"
	"time"
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

	StudentEnrollAddParam struct {
		SessionID uuid.UUID `json:"session_id" value:"required"`
	}

	StudentEnrollDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}

	StudentEnrollListBySessionParam struct {
		SessionID uuid.UUID `json:"session_id"`
	}

	StudentEnrollListByStudentParam struct {
		StudentID uuid.UUID `json:"student_id"`
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

func (s StudentEnrollModule) Add(ctx context.Context, param StudentEnrollAddParam) (interface{}, *helpers.Error) {

	studentID := uuid.FromStringOrNil(ctx.Value("user_id").(string))
	student, err := models.GetOneStudent(ctx, s.db, studentID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneStudent", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	studentProgramID := student.ProgramID

	session, err := models.GetOneSession(ctx, s.db, param.SessionID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneSession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	studentEnroll, err := models.GetOneStudentEnrollBySessionAndStudentID(ctx, s.db, param.SessionID, studentID)
	if err != nil && err != sql.ErrNoRows {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneStudentEnrollBySessionAndStudentID",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	if studentEnroll.SessionID == param.SessionID {
		return nil, helpers.ErrorWrap(errors.New("You have already enroll this session"), s.name,
			"Add/ValidationSession",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	sessionProgramID := session.ProgramID

	fmt.Println(sessionProgramID)
	fmt.Println(studentProgramID)

	intake, err := models.GetOneIntake(ctx, s.db, session.IntakeID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneIntake", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	startDate := intake.StartDate

	enrollDateStart := startDate.AddDate(0, 0, -5)
	enrollDateEnd := startDate.Add(-1 * time.Second)

	//fmt.Println(enrollDateStart)
	//fmt.Println(enrollDateEnd)

	now := time.Now()
	//)fmt.Println(now)

	if now.After(enrollDateEnd) || now.Before(enrollDateStart) {
		return nil, helpers.ErrorWrap(errors.New("Invalid Time To Enroll"), s.name, "Add/ValidationDate",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	if sessionProgramID != studentProgramID {
		return nil, helpers.ErrorWrap(errors.New("This Session Is Not For Your Program"), s.name,
			"Add/ValidationProgram",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}
	////fmt.Println(now.Before(enrollDateEnd) && now.After(enrollDateStart))
	studentEnroll = models.StudentEnrollModel{
		SessionID: param.SessionID,
		StudentID: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err = studentEnroll.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/StudentEnrollInsert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	result := models.ResultModel{
		StudentEnrollID: studentEnroll.ID,
		CreatedBy:       uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err = result.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/ResultInsert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := studentEnroll.Response(ctx, s.db, s.logger)
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

func (s StudentEnrollModule) ListBySession(ctx context.Context, filter helpers.Filter,
	param StudentEnrollListBySessionParam) (
	interface{}, *helpers.Error) {

	studentEnrolls, err := models.GetAllStudentEnrollBySession(ctx, s.db, helpers.Filter{
		FilterOption: helpers.FilterOption{
			Limit:  999,
			Offset: 0,
		},

		SessionID: param.SessionID,
	})

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListBySession/GetAllStudentEnrollBySession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var studentEnrollsResponse []models.StudentEnrollResponse
	for _, studentEnroll := range studentEnrolls {
		response, err := studentEnroll.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListBySession/StudentEnrollResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		studentEnrollsResponse = append(studentEnrollsResponse, response)
	}

	return studentEnrollsResponse, nil
}

//func (s StudentEnrollModule) ListByStudent(ctx context.Context, filter helpers.Filter,
//	param StudentEnrollListByStudentParam) (
//	interface{}, *helpers.Error) {
//
//	studentEnrolls, err := models.GetAllStudentEnrollByStudent(ctx, s.db, helpers.Filter{
//		FilterOption: helpers.FilterOption{
//			Limit:  999,
//			Offset: 0,
//		},
//
//		StudentID: param.StudentID,
//	})
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "ListByStudent/GetAllStudentEnrollByStudent", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	var studentEnrollsResponse []models.StudentEnrollResponse
//	for _, studentEnroll := range studentEnrolls {
//		response, err := studentEnroll.Response(ctx, s.db, s.logger)
//		if err != nil {
//			return nil, helpers.ErrorWrap(err, s.name, "ListByStudent/StudentEnrollResponse", helpers.InternalServerError,
//				http.StatusInternalServerError)
//		}
//		studentEnrollsResponse = append(studentEnrollsResponse, response)
//	}
//
//	return studentEnrollsResponse, nil
//}

func (s StudentEnrollModule) Delete(ctx context.Context, param StudentEnrollDeleteParam) (interface{}, *helpers.Error) {

	studentEnroll := models.StudentEnrollModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := studentEnroll.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	results, err := models.GetAllResultByStudentEnroll(ctx, s.db, helpers.Filter{
		FilterOption: helpers.FilterOption{
			Limit:  999,
			Offset: 0,
			Dir:    "asc",
		},
		StudentEnrollID: param.ID,
	})

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/GetAllResultByStudentEnroll", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	for _, result := range results {
		result := models.ResultModel{
			ID: result.ID,
			UpdatedBy: uuid.NullUUID{
				UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
				Valid: true,
			},
		}

		err = result.Delete(ctx, s.db)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/ResultDelete", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
	}

	return nil, nil

}
