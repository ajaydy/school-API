package session

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"school/helpers"
)

var (
	dbPool    *sql.DB
	cachePool *redis.Pool
	logger    *helpers.Logger
)

func Init(db *sql.DB, cache *redis.Pool, log *helpers.Logger) {
	dbPool = db
	cachePool = cache
	logger = log
}
