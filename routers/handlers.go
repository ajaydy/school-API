package routers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"school/helpers"
	"school/middleware"
	"school/session"
)

type (
	HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, *helpers.Error)
)

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errs []string
	r.ParseForm()
	data, err := fn(w, r)
	if err != nil {
		errs = append(errs, err.Error())
		w.WriteHeader(err.StatusCode)
	}
	resp := helpers.Response{
		Data: data,
		BaseResponse: helpers.BaseResponse{
			Errors: errs,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		return
	}
}

func InitHandlers() *mux.Router {
	r := mux.NewRouter()

	http.Handle("/", r)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	//apiV1.Use(middleware.BasicAuthMiddleware)
	//apiV1.Handle("/student/timetable", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerTimetable), session.STUDENT_ROLE))).Methods(http.MethodGet)
	//apiV1.Handle("/student/enroll", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerAddStudentEnroll), session.STUDENT_ROLE))).Methods(http.MethodPost)
	//apiV1.Handle("/student/result", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerResult), session.STUDENT_ROLE))).Methods(http.MethodGet)
	//apiV1.Handle("/student/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerUpdateStudent), session.STUDENT_ROLE))).Methods(http.MethodPut)
	//apiV1.Handle("/student/{id}", HandlerFunc(HandlerStudentDetail)).Methods(http.MethodGet)
	//apiV1.Handle("/student", HandlerFunc(HandlerStudent)).Methods(http.MethodGet)
	//apiV1.Handle("/student/add", HandlerFunc(HandlerAddStudent)).Methods(http.MethodPost)
	//apiV1.Handle("/student/login", HandlerFunc(HandlerLoginStudent)).Methods(http.MethodPost)

	apiV1.Handle("/lecturer/attendance/{id}/is-attend", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerUpdateAttendanceIsAttend), session.LECTURER_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/lecturer/session/class/attendance", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerAttendanceByClass), session.LECTURER_ROLE))).Methods(http.MethodGet)
	apiV1.Handle("/lecturer/session/class/attendance", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerAddAttendance), session.LECTURER_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/lecturer/session/class/attendance", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerAttendanceByClass), session.LECTURER_ROLE))).Methods(http.MethodGet)
	apiV1.Handle("/lecturer/session/class", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerClassBySession), session.LECTURER_ROLE))).Methods(http.MethodGet)
	apiV1.Handle("/lecturer/session/student/result", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerLecturerUpdateResult), session.LECTURER_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/lecturer/session/student", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerStudentEnrollBySession), session.LECTURER_ROLE))).Methods(http.MethodGet)
	apiV1.Handle("/lecturer/session", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSessionByLecturer), session.LECTURER_ROLE))).Methods(http.MethodGet)

	//apiV1.Handle("/students", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerStudent), session.ADMIN_ROLE, session.STUDENT_ROLE, session.LECTURER_ROLE))).Methods(http.MethodGet)
	//apiV1.Handle("/admin/student/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerStudentDetail), session.ADMIN_ROLE))).Methods(http.MethodGet)
	//apiV1.Handle("/admin/student/{id}/update", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerStudentDetail), session.ADMIN_ROLE))).Methods(http.MethodGet)
	//apiV1.Handle("/admin/student/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerStudentDetail), session.ADMIN_ROLE))).Methods(http.MethodGet)

	apiV1.Handle("/student-enrolls", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerStudentEnrollAdd), session.STUDENT_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/student-enrolls/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerStudentEnrollDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/lecturers", middleware.SessionMiddleware(HandlerFunc(HandlerLecturerList))).Methods(http.MethodGet)
	apiV1.Handle("/lecturers/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerLecturerDetail))).Methods(http.MethodGet)
	apiV1.Handle("/lecturers", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerLecturerAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/lecturers/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerLecturerUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/lecturers/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerLecturerDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/students", middleware.SessionMiddleware(HandlerFunc(HandlerStudentList))).Methods(http.MethodGet)
	apiV1.Handle("/students/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerStudentDetail))).Methods(http.MethodGet)
	apiV1.Handle("/students", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerStudentAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/students/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerStudentUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/students/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerStudentDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/sessions", middleware.SessionMiddleware(HandlerFunc(HandlerSessionList))).Methods(http.MethodGet)
	apiV1.Handle("/sessions/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerSessionDetail))).Methods(http.MethodGet)
	apiV1.Handle("/sessions", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSessionAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/sessions/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSessionUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/sessions/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSessionDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/results", middleware.SessionMiddleware(HandlerFunc(HandlerResultList))).Methods(http.MethodGet)
	apiV1.Handle("/results/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerResultDetail))).Methods(http.MethodGet)
	apiV1.Handle("/results", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerResultAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/results/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerResultUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/results/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerResultDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/programs", middleware.SessionMiddleware(HandlerFunc(HandlerProgramList))).Methods(http.MethodGet)
	apiV1.Handle("/programs/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerProgramDetail))).Methods(http.MethodGet)
	apiV1.Handle("/programs", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerProgramAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/programs/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerProgramUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/programs/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerProgramDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/intakes", middleware.SessionMiddleware(HandlerFunc(HandlerIntakeList))).Methods(http.MethodGet)
	apiV1.Handle("/intakes/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerIntakeDetail))).Methods(http.MethodGet)
	apiV1.Handle("/intakes", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerIntakeAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/intakes/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerIntakeUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/intakes/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerIntakeDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/subjects", middleware.SessionMiddleware(HandlerFunc(HandlerSubjectList))).Methods(http.MethodGet)
	apiV1.Handle("/subjects/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerSubjectDetail))).Methods(http.MethodGet)
	apiV1.Handle("/subjects", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSubjectAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/subjects/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSubjectUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/subjects/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSubjectDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/classrooms/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerClassroomDetail))).Methods(http.MethodGet)
	apiV1.Handle("/classrooms", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerClassroomAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/classrooms/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerClassroomUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/classrooms/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerClassroomDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/faculties", middleware.SessionMiddleware(HandlerFunc(HandlerFacultyList))).Methods(http.MethodGet)
	apiV1.Handle("/faculties/{id}", middleware.SessionMiddleware(HandlerFunc(HandlerFacultyDetail))).Methods(http.MethodGet)
	apiV1.Handle("/faculties", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerFacultyAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/faculties/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerFacultyUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/faculties/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerFacultyDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/lecturer/login", HandlerFunc(HandlerLecturerLogin)).Methods(http.MethodPost)
	apiV1.Handle("/admin/login", HandlerFunc(HandlerAdminLogin)).Methods(http.MethodPost)
	apiV1.Handle("/student/login", HandlerFunc(HandlerStudentLogin)).Methods(http.MethodPost)

	return r
}
