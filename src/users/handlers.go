package users

import (
	"app/db"
	userModel "app/src/users/model"
	"app/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func indexViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = indexView(&indexViewProps{
		Ctx: ctx,
	}).Render(w)
}

func userViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ctx := utils.GetContext(r)
	_ = userView(&userViewProps{
		Id:  id,
		Ctx: ctx,
	}).Render(w)
}

func addUserViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	_ = addUserView(&addUserViewProps{
		Ctx: ctx,
	}).Render(w)
}

func validateAddUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser userModel.NewUser

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	err = utils.DecodeForm(r.Form, &newUser)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	_, validationErrors := userModel.ValidateNewUser(newUser)

	_ = addUserForm(&addUserFormProps{
		values:           r.Form,
		validationErrors: validationErrors,
	}).Render(w)

	return
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var newUser userModel.NewUser
	err = utils.DecodeForm(r.Form, &newUser)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := userModel.ValidateNewUser(newUser)
	if !valid {
		_ = addUserForm(&addUserFormProps{
			values:           r.Form,
			validationErrors: validationErrors,
			isSubmission:     true,
		}).Render(w)
		return
	}

	db := db.UseDB()
	err = userModel.Add(db, newUser)

	if err != nil {
		log.Panic(err)
	}

	// Redirect to users view
	w.Header().Set("hx-redirect", "/users")
}

func editUserViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	db := db.UseDB()
	user, err := userModel.ByID(db, id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusBadRequest)
		return
	}

	ctx := utils.GetContext(r)
	_ = editUserView(&editUserViewProps{
		User: user,
		Ctx:  ctx,
	}).Render(w)
}

func validateEditUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	db := db.UseDB()
	user, err := userModel.ByID(db, id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var userUpdate userModel.UserUpdate

	err = utils.DecodeForm(r.Form, &userUpdate)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	_, validationErrors := userModel.ValidateUserUpdate(userUpdate)

	_ = editUserForm(&editUserFormProps{
		user:             user,
		values:           r.Form,
		validationErrors: validationErrors,
	}).Render(w)

	return
}

func editUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	db := db.UseDB()
	user, err := userModel.ByID(db, id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var userUpdate userModel.UserUpdate
	err = utils.DecodeForm(r.Form, &userUpdate)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := userModel.ValidateUserUpdate(userUpdate)
	if !valid {
		_ = editUserForm(&editUserFormProps{
			user:             user,
			values:           r.Form,
			validationErrors: validationErrors,
			isSubmission:     true,
		}).Render(w)
		return
	}

	err = userModel.Update(db, id, userUpdate)

	if err != nil {
		log.Panic(err)
	}

	// Redirect to user view
	w.Header().Set("hx-redirect", "/users")

	return
}
