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
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_ = indexView(&indexViewProps{
		Ctx: ctx,
	}).Render(w)
}

func userViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	_ = userView(&userViewProps{
		Id:  id,
		Ctx: ctx,
	}).Render(w)
}

func addUserViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_ = addUserView(&addUserViewProps{
		Ctx: ctx,
	}).Render(w)
}

func addUserAPIViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_ = addUserAPIView(&addUserAPIViewProps{
		Ctx: ctx,
	}).Render(w)
}

func validateAddUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var newUser userModel.NewUser
	_ = utils.DecodeForm(r.Form, &newUser)

	_, validationErrors := userModel.ValidateNewUser(newUser)

	_ = addUserForm(&addUserFormProps{
		values:           r.Form,
		validationErrors: validationErrors,
	}).Render(w)

}

func validateAddAPIUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var newAPIUser userModel.NewAPIUser
	_ = utils.DecodeForm(r.Form, &newAPIUser)

	_, validationErrors := userModel.ValidateNewAPIUser(newAPIUser)

	_ = addApiUserForm(&addApiUserFormProps{
		values:           r.Form,
		validationErrors: validationErrors,
	}).Render(w)

}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		fmt.Println("Error:", "Forbidden")
		return
	}

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
		http.Error(w, "Error adding user", http.StatusInternalServerError)
	}

	// Redirect to users view
	w.Header().Set("hx-redirect", "/users")
}

func addUserAPIHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		fmt.Println("Error:", "Forbidden")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var newAPIUser userModel.NewAPIUser
	_ = utils.DecodeForm(r.Form, &newAPIUser)

	password, err := GenerateRandomPassword(24)
	if err != nil {
		log.Panic(err)
	}
	newAPIUser.Password = password
	newAPIUser.Roles = []string{"User Admin"}

	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := userModel.ValidateNewAPIUser(newAPIUser)
	if !valid {
		_ = addApiUserForm(&addApiUserFormProps{
			values:           r.Form,
			validationErrors: validationErrors,
			isSubmission:     true,
		}).Render(w)
		return
	}

	db := db.UseDB()
	err = userModel.AddAPIUser(db, newAPIUser)

	if err != nil {
		http.Error(w, "Error adding API user", http.StatusInternalServerError)
	}

	w.Header().Set("HX-Reswap", "outerHTML")
	w.Header().Set("HX-Reselect", ".card")

	_ = APIUserCredentials(&APIUserCredentialsProps{
		Username: newAPIUser.Username,
		Password: password,
	}).Render(w)
}

func editUserViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

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

	_ = editUserView(&editUserViewProps{
		User: user,
		Ctx:  ctx,
	}).Render(w)
}

func validateEditUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	db := db.UseDB()
	user, err := userModel.ByID(db, id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusBadRequest)
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

}

func editUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
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
			values:           r.Form,
			validationErrors: validationErrors,
			isSubmission:     true,
		}).Render(w)
		return
	}

	db := db.UseDB()
	err = userModel.Update(db, id, userUpdate)

	if err != nil {
		log.Panic(err)
	}

	// Redirect to user view
	w.Header().Set("hx-redirect", "/users")

}

func resetPasswordViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

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

	_ = resetPasswordView(&resetPasswordViewProps{
		User: user,
		Ctx:  ctx,
	}).Render(w)
}

func validateResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var passwordReset userModel.PasswordReset

	err = utils.DecodeForm(r.Form, &passwordReset)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	_, validationErrors := userModel.ValidatePasswordReset(passwordReset)

	_ = resetPasswordForm(&resetPasswordFormProps{
		userID:           id,
		values:           r.Form,
		validationErrors: validationErrors,
	}).Render(w)

}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var passwordReset userModel.PasswordReset
	err = utils.DecodeForm(r.Form, &passwordReset)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := userModel.ValidatePasswordReset(passwordReset)
	if !valid {
		_ = resetPasswordForm(&resetPasswordFormProps{
			userID:           id,
			values:           r.Form,
			validationErrors: validationErrors,
			isSubmission:     true,
		}).Render(w)
		return
	}

	db := db.UseDB()
	err = userModel.ResetPassword(db, id, passwordReset)

	if err != nil {
		http.Error(w, "Error resetting password", http.StatusInternalServerError)
		return
	}

	// Redirect to user view
	w.Header().Set("hx-redirect", fmt.Sprintf("/users/%d", id))
}
