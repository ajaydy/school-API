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
	IntakeModule struct {
		db     *sql.DB
		cache  *redis.Pool
		name   string
		logger *helpers.Logger
	}

	IntakeDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

//	IntakeParamAdd struct {
//		Year  string `json:"year"valid:"required"`
//		Month int    `json:"month" valid:"range(1|12),required"`
//	}
//
//	IntakeParamUpdate struct {
//		ID    uuid.UUID `json:"id"`
//		Year  string    `json:"year"valid:"required"`
//		Month int       `json:"month" valid:"range(1|12),required"`
//	}
)

func NewIntakeModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *IntakeModule {
	return &IntakeModule{
		db:     db,
		cache:  cache,
		name:   "module/intake",
		logger: logger,
	}
}

//func (s IntakeModule) List(ctx context.Context) (interface{}, *helpers.Error) {
//	intakes, err := models.GetAllIntake(ctx, s.db)
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllIntake", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	var intakeResponse []models.IntakeResponse
//	for _, intake := range intakes {
//		intakeResponse = append(intakeResponse, intake.Response())
//	}
//
//	return intakeResponse, nil
//}

func (s IntakeModule) Detail(ctx context.Context, param IntakeDetailParam) (interface{}, *helpers.Error) {
	intake, err := models.GetOneIntake(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneIntake", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return intake.Response(), nil
}

//func (s IntakeModule) Add(ctx context.Context, param IntakeParamAdd) (interface{}, *helpers.Error) {
//	intake := models.IntakeModel{
//		Year:      param.Year,
//		Month:     param.Month,
//		CreatedBy: uuid.NewV4(),
//	}
//
//	err := intake.Insert(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return intake, nil
//}
//
//func (s IntakeModule) Update(ctx context.Context, param IntakeParamUpdate) (interface{}, *helpers.Error) {
//
//	intake := models.IntakeModel{
//		ID:    param.ID,
//		Year:  param.Year,
//		Month: param.Month,
//		UpdatedBy: uuid.NullUUID{
//			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
//			Valid: true,
//		},
//	}
//	err := intake.Update(ctx, s.db)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
//			http.StatusInternalServerError)
//	}
//
//	return intake, nil
//
//}
