package routers

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"school/api"
	"school/helpers"
)

var (
	dbPool               *sql.DB
	cachePool            *redis.Pool
	logger               *helpers.Logger
	studentService       *api.StudentsModule
	lecturerService      *api.LecturersModule
	subjectService       *api.SubjectModule
	intakeService        *api.IntakeModule
	studentenrollService *api.StudentEnrollModule
)

func Init(db *sql.DB, cache *redis.Pool, log *helpers.Logger) {
	dbPool = db
	cachePool = cache
	logger = log
	studentService = api.NewStudentsModule(dbPool, cachePool)
	lecturerService = api.NewLecturersModule(dbPool, cachePool)
	subjectService = api.NewSubjectsModule(dbPool, cachePool)
	intakeService = api.NewIntakesModule(dbPool, cachePool)
	studentenrollService = api.NewStudentEnrollModule(dbPool, cachePool)

}
