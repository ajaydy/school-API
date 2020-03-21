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
)

func NewSubjectsModule(db *sql.DB, cache *redis.Pool) *SubjectModule {
	return &SubjectModule{
		db:    db,
		cache: cache,
		name:  "module/subjects",
	}
}
func (s SubjectModule) List(ctx context.Context) (interface{}, *helpers.Error) {
	subjects, err := models.GetAllSubject(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSubject", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var subjectResponse []models.SubjectResponse
	for _, subject := range subjects {
		subjectResponse = append(subjectResponse, subject.Response())
	}

	return subjectResponse, nil
}

func (s SubjectModule) Detail(ctx context.Context, subjectID uuid.UUID) (interface{}, *helpers.Error) {
	subject, err := models.GetOneSubject(ctx, s.db, subjectID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneSubject", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return subject.Response(), nil
}
