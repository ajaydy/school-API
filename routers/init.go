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
	studentService       *api.StudentModule
	lecturerService      *api.LecturerModule
	subjectService       *api.SubjectModule
	intakeService        *api.IntakeModule
	studentEnrollService *api.StudentEnrollModule
	attendanceService    *api.AttendanceModule
	classroomService     *api.ClassroomModule
	facultyService       *api.FacultyModule
	programService       *api.ProgramModule
	resultService        *api.ResultModule
	sessionService       *api.SessionModule
	adminService         *api.AdminModule
	classService         *api.ClassModule
)

func Init(db *sql.DB, cache *redis.Pool, log *helpers.Logger) {
	dbPool = db
	cachePool = cache
	logger = log
	studentService = api.NewStudentModule(dbPool, cachePool, logger)
	lecturerService = api.NewLecturerModule(dbPool, cachePool, logger)
	subjectService = api.NewSubjectModule(dbPool, cachePool, logger)
	intakeService = api.NewIntakeModule(dbPool, cachePool, logger)
	studentEnrollService = api.NewStudentEnrollModule(dbPool, cachePool, logger)
	attendanceService = api.NewAttendanceModule(dbPool, cachePool, logger)
	classroomService = api.NewClassroomModule(dbPool, cachePool, logger)
	facultyService = api.NewFacultyModule(dbPool, cachePool, logger)
	programService = api.NewProgramModule(dbPool, cachePool, logger)
	resultService = api.NewResultModule(dbPool, cachePool, logger)
	sessionService = api.NewSessionModule(dbPool, cachePool, logger)
	adminService = api.NewAdminModule(dbPool, cachePool, logger)
	classService = api.NewClassModule(dbPool, cachePool, logger)
}
