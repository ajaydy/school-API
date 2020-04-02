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
	ResultModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	ResultDetailParam struct {
		ID uuid.UUID `json:"id"`
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

	studentID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	result, err := models.GetAllResultForOneStudent(ctx, s.db, filter, studentID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllResult", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var resultsResponse []models.ResultResponse
	for _, results := range result {
		response, err := results.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/ResultResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		resultsResponse = append(resultsResponse, response)
	}

	return resultsResponse, nil
}
