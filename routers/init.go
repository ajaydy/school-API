package routers

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"school/api"
	"school/helpers"
)

var (
	dbPool         *sql.DB
	cachePool      *redis.Pool
	logger         *helpers.Logger
	studentService *api.StudentsModule
)

func Init(db *sql.DB, cache *redis.Pool, log *helpers.Logger) {
	dbPool = db
	cachePool = cache
	logger = log
	studentService = api.NewStudentsModule(dbPool, cachePool)
}
