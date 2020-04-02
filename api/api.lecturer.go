package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"school/helpers"
	"school/models"
	"school/session"
	"school/util"
)

type (
	LecturerModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	LecturerDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	LecturerParamLogin struct {
		Email    string `json:"email" valid:"required"`
		Password string `json:"password" valid:"length(1|50),required"`
	}

	LecturerWithSession struct {
		Lecturer models.LecturerResponse `json:"lecturer"`
		Session  string                  `json:"session"`
	}

	LecturerAddResultParam struct {
		StudentEnrollID uuid.UUID `json:"student_enroll_id"`
		Marks           int       `json:"marks"`
	}

	LecturerUpdateResultParam struct {
		StudentEnrollID uuid.UUID `json:"student_enroll_id"`
		Marks           int       `json:"marks"`
	}

	LecturerParamAdd struct {
		Name      string    `json:"name" valid:"length(3|50),required"`
		ProgramID uuid.UUID `json:"program_id" valid:"required"`
		Address   string    `json:"address" valid:"optional"`
		Gender    int       `json:"gender" valid:"required"`
		PhoneNo   string    `json:"phone_no" valid:"length(0|15),required"`
		Email     string    `json:"email" valid:"email,required"`
	}

	//LecturerParamUpdate struct {
	//	ID      uuid.UUID `json:"id"`
	//	Name    string    `json:"name" valid:"length(3|50),required"`
	//	Address string    `json:"address" valid:"optional"`
	//	PhoneNo string    `json:"phone_no" valid:"length(10|15),required"`
	//	Email   string    `json:"email" valid:"email,required"`
	//}
)

func NewLecturerModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *LecturerModule {
	return &LecturerModule{
		db:     db,
		cache:  cache,
		name:   "module/lecturer",
		logger: logger,
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

func (s LecturerModule) Detail(ctx context.Context, param LecturerDetailParam) (interface{}, *helpers.Error) {
	lecturer, err := models.GetOneLecturer(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneLecturer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	lecturers, err := lecturer.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/LecturerResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return lecturers, nil
}

func (s LecturerModule) Add(ctx context.Context, param LecturerParamAdd) (interface{}, *helpers.Error) {

	password := util.RandomString(12)

	lecturer := models.LecturerModel{
		Name:      param.Name,
		ProgramID: param.ProgramID,
		Address:   param.Address,
		Email:     param.Email,
		Gender:    param.Gender,
		PhoneNo:   param.PhoneNo,
		Password:  password,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := lecturer.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	lecturers, err := lecturer.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return lecturers, nil
}

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

func (s LecturerModule) Login(ctx context.Context, param LecturerParamLogin) (interface{}, *helpers.Error) {

	lecturer, err := models.GetOneLecturerByEmail(ctx, s.db, param.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helpers.ErrorWrap(err, s.name, "Session/Login", helpers.IncorrectEmailMessage,
				http.StatusInternalServerError)
		}
		return nil, helpers.ErrorWrap(err, s.name, "List/GetOneLecturerByEmail", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(lecturer.Password), []byte(param.Password))
	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Invalid Password"), s.name, "Login/CompareHashAndPassword",
			helpers.IncorrectPasswordMessage,
			http.StatusInternalServerError)
	}

	session := session.Session{
		UserID:     lecturer.ID,
		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
		Expiry:     86400,
		Role:       session.LECTURER_ROLE,
	}

	err = session.Store(ctx)

	lecturers, err := lecturer.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Login/LecturerResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	lecturerSession := LecturerWithSession{
		Lecturer: lecturers,
		Session:  session.SessionKey,
	}

	return lecturerSession, nil

}

func (s LecturerModule) SessionListByLecturer(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

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
			return nil, helpers.ErrorWrap(err, s.name, "SessionListByLecturer/SessionResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		sessionResponse = append(sessionResponse, response)
	}

	return sessionResponse, nil
}

func (s LecturerModule) StudentEnrollListBySession(ctx context.Context, filter helpers.Filter) (
	interface{}, *helpers.Error) {

	enrolls, err := models.GetAllStudentEnrollBySession(ctx, s.db, filter)
	fmt.Println(filter)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "StudentEnrollListBySession/GetAllStudentEnrollBySession", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var enrollResponse []models.StudentEnrollResponse
	for _, enroll := range enrolls {
		response, err := enroll.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "StudentEnrollListBySession/StudentEnrollResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		enrollResponse = append(enrollResponse, response)
	}

	return enrollResponse, nil
}

func (s LecturerModule) LecturerAddResult(ctx context.Context, param LecturerAddResultParam) (interface{}, *helpers.Error) {

	score := util.GetGrade(param.Marks)

	result := models.ResultModel{
		StudentEnrollID: param.StudentEnrollID,
		Marks:           param.Marks,
		Grade:           score,
		CreatedBy:       uuid.NewV4(),
	}

	err := result.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "LecturerAddResult/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	results, err := result.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "LecturerAddResult/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return results, nil
}

func (s LecturerModule) LecturerUpdateResult(ctx context.Context, param LecturerUpdateResultParam) (interface{}, *helpers.Error) {

	score := util.GetGrade(param.Marks)

	result := models.ResultModel{
		StudentEnrollID: param.StudentEnrollID,
		Marks:           param.Marks,
		Grade:           score,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := result.UpdateByStudentEnroll(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, " LecturerUpdateResult/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	results, err := result.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "LecturerUpdateResult/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return results, nil
}
