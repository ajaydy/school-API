package api

import (
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
	"school/models"
	"time"
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

	SessionAddParam struct {
		SubjectID   uuid.UUID `json:"subject_id"`
		LecturerID  uuid.UUID `json:"lecturer_id"`
		ClassroomID uuid.UUID `json:"classroom_id"`
		IntakeID    uuid.UUID `json:"intake_id"`
		ProgramID   uuid.UUID `json:"program_id"`
		Day         int       `json:"day"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
	}

	SessionUpdateParam struct {
		ID          uuid.UUID `json:"id"`
		SubjectID   uuid.UUID `json:"subject_id"`
		LecturerID  uuid.UUID `json:"lecturer_id"`
		ClassroomID uuid.UUID `json:"classroom_id"`
		IntakeID    uuid.UUID `json:"intake_id"`
		ProgramID   uuid.UUID `json:"program_id"`
		Day         int       `json:"day"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
	}

	SessionDeleteParam struct {
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

func (s SessionModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	sessions, err := models.GetAllSession(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var sessionsResponse []models.SessionResponse
	for _, session := range sessions {
		response, err := session.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		sessionsResponse = append(sessionsResponse, response)
	}

	return sessionsResponse, nil
}

func (s SessionModule) SessionListByLecturer(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	lecturerID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	sessions, err := models.GetAllSessionByLecturer(ctx, s.db, filter, lecturerID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "SessionListByLecturer/GetAllSessionByLecturer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var sessionResponse []models.SessionResponse
	for _, session := range sessions {
		response, err := session.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "SessionListByLecturer/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		sessionResponse = append(sessionResponse, response)
	}

	return sessionResponse, nil
}

func (s SessionModule) Add(ctx context.Context, param SessionAddParam) (interface{}, *helpers.Error) {

	session := models.SessionModel{

		SubjectID:   param.SubjectID,
		LecturerID:  param.LecturerID,
		IntakeID:    param.IntakeID,
		ClassroomID: param.ClassroomID,
		ProgramID:   param.ProgramID,
		Day:         param.Day,
		StartTime:   param.StartTime,
		EndTime:     param.EndTime,
		CreatedBy:   uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := session.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := session.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s SessionModule) Update(ctx context.Context, param SessionUpdateParam) (interface{}, *helpers.Error) {

	session := models.SessionModel{
		ID:          param.ID,
		SubjectID:   param.SubjectID,
		LecturerID:  param.LecturerID,
		IntakeID:    param.IntakeID,
		ClassroomID: param.ClassroomID,
		ProgramID:   param.ProgramID,
		Day:         param.Day,
		StartTime:   param.StartTime,
		EndTime:     param.EndTime,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := session.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := session.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (s SessionModule) Delete(ctx context.Context, param SessionDeleteParam) (interface{}, *helpers.Error) {

	session := models.SessionModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := session.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
