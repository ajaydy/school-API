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
	SessionModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	SessionDetailParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewSessionModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *SessionModule {
	return &SessionModule{
		db:     db,
		cache:  cache,
		name:   "module/session",
		logger: logger,
	}
}

func (s SessionModule) Detail(ctx context.Context, param SessionDetailParam) (interface{}, *helpers.Error) {
	session, err := models.GetOneSession(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneSession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := session.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/SessionResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s SessionModule) List(ctx context.Context) (interface{}, *helpers.Error) {
	sessions, err := models.GetAllSession(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var sessionResponse []models.SessionResponse
	for _, session := range sessions {
		response, err := session.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSessionResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		sessionResponse = append(sessionResponse, response)
	}

	return sessionResponse, nil
}
