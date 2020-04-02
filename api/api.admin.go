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
