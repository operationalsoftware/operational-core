package handlers

import (
	"log"
	"net/http"
	"operationalcore/components"
	"operationalcore/db"
	"operationalcore/model"

	bcrypt "golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Hash password
	confirmPassword := r.FormValue("confirmPassword")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(confirmPassword), bcrypt.DefaultCost)

	if err != nil {
		// TODO: Handle error in UI
		_ = components.InputHelper(&components.InputHelperProps{
			Label: "Submission Error",
			Type:  components.InputHelperTypeError,
		}).Render(w)
		log.Fatal(err)
	}

	hashedPasswordString := string(hashedPassword)
	// Create user in db
	dbInstance := db.UseDB()
	if dbInstance != nil {
		defer dbInstance.Close()
	}

	if err != nil {
		log.Fatal(err)
	}

	queryErr := model.AddUser(dbInstance, model.User{
		Username:       r.FormValue("username"),
		Email:          r.FormValue("email"),
		FirstName:      r.FormValue("first_name"),
		LastName:       r.FormValue("last_name"),
		HashedPassword: hashedPasswordString,
	})

	if queryErr != nil {
		log.Fatal(queryErr)
	}

	// Redirect to user view
	w.Header().Set("hx-redirect", "/users")

}
