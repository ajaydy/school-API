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
	StudentsModule struct {
		db    *sql.DB
		cache *redis.Pool
		name  string
	}
)

func NewStudentsModule(db *sql.DB, cache *redis.Pool) *StudentsModule {
	return &StudentsModule{
		db:    db,
		cache: cache,
		name:  "module/students",
	}
}
func (s StudentsModule) List(ctx context.Context) (interface{}, *helpers.Error) {
	students, err := models.GetAllStudent(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllStudent", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var studentsResponse []models.StudentResponse
	for _, student := range students {
		studentsResponse = append(studentsResponse, student.Response())
	}

	return studentsResponse, nil
}

func (s StudentsModule) Detail(ctx context.Context, studentID uuid.UUID) (interface{}, *helpers.Error) {
	student, err := models.GetOneStudent(ctx, s.db, studentID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneStudent", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return student.Response(), nil
}
