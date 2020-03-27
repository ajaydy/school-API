package routers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"school/helpers"
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

	apiV1.Handle("/student/{id}", HandlerFunc(HandlerStudentDetail)).Methods(http.MethodGet)
	apiV1.Handle("/lecturer/{id}", HandlerFunc(HandlerLecturerDetail)).Methods(http.MethodGet)
	apiV1.Handle("/subjects/{id}", HandlerFunc(HandlerSubjectDetail)).Methods(http.MethodGet)
	apiV1.Handle("/attendance/{id}", HandlerFunc(HandlerAttendanceDetail)).Methods(http.MethodGet)
	apiV1.Handle("/classroom/{id}", HandlerFunc(HandlerClassroomDetail)).Methods(http.MethodGet)
	apiV1.Handle("/faculty/{id}", HandlerFunc(HandlerFacultyDetail)).Methods(http.MethodGet)
	apiV1.Handle("/intake/{id}", HandlerFunc(HandlerIntakeDetail)).Methods(http.MethodGet)
	apiV1.Handle("/program/{id}", HandlerFunc(HandlerProgramDetail)).Methods(http.MethodGet)
	apiV1.Handle("/result/{id}", HandlerFunc(HandlerResultDetail)).Methods(http.MethodGet)
	apiV1.Handle("/session/{id}", HandlerFunc(HandlerSessionDetail)).Methods(http.MethodGet)
	apiV1.Handle("/studentEnroll/{id}", HandlerFunc(HandlerStudentEnrollDetail)).Methods(http.MethodGet)

	//apiV1.Use(middleware.BasicAuthMiddleware)

	//apiV1.Handle("/students", HandlerFunc(HandlerStudents)).Methods(http.MethodGet)
	//apiV1.Handle("/students/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
	//	HandlerFunc(HandlerStudentsDetail), session.STUDENT_ROLE, session.ADMIN_ROLE))).Methods(http.MethodGet)
	//apiV1.Handle("/students", HandlerFunc(HandlerAddStudents)).Methods(http.MethodPost)
	//apiV1.Handle("/students/{id}", HandlerFunc(HandlerUpdateStudents)).Methods(http.MethodPut)
	//apiV1.Handle("/students/register", HandlerFunc(HandlerRegisterStudents)).Methods(http.MethodPost)
	//apiV1.Handle("/students/login", HandlerFunc(HandlerLoginStudents)).Methods(http.MethodPost)
	//
	//apiV1.Handle("/lecturers", HandlerFunc(HandlerLecturers)).Methods(http.MethodGet)
	//apiV1.Handle("/lecturers/{id}", HandlerFunc(HandlerLecturersDetail)).Methods(http.MethodGet)
	//apiV1.Handle("/lecturers", HandlerFunc(HandlerAddLecturers)).Methods(http.MethodPost)
	//apiV1.Handle("/lecturers/{id}", HandlerFunc(HandlerUpdateLecturers)).Methods(http.MethodPut)
	//
	//apiV1.Handle("/subjects", HandlerFunc(HandlerSubjects)).Methods(http.MethodGet)
	//apiV1.Handle("/subjects/{id}", HandlerFunc(HandlerSubjectsDetail)).Methods(http.MethodGet)
	//apiV1.Handle("/subjects", HandlerFunc(HandlerAddSubjects)).Methods(http.MethodPost)
	//apiV1.Handle("/subjects/{id}", HandlerFunc(HandlerUpdateSubjects)).Methods(http.MethodPut)
	//
	//apiV1.Handle("/score", HandlerFunc(HandlerStudentsWithScore)).Methods(http.MethodGet)
	//apiV1.Handle("/score/{id}", HandlerFunc(HandlerStudentsWithScoreDetail)).Methods(http.MethodGet)
	//
	//apiV1.Handle("/enroll", HandlerFunc(HandlerAddStudentEnroll)).Methods(http.MethodPost)

	return r
}
