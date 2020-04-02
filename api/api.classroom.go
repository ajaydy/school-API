package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
	"school/models"
)

type (
	ClassroomModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	ClassroomDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	ClassroomAddParam struct {
		FacultyID uuid.UUID `json:"faculty_id"valid:"required"`
		Floor     int       `json:"floor"valid:"required"`
		RoomNo    int       `json:"room_no"valid:"required"`
	}

	ClassroomUpdateParam struct {
		ID        uuid.UUID `json:"id"`
		FacultyID uuid.UUID `json:"faculty_id"valid:"required"`
		Floor     int       `json:"floor"valid:"required"`
		RoomNo    int       `json:"room_no"valid:"required"`
	}

	ClassroomDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewClassroomModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *ClassroomModule {
	return &ClassroomModule{
		db:     db,
		cache:  cache,
		name:   "module/classroom",
		logger: logger,
	}
}

func (s ClassroomModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	classrooms, err := models.GetAllClassroom(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllClassroom", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var classroomsResponse []models.ClassRoomResponse
	for _, classroom := range classrooms {
		response, err := classroom.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/ClassroomResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		classroomsResponse = append(classroomsResponse, response)
	}

	return classroomsResponse, nil
}

func (s ClassroomModule) Detail(ctx context.Context, param ClassroomDetailParam) (interface{}, *helpers.Error) {
	classroom, err := models.GetOneClassroom(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneClassroom", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := classroom.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/ClassroomResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s ClassroomModule) Add(ctx context.Context, param ClassroomAddParam) (interface{}, *helpers.Error) {

	faculty, err := models.GetOneFaculty(ctx, s.db, param.FacultyID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneFaculty", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	roomNo := fmt.Sprintf("%03d", param.RoomNo)

	facultyAbbreviation := faculty.Abbreviation

	roomCode := fmt.Sprintf("%s%d%s", facultyAbbreviation, param.Floor, roomNo)

	classroom := models.ClassRoomModel{
		FacultyID: param.FacultyID,
		Floor:     param.Floor,
		RoomNo:    param.RoomNo,
		Code:      roomCode,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err = classroom.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := classroom.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s ClassroomModule) Update(ctx context.Context, param ClassroomUpdateParam) (interface{}, *helpers.Error) {

	faculty, err := models.GetOneFaculty(ctx, s.db, param.FacultyID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOneFaculty", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	roomNo := fmt.Sprintf("%03d", param.RoomNo)

	facultyAbbreviation := faculty.Abbreviation

	roomCode := fmt.Sprintf("%s%d%s", facultyAbbreviation, param.Floor, roomNo)

	classroom := models.ClassRoomModel{
		ID:        param.ID,
		FacultyID: param.FacultyID,
		Floor:     param.Floor,
		RoomNo:    param.RoomNo,
		Code:      roomCode,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err = classroom.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := classroom.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (s ClassroomModule) Delete(ctx context.Context, param ClassroomDeleteParam) (interface{}, *helpers.Error) {

	classroom := models.ClassRoomModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := classroom.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
