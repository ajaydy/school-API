package api

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
	"school/models"
	"school/util"
	"time"
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

	IntakeAddParam struct {
		Year      string    `json:"year"valid:"required"`
		Month     int       `json:"month" valid:"required"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	IntakeUpdateParam struct {
		ID        uuid.UUID `json:"id"`
		Year      string    `json:"year"valid:"required"`
		Month     int       `json:"month" valid:"required"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	IntakeDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewIntakeModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *IntakeModule {
	return &IntakeModule{
		db:     db,
		cache:  cache,
		name:   "module/intake",
		logger: logger,
	}
}

func (s IntakeModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {
	intakes, err := models.GetAllIntake(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllIntake", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var intakesResponse []models.IntakeResponse
	for _, intake := range intakes {
		intakesResponse = append(intakesResponse, intake.Response())
	}

	return intakesResponse, nil
}

func (s IntakeModule) Detail(ctx context.Context, param IntakeDetailParam) (interface{}, *helpers.Error) {
	intake, err := models.GetOneIntake(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneIntake", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return intake.Response(), nil
}

func (s IntakeModule) Add(ctx context.Context, param IntakeAddParam) (interface{}, *helpers.Error) {

	if param.Month != 4 && param.Month != 7 && param.Month != 11 {
		return nil, helpers.ErrorWrap(errors.New("Invalid Month"), s.name, "Add/Month",
			helpers.IncorrectMonthMessage,
			http.StatusBadRequest)
	}

	trimester := util.GetTrimester(param.Month)

	intake := models.IntakeModel{
		Trimester: trimester,
		Year:      param.Year,
		Month:     param.Month,
		StartDate: param.StartDate,
		EndDate:   param.EndDate,
		CreatedBy: uuid.NewV4(),
	}

	err := intake.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return intake.Response(), nil
}

func (s IntakeModule) Update(ctx context.Context, param IntakeUpdateParam) (interface{}, *helpers.Error) {

	if param.Month != 4 && param.Month != 7 && param.Month != 11 {
		return nil, helpers.ErrorWrap(errors.New("Invalid Month"), s.name, "Update/Month",
			helpers.IncorrectMonthMessage,
			http.StatusBadRequest)
	}
	trimester := util.GetTrimester(param.Month)

	intake := models.IntakeModel{
		ID:        param.ID,
		Trimester: trimester,
		Year:      param.Year,
		Month:     param.Month,
		StartDate: param.StartDate,
		EndDate:   param.EndDate,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}
	err := intake.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return intake.Response(), nil

}

func (s IntakeModule) Delete(ctx context.Context, param IntakeDeleteParam) (interface{}, *helpers.Error) {

	intake := models.IntakeModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := intake.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
