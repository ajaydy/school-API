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
	FacultyModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	FacultyDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	FacultyAddParam struct {
		Code         int    `json:"code" valid:"required"`
		Abbreviation string `json:"abbreviation" valid:"required"`
		Name         string `json:"name" valid:"required"`
		Description  string `json:"description" valid:"required"`
	}

	FacultyUpdateParam struct {
		ID           uuid.UUID `json:"id" valid:"required"`
		Code         int       `json:"code" valid:"required"`
		Abbreviation string    `json:"abbreviation" valid:"required"`
		Name         string    `json:"name" valid:"required"`
		Description  string    `json:"description" valid:"required"`
	}

	FacultyDeleteParam struct {
		ID uuid.UUID `json:"id" valid:"required"`
	}
)

func NewFacultyModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *FacultyModule {
	return &FacultyModule{
		db:     db,
		cache:  cache,
		name:   "module/faculty",
		logger: logger,
	}
}

func (s FacultyModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	faculties, err := models.GetAllFaculty(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllFaculty", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var facultiesResponse []models.FacultyResponse
	for _, faculty := range faculties {
		facultiesResponse = append(facultiesResponse, faculty.Response())
	}

	return facultiesResponse, nil
}

func (s FacultyModule) Detail(ctx context.Context, param FacultyDetailParam) (interface{}, *helpers.Error) {
	faculty, err := models.GetOneFaculty(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneFaculty", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return faculty.Response(), nil
}

func (s FacultyModule) Add(ctx context.Context, param FacultyAddParam) (interface{}, *helpers.Error) {

	faculty := models.FacultyModel{
		Code:         param.Code,
		Abbreviation: param.Abbreviation,
		Name:         param.Name,
		Description:  param.Description,
		CreatedBy:    uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := faculty.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return faculty.Response(), nil
}

func (s FacultyModule) Update(ctx context.Context, param FacultyUpdateParam) (interface{}, *helpers.Error) {

	faculty := models.FacultyModel{
		ID:           param.ID,
		Code:         param.Code,
		Abbreviation: param.Abbreviation,
		Name:         param.Name,
		Description:  param.Description,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := faculty.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return faculty.Response(), nil

}

func (s FacultyModule) Delete(ctx context.Context, param FacultyDeleteParam) (interface{}, *helpers.Error) {

	faculty := models.FacultyModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := faculty.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
