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
		FacultyID uuid.UUID `json:"faculty_id"`
		Floor     int       `json:"floor"`
		RoomNo    int       `json:"room_no"`
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

	classrooms, err := classroom.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return classrooms, nil
}
