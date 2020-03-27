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
	StudentModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	StudentDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	//StudentParamAdd struct {
	//	Name        string    `json:"name" valid:"length(3|50),required"`
	//	Address     string    `json:"address" valid:"optional"`
	//	DateOfBirth time.Time `json:"date_of_birth" valid:"required"`
	//	Gender      int       `json:"gender" valid:"range(0|1),required"`
	//	Email       string    `json:"email" valid:"email,required"`
	//	PhoneNo     string    `json:"phone_no" valid:"length(10|15),required"`
	//}
	//
	//StudentParamUpdate struct {
	//	ID          uuid.UUID `json:"id"`
	//	Name        string    `json:"name" valid:"length(3|50),required"`
	//	Address     string    `json:"address" valid:"optional"`
	//	DateOfBirth time.Time `json:"date_of_birth" valid:"required"`
	//	Gender      int       `json:"gender" valid:"range(0|1),required"`
	//	Email       string    `json:"email" valid:"email,required"`
	//	PhoneNo     string    `json:"phone_no" valid:"length(10|15),required"`
	//}
	//
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
	//StudentParamLogin struct {
	//	Email    string `json:"email" valid:"email,required"`
	//	Password string `json:"password" valid:"length(6|15),required"`
	//}
	//
	//StudentWithSession struct {
	//	Student models.StudentResponse `json:"student"`
	//	Session string                 `json:"session"`
	//}
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
//func (s StudentsModule) Login(ctx context.Context, param StudentParamLogin) (interface{}, *helpers.Error) {
//
//	students, err := models.GetOneStudentByEmail(ctx, s.db, param.Email)
//
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, helpers.ErrorWrap(err, s.name, "Session/Login", helpers.IncorrectEmailMessage,
//				http.StatusInternalServerError)
//		}
//		return nil, helpers.ErrorWrap(err, s.name, "List/GetOneStudentByEmail", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(students.Password), []byte(param.Password))
//	if err != nil {
//		return nil, helpers.ErrorWrap(errors.New("Invalid Password"), s.name, "Login/CompareHashAndPassword",
//			helpers.IncorrectPasswordMessage,
//			http.StatusInternalServerError)
//	}
//
//	session := session.Session{
//		UserID:     students.ID,
//		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
//		Expiry:     86400,
//		Role:       session.STUDENT_ROLE,
//	}
//
//	err = session.Store(ctx)
//
//	studentSession := StudentWithSession{
//		Student: students.Response(),
//		Session: session.SessionKey,
//	}
//
//	return studentSession, nil
//
//}
//func (s StudentsModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
//	students, err := models.GetAllStudent(ctx, s.db, filter)
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllStudent", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	var studentsResponse []models.StudentResponse
//	for _, student := range students {
//		response, err := student.Response(ctx, s.db, s.logger)
//		if err != nil {
//			return nil, helpers.ErrorWrap(err, s.name, "List/StudentResponse", helpers.InternalServerError,
//				http.StatusInternalServerError)
//		}
//		studentsResponse = append(studentsResponse, response)
//	}
//
//	return studentsResponse, nil
//}

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

//
//func (s StudentsModule) Add(ctx context.Context, param StudentParamAdd) (interface{}, *helpers.Error) {
//	student := models.StudentModel{
//		Name:        param.Name,
//		Address:     param.Address,
//		DateOfBirth: param.DateOfBirth,
//		Gender:      param.Gender,
//		Email:       param.Email,
//		PhoneNo:     param.PhoneNo,
//		CreatedBy:   uuid.NewV4(),
//	}
//
//	err := student.Insert(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return student, nil
//}
//
//func (s StudentsModule) Update(ctx context.Context, param StudentParamUpdate) (interface{}, *helpers.Error) {
//
//	student := models.StudentModel{
//		ID:          param.ID,
//		Name:        param.Name,
//		Address:     param.Address,
//		DateOfBirth: param.DateOfBirth,
//		Gender:      param.Gender,
//		Email:       param.Email,
//		PhoneNo:     param.PhoneNo,
//		UpdatedBy: uuid.NullUUID{
//			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
//			Valid: true,
//		},
//	}
//	err := student.Update(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return student, nil
//
//}
