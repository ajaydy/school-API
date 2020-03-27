package middleware

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"school/models"
)

func basicAuth(w http.ResponseWriter, r *http.Request) (models.AdminModel, bool) {
	username, password, ok := r.BasicAuth()

	if !ok {
		return models.AdminModel{}, false
	}

	user, err := models.GetOneUserByUsername(r.Context(), dbPool, username)

	if err != nil {
		logger.Err.Printf(`%v`, err)
		return models.AdminModel{}, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		logger.Err.Printf(`%v`, err)
		return models.AdminModel{}, false
	}

	return user, true

}
