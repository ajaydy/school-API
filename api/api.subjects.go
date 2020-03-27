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
		db    *sql.DB
		cache *redis.Pool
		name  string
	}

	SubjectDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

//	SubjectParamAdd struct {
//		//ID          uuid.UUID `json:"id" valid:"uuid,required"`
//		Name        string `json:"name" valid:"length(3|50),required"`
//		Description string `json:"description" valid:"required"`
//		Duration    int    `json:"duration" valid:"required"`
//	}
//
//	SubjectParamUpdate struct {
//		ID          uuid.UUID `json:"id"`
//		Name        string    `json:"name" valid:"length(3|50),required"`
//		Description string    `json:"description" valid:"required"`
//		Duration    int       `json:"duration" valid:"required"`
//	}
)

func NewSubjectsModule(db *sql.DB, cache *redis.Pool) *SubjectModule {
	return &SubjectModule{
		db:    db,
		cache: cache,
		name:  "module/subjects",
	}
}

//func (s SubjectModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
//	subjects, err := models.GetAllSubject(ctx, s.db, filter)
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSubject", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	var subjectResponse []models.SubjectResponse
//	for _, subject := range subjects {
//		subjectResponse = append(subjectResponse, subject.Response())
//	}
//
//	return subjectResponse, nil
//}

func (s SubjectModule) Detail(ctx context.Context, param SubjectDetailParam) (interface{}, *helpers.Error) {
	subject, err := models.GetOneSubject(ctx, s.db, subjectID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneSubject", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return subject.Response(), nil
}

//func (s SubjectModule) Add(ctx context.Context, param SubjectParamAdd) (interface{}, *helpers.Error) {
//	subjects := models.SubjectModel{
//		Name:        param.Name,
//		Description: param.Description,
//		Duration:    param.Duration,
//		CreatedBy:   uuid.NewV4(),
//	}
//
//	err := subjects.Insert(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return subjects.Response(), nil
//}
//
//func (s SubjectModule) Update(ctx context.Context, param SubjectParamUpdate) (interface{}, *helpers.Error) {
//
//	subject := models.SubjectModel{
//		ID:          param.ID,
//		Name:        param.Name,
//		Description: param.Description,
//		Duration:    param.Duration,
//		UpdatedBy: uuid.NullUUID{
//			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
//			Valid: true,
//		},
//	}
//	err := subject.Update(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return subject.Response(), nil
//
//}
