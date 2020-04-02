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
	SubjectModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	SubjectDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	SubjectAddParam struct {
		Name        string `json:"name" valid:"required"`
		Description string `json:"description" valid:"required"`
		Duration    int    `json:"duration" valid:"required"`
	}

	SubjectUpdateParam struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name" valid:"required"`
		Description string    `json:"description" valid:"required"`
		Duration    int       `json:"duration" valid:"required"`
	}

	SubjectDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewSubjectModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *SubjectModule {
	return &SubjectModule{
		db:     db,
		cache:  cache,
		name:   "module/subject",
		logger: logger,
	}
}

func (s SubjectModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	subjects, err := models.GetAllSubject(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSubject", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var subjectsResponse []models.SubjectResponse
	for _, subject := range subjects {
		subjectsResponse = append(subjectsResponse, subject.Response())
	}

	return subjectsResponse, nil
}

func (s SubjectModule) Detail(ctx context.Context, param SubjectDetailParam) (interface{}, *helpers.Error) {
	subject, err := models.GetOneSubject(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneSubject", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return subject.Response(), nil
}

func (s SubjectModule) Add(ctx context.Context, param SubjectAddParam) (interface{}, *helpers.Error) {
	subject := models.SubjectModel{
		Name:        param.Name,
		Description: param.Description,
		Duration:    param.Duration,
		CreatedBy:   uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := subject.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return subject.Response(), nil
}

func (s SubjectModule) Update(ctx context.Context, param SubjectUpdateParam) (interface{}, *helpers.Error) {

	subject := models.SubjectModel{
		ID:          param.ID,
		Name:        param.Name,
		Description: param.Description,
		Duration:    param.Duration,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := subject.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return subject.Response(), nil

}

func (s SubjectModule) Delete(ctx context.Context, param SubjectDeleteParam) (interface{}, *helpers.Error) {

	subject := models.SubjectModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := subject.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
