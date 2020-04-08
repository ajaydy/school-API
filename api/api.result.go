package api

import (
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
	"school/models"
	"school/util"
)

type (
	ResultModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	ResultDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	ResultUpdateParam struct {
		ID    uuid.UUID `json:"id"`
		Marks int       `json:"marks"`
	}

	ResultDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}

	ResultListByStudentEnrollParam struct {
		StudentEnrollID uuid.UUID `json:"student_enroll_id"`
	}
)

func NewResultModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *ResultModule {
	return &ResultModule{
		db:     db,
		cache:  cache,
		name:   "module/result",
		logger: logger,
	}
}

func (s ResultModule) Detail(ctx context.Context, param ResultDetailParam) (interface{}, *helpers.Error) {
	result, err := models.GetOneResult(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneResult", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := result.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/ResultResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s ResultModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	results, err := models.GetAllResult(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllResult", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var resultsResponse []models.ResultResponse
	for _, result := range results {
		response, err := result.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/ResultResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		resultsResponse = append(resultsResponse, response)
	}

	return resultsResponse, nil
}

func (s ResultModule) ListByOneStudent(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	studentID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	results, err := models.GetAllResultForOneStudent(ctx, s.db, helpers.Filter{
		FilterOption: helpers.FilterOption{
			Limit:  999,
			Offset: 0,
		},
		StudentID: studentID,
	})

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListByOneStudent/GetAllResultForOneStudent",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var resultsResponse []models.ResultResponse
	for _, result := range results {
		response, err := result.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListByOneStudent/ResultResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		resultsResponse = append(resultsResponse, response)
	}

	return resultsResponse, nil
}

//func (s ResultModule) Add(ctx context.Context, param ResultAddParam) (interface{}, *helpers.Error) {
//
//	grade := util.GetGrade(param.Marks)
//
//	result := models.ResultModel{
//		StudentEnrollID: param.StudentEnrollID,
//		Grade:           grade,
//		Marks:           param.Marks,
//		CreatedBy:       uuid.FromStringOrNil(ctx.Value("user_id").(string)),
//	}
//
//	err := result.Insert(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	response, err := result.Response(ctx, s.db, s.logger)
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return response, nil
//}

func (s ResultModule) Update(ctx context.Context, param ResultUpdateParam) (interface{}, *helpers.Error) {

	grade := util.GetGrade(param.Marks)

	result := models.ResultModel{
		ID:    param.ID,
		Grade: grade,
		Marks: param.Marks,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := result.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := result.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/ResultResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (s ResultModule) Delete(ctx context.Context, param ResultDeleteParam) (interface{}, *helpers.Error) {

	result := models.ResultModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := result.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}

func (s ResultModule) ListByStudentEnroll(ctx context.Context, filter helpers.Filter,
	param ResultListByStudentEnrollParam) (
	interface{}, *helpers.Error) {

	results, err := models.GetAllResultByStudentEnroll(ctx, s.db, helpers.Filter{
		FilterOption: helpers.FilterOption{
			Limit:  999,
			Offset: 0,
		},

		StudentEnrollID: param.StudentEnrollID,
	})

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListByStudentEnroll/GetAllResultByStudentEnroll",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var resultsResponse []models.ResultResponse
	for _, result := range results {
		response, err := result.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListByStudentEnroll/ResultResponse",
				helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		resultsResponse = append(resultsResponse, response)
	}

	return resultsResponse, nil
}
