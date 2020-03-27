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
	LecturersModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}
	LecturerDetailParam struct {
		ID uuid.UUID `json:"id"`
	}
	//LecturerParamAdd struct {
	//	Name    string `json:"name" valid:"length(3|50),required"`
	//	Address string `json:"address" valid:"optional"`
	//	PhoneNo string `json:"phone_no" valid:"length(10|15),required"`
	//	Email   string `json:"email" valid:"email,required"`
	//}
	//
	//LecturerParamUpdate struct {
	//	ID      uuid.UUID `json:"id"`
	//	Name    string    `json:"name" valid:"length(3|50),required"`
	//	Address string    `json:"address" valid:"optional"`
	//	PhoneNo string    `json:"phone_no" valid:"length(10|15),required"`
	//	Email   string    `json:"email" valid:"email,required"`
	//}
)

func NewLecturersModule(db *sql.DB, cache *redis.Pool) *LecturersModule {
	return &LecturersModule{
		db:    db,
		cache: cache,
		name:  "module/lecturers",
	}
}

//func (s LecturersModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
//	lecturers, err := models.GetAllLecturer(ctx, s.db, filter)
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllLecturer", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	var lecturerResponse []models.LecturerResponse
//	for _, lecturer := range lecturers {
//		lecturerResponse = append(lecturerResponse, lecturer.Response())
//	}
//
//	return lecturerResponse, nil
//}

func (s LecturersModule) Detail(ctx context.Context, studentID uuid.UUID) (interface{}, *helpers.Error) {
	lecturer, err := models.GetOneLecturer(ctx, s.db, studentID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneLecturer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return lecturer.Response(), nil
}

//func (s LecturersModule) Add(ctx context.Context, param LecturerParamAdd) (interface{}, *helpers.Error) {
//	lecturers := models.LecturerModel{
//		Name:      param.Name,
//		Address:   param.Address,
//		Email:     param.Email,
//		PhoneNo:   param.PhoneNo,
//		CreatedBy: uuid.NewV4(),
//	}
//
//	err := lecturers.Insert(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return lecturers.Response(), nil
//}
//
//func (s LecturersModule) Update(ctx context.Context, param LecturerParamUpdate) (interface{}, *helpers.Error) {
//
//	lecturer := models.LecturerModel{
//		ID:      param.ID,
//		Name:    param.Name,
//		Address: param.Address,
//		Email:   param.Email,
//		PhoneNo: param.PhoneNo,
//		UpdatedBy: uuid.NullUUID{
//			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
//			Valid: true,
//		},
//	}
//	err := lecturer.Update(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return lecturer.Response(), nil
//
//}
