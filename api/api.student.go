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
	"time"
)

type (
	StudentModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	StudentDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	StudentParamAdd struct {
		Name        string    `json:"name" valid:"length(3|50),required"`
		ProgramID   uuid.UUID `json:"program_id" valid:"required"`
		Address     string    `json:"address" valid:"optional"`
		DateOfBirth time.Time `json:"date_of_birth" valid:"required"`
		Gender      int       `json:"gender" valid:"range(0|1),required"`
		Email       string    `json:"email" valid:"email,required"`
		PhoneNo     string    `json:"phone_no" valid:"length(10|15),required"`
	}

	StudentParamUpdate struct {
		ID          uuid.UUID `json:"id" value:"required"`
		ProgramID   uuid.UUID `json:"program_id" value:"required"`
		Name        string    `json:"name" valid:"length(3|50),required"`
		Address     string    `json:"address" valid:"optional"`
		DateOfBirth time.Time `json:"date_of_birth" valid:"required"`
		Gender      int       `json:"gender" valid:"range(0|1),required"`
		Email       string    `json:"email" valid:"email,required"`
		PhoneNo     string    `json:"phone_no" valid:"length(10|15),required"`
		IsActive    bool      `json:"is_active" valid:"required"`
	}

	//StudentParamRegister struct {
	//	ID              uuid.UUID `json:"id"`
	//	Name            string    `json:"name" valid:"length(3|50),required"`
	//	Address         string    `json:"address" valid:"optional"`
	//	DateOfBirth     time.Time `json:"date_of_birth" valid:"required"`
	//	Gender          int       `json:"gender" valid:"range(0|1),required"`
	//	Email           string    `json:"email" valid:"email,required"`
	//	PhoneNo         string    `json:"phone_no" valid:"length(10|15),required"`
	//	Password        string    `json:"password" valid:"length(6|15),required"`
	//	ConfirmPassword string    `json:"confirm_password" valid:"length(6|15),required"`
	//}
	//
	StudentParamLogin struct {
		StudentCode string `json:"student_code" valid:"required"`
		Password    string `json:"password" valid:"length(1|50),required"`
	}

	StudentWithSession struct {
		Student models.StudentResponse `json:"student"`
		Session string                 `json:"session"`
	}
)

func NewStudentModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *StudentModule {
	return &StudentModule{
		db:     db,
		cache:  cache,
		name:   "module/student",
		logger: logger,
	}

}

//
//func (s StudentsModule) Register(ctx context.Context, param StudentParamRegister) (interface{}, *helpers.Error) {
//
//	if param.Password != param.ConfirmPassword {
//		return nil, helpers.ErrorWrap(errors.New("Invalid Password"), s.name, "Student/Register", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	password, err := bcrypt.GenerateFromPassword([]byte(param.Password), 12)
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Student/Register", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	student := models.StudentModel{
//		Name:        param.Name,
//		Address:     param.Address,
//		DateOfBirth: param.DateOfBirth,
//		Gender:      param.Gender,
//		Email:       param.Email,
//		PhoneNo:     param.PhoneNo,
//		Password:    string(password),
//		CreatedBy:   uuid.NewV4(),
//	}
//
//	err = student.Insert(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Student/Register", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//	session := session.Session{
//		UserID:     student.ID,
//		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
//		Expiry:     86400,
//		Role:       session.STUDENT_ROLE,
//	}
//
//	err = session.Store(ctx)
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Session/Register", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	studentSession := StudentWithSession{
//		Student: student.Response(),
//		Session: session.SessionKey,
//	}
//
//	return studentSession, nil
//}
//
func (s StudentModule) Login(ctx context.Context, param StudentParamLogin) (interface{}, *helpers.Error) {

	student, err := models.GetOneStudentByCode(ctx, s.db, param.StudentCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helpers.ErrorWrap(err, s.name, "Session/Login", helpers.IncorrectEmailMessage,
				http.StatusInternalServerError)
		}
		return nil, helpers.ErrorWrap(err, s.name, "List/GetOneStudentByCode", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(param.Password))
	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Invalid Password"), s.name, "Login/CompareHashAndPassword",
			helpers.IncorrectPasswordMessage,
			http.StatusInternalServerError)
	}

	session := session.Session{
		UserID:     student.ID,
		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
		Expiry:     86400,
		Role:       session.STUDENT_ROLE,
	}

	err = session.Store(ctx)

	students, err := student.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Login/StudentResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	studentSession := StudentWithSession{
		Student: students,
		Session: session.SessionKey,
	}

	return studentSession, nil

}

func (s StudentModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	student, err := models.GetAllStudent(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllStudent", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var studentsResponse []models.StudentResponse
	for _, students := range student {
		response, err := students.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/StudentResponse", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		studentsResponse = append(studentsResponse, response)
	}

	return studentsResponse, nil
}

func (s StudentModule) Detail(ctx context.Context, param StudentDetailParam) (interface{}, *helpers.Error) {
	student, err := models.GetOneStudent(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneStudent", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := student.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/StudentResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s StudentModule) Add(ctx context.Context, param StudentParamAdd) (interface{}, *helpers.Error) {

	password := util.RandomString(12)

	program, err := models.GetOneProgram(ctx, s.db, param.ProgramID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneProgram", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	programCode := program.Code

	facultyID := program.FacultyID

	faculty, err := models.GetOneFaculty(ctx, s.db, facultyID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneFaculty", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	facultyCode := faculty.Code

	yearCode, err := util.GetYearCode()
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetYearCode", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	studentCode := fmt.Sprintf("%d%d%d%d%d%s", 1, yearCode, programCode, 0, facultyCode, "000")

	student := models.StudentModel{
		Name:        param.Name,
		ProgramID:   param.ProgramID,
		StudentCode: studentCode,
		Address:     param.Address,
		DateOfBirth: param.DateOfBirth,
		Gender:      param.Gender,
		Password:    password,
		Email:       param.Email,
		PhoneNo:     param.PhoneNo,
		CreatedBy:   uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err = student.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := student.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s StudentModule) Update(ctx context.Context, param StudentParamUpdate) (interface{}, *helpers.Error) {

	student := models.StudentModel{
		ID:          param.ID,
		ProgramID:   param.ProgramID,
		Name:        param.Name,
		Address:     param.Address,
		DateOfBirth: param.DateOfBirth,
		Gender:      param.Gender,
		IsActive:    param.IsActive,
		Email:       param.Email,
		PhoneNo:     param.PhoneNo,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}
	err := student.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	students, err := student.Response(ctx, s.db, s.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/StudentResponse", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return students, nil

}
