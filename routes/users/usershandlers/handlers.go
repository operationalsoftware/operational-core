package usershandlers

import (
	"app/internal/appsort"
	"app/internal/db"
	"app/internal/reqcontext"
	"app/internal/urlvalues"
	"app/models/usermodel"
	"app/routes/users/usersviews"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func UsersHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	db := db.UseDB()
	sort := appsort.Sort{}
	sort.ParseQueryParam(ctx.Req.URL.Query().Get("Sort"), usermodel.ListSortableKeys)

	users, err := usermodel.List(db, usermodel.ListQuery{
		Sort: sort,
	})
	if err != nil {
		http.Error(w, "Error listing users", http.StatusInternalServerError)
		return
	}
	count, err := usermodel.Count(db)
	if err != nil {
		http.Error(w, "Error counting users", http.StatusInternalServerError)
		return
	}

	_ = usersviews.UsersHomePage(&usersviews.UsersHomePageProps{
		Ctx:       ctx,
		Users:     users,
		UserCount: count,
		Sort:      sort,
		Page:      1,
		PageSize:  10,
	}).Render(w)
}

func UserPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	_ = usersviews.UserPage(&usersviews.UserPageProps{
		Id:  id,
		Ctx: ctx,
	}).Render(w)
}

func AddUserPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_ = usersviews.AddUserPage(&usersviews.AddUserPageProps{
		Ctx: ctx,
	}).Render(w)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var newUser usermodel.NewUser
	err = urlvalues.Unmarshal(r.Form, &newUser)

	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := usermodel.ValidateNewUser(newUser)
	if !valid {
		_ = usersviews.AddUserPage(&usersviews.AddUserPageProps{
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	db := db.UseDB()
	err = usermodel.Add(db, newUser)

	if err != nil {
		http.Error(w, "Error adding user", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func AddAPIUserPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = usersviews.AddAPIUserPage(&usersviews.AddAPIUserPageProps{
		Ctx: ctx,
	}).Render(w)
}

func AddAPIUser(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var newAPIUser usermodel.NewAPIUser
	_ = urlvalues.Unmarshal(r.Form, &newAPIUser)

	password, err := usermodel.GenerateRandomPassword(24)
	if err != nil {
		log.Panic(err)
	}
	newAPIUser.Password = password

	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := usermodel.ValidateNewAPIUser(newAPIUser)
	if !valid {
		_ = usersviews.AddAPIUserPage(&usersviews.AddAPIUserPageProps{
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	db := db.UseDB()
	err = usermodel.AddAPIUser(db, newAPIUser)

	if err != nil {
		http.Error(w, "Error adding API user", http.StatusInternalServerError)
	}

	_ = usersviews.APIUserCredentialsPage(&usersviews.APIUserCredentialsPageProps{
		Username: newAPIUser.Username,
		Password: password,
	}).Render(w)
}

func EditUserPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	db := db.UseDB()
	user, err := usermodel.ByID(db, id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusBadRequest)
		return
	}

	_ = usersviews.EditUserPage(&usersviews.EditUserPageProps{
		User: user,
		Ctx:  ctx,
	}).Render(w)
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var userUpdate usermodel.UserUpdate
	err = urlvalues.Unmarshal(r.Form, &userUpdate)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := usermodel.ValidateUserUpdate(userUpdate)
	if !valid {
		_ = usersviews.EditUserPage(&usersviews.EditUserPageProps{
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	db := db.UseDB()
	err = usermodel.Update(db, id, userUpdate)

	if err != nil {
		log.Panic(err)
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	db := db.UseDB()
	user, err := usermodel.ByID(db, id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusBadRequest)
		return
	}

	_ = usersviews.ResetPasswordPage(&usersviews.ResetPasswordPageProps{
		User: user,
		Ctx:  ctx,
	}).Render(w)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	db := db.UseDB()

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var passwordReset usermodel.PasswordReset
	err = urlvalues.Unmarshal(r.Form, &passwordReset)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := usermodel.ValidatePasswordReset(passwordReset)
	if !valid {
		user, err := usermodel.ByID(db, id)
		if err != nil {
			http.Error(w, "Error getting user", http.StatusBadRequest)
			return
		}

		_ = usersviews.ResetPasswordPage(&usersviews.ResetPasswordPageProps{
			Ctx:              ctx,
			User:             user,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	err = usermodel.ResetPassword(db, id, passwordReset)

	if err != nil {
		http.Error(w, "Error resetting password", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d", id), http.StatusSeeOther)
}
