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
)

type (
	AdminModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	AdminDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	AdminLoginParam struct {
		Username string `json:"username"valid:"required"`
		Password string `json:"password"valid:"required"`
	}

	AdminWithSession struct {
		Admin   models.AdminResponse `json:"admin"`
		Session string               `json:"session"`
	}

	AdminPasswordUpdateParam struct {
		ID                 uuid.UUID `json:"id" valid:"required"`
		CurrentPassword    string    `json:"current_password" valid:"required"`
		NewPassword        string    `json:"new_password" valid:"length(5|50),required"`
		ConfirmNewPassword string    `json:"confirm_new_password" valid:"required"`
	}
)

func NewAdminModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *AdminModule {
	return &AdminModule{
		db:     db,
		cache:  cache,
		name:   "module/admin",
		logger: logger,
	}
}

func (s AdminModule) Login(ctx context.Context, param AdminLoginParam) (interface{}, *helpers.Error) {

	admin, err := models.GetOneAdminByUsername(ctx, s.db, param.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helpers.ErrorWrap(err, s.name, "Session/Login", helpers.IncorrectEmailMessage,
				http.StatusInternalServerError)
		}
		return nil, helpers.ErrorWrap(err, s.name, "Login/GetOneAdminByUsername", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(param.Password))
	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Invalid Password"), s.name, "Login/CompareHashAndPassword",
			helpers.IncorrectPasswordMessage,
			http.StatusInternalServerError)
	}

	session := session.Session{
		UserID:     admin.ID,
		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
		Expiry:     86400,
		Role:       session.ADMIN_ROLE,
	}

	err = session.Store(ctx)

	adminSession := AdminWithSession{
		Admin:   admin.Response(),
		Session: session.SessionKey,
	}

	return adminSession, nil

}

func (s AdminModule) PasswordUpdate(ctx context.Context, param AdminPasswordUpdateParam) (interface{}, *helpers.Error) {

	admin, err := models.GetOneAdmin(ctx, s.db, param.ID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "PasswordUpdate/GetOneAdmin", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(param.CurrentPassword))

	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Current Password Is Incorrect!"), s.name,
			"PasswordUpdate/CompareHashAndPassword",
			helpers.IncorrectPasswordMessage,
			http.StatusInternalServerError)
	}

	if param.NewPassword == param.CurrentPassword {
		return nil, helpers.ErrorWrap(errors.New("New Password Cannot Be Same With Current Password"), s.name,
			"PasswordUpdate/CurrentPasswordComparison", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	if param.NewPassword != param.ConfirmNewPassword {
		return nil, helpers.ErrorWrap(errors.New("New Password Does Not Match"), s.name,
			"PasswordUpdate/NewPassword", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(param.NewPassword), 12)

	admin = models.AdminModel{
		ID:       param.ID,
		Password: string(password),
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}
	err = admin.PasswordUpdate(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "PasswordUpdate/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	adminResponse := models.AdminUpdatePasswordResponse{
		Message: "Password Successfully Changed",
	}

	return adminResponse, nil

}
