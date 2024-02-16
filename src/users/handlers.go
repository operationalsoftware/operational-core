package users

import (
	"app/components"
	"app/db"
	userModel "app/src/users/model"
	"app/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
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
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_ = userView(&userViewProps{
		Id:  id,
		Ctx: ctx,
	}).Render(w)
}

func resetPasswordUserViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_ = resetPasswordUserView(&resetPasswordUserViewProps{
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

func editUserViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_ = editUserView(&editUserViewProps{
		Id:  id,
		Ctx: ctx,
	}).Render(w)
}

func validateConfirmPasswordHandler(w http.ResponseWriter, r *http.Request) {
	confirmPassword := r.FormValue("ConfirmPassword")
	password := r.FormValue("Password")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(confirmPassword, "required,min=3,max=20")

	var helperText string

	if err != nil {
		helperText = "Passwords must be between 3 and 20 characters"
	}

	if confirmPassword != password {
		helperText = "Passwords do not match"
	}

	_ = confirmPasswordInput(&confirmPasswordInputProps{
		ValidationError: helperText,
		Value:           confirmPassword,
	}).Render(w)
}

func validateEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("Email")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(email, "email")

	var helperText string

	if err != nil {
		helperText = "Email must be a valid email address"
	}

	_ = emailInput(&emailInputProps{
		ValidationError: helperText,
		Value:           email,
	}).Render(w)
}

func validateFirstNameHandler(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("FirstName")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(firstName, "required,gte=3,lte=20")

	var helperText string

	if err != nil {
		helperText = "First name must be between 3 and 20 characters"
	}

	_ = firstNameInput(&firstNameInputProps{
		ValidationError: helperText,
		Value:           firstName,
	}).Render(w)
}

func validateLastNameHandler(w http.ResponseWriter, r *http.Request) {
	lastName := r.FormValue("LastName")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(lastName, `required,min=3,max=20`)

	var helperText string

	if err != nil {
		helperText = "Last name must be between 3 and 20 characters"
	}

	_ = lastNameInput(&lastNameInputProps{
		ValidationError: helperText,
		Value:           lastName,
	}).Render(w)
}

func validatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("Password")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(password, "required,gte=3,lte=20")

	var helperText string

	if err != nil {
		helperText = "Password must be between 3 and 20 characters"
	}

	_ = passwordInput(&passwordInputProps{
		ValidationError: helperText,
		Value:           password,
	}).Render(w)
}

func validateUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("Username")

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Var(username, "required,gte=3,lte=20")

	var helperText string

	if err != nil {
		helperText = "Username must be between 3 and 20 characters"
	}

	_ = usernameInput(&usernameInputProps{
		ValidationError: helperText,
		Value:           username,
	}).Render(w)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		fmt.Println("Error:", "Forbidden")
		return
	}
	// Create user in db
	db := db.UseDB()

	err := userModel.Add(db, userModel.NewUser{
		Username:  r.FormValue("Username"),
		Email:     utils.StringToNullString(r.FormValue("Email")),
		FirstName: utils.StringToNullString(r.FormValue("FirstName")),
		LastName:  utils.StringToNullString(r.FormValue("LastName")),
		Password:  r.FormValue("Password"),
	})

	if err != nil {
		log.Panic(err)
	}

	// Redirect to users view
	w.Header().Set("hx-redirect", "/users")
}

func editUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")

	if !isAdmin {
		fmt.Println("Error:", "Forbidden")
		return
	}

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	type User struct {
		FirstName string `validate:"required"`
		LastName  string `validate:"required"`
		Email     string `validate:"email"`
		Username  string `validate:"required,gte=3,lte=20"`
	}

	var user User = User{
		FirstName: r.FormValue("FirstName"),
		LastName:  r.FormValue("LastName"),
		Email:     r.FormValue("Email"),
		Username:  r.FormValue("Username"),
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(user)

	if err != nil {
		_ = components.InputHelper(&components.InputHelperProps{
			Label: "Submission Error",
			Type:  components.InputHelperTypeError,
		}).Render(w)
		return
	}

	db := db.UseDB()

	if err != nil {
		log.Fatal(err)
	}

	roles := r.FormValue("roles")
	roles = strings.Trim(roles, "[]")
	userRoles := strings.Split(roles, ",")

	query := userModel.Update(db, id, userModel.UserUpdate{
		FirstName: utils.StringToNullString(user.FirstName),
		LastName:  utils.StringToNullString(user.LastName),
		Email:     utils.StringToNullString(user.Email),
		Username:  user.Username,
		Roles:     userRoles,
	})

	if query != nil {
		_ = components.InputHelper(&components.InputHelperProps{
			Label: "Submission Error",
			Type:  components.InputHelperTypeError,
		}).Render(w)
		return
	}

	// Redirect to user view
	w.Header().Set("hx-redirect", "/users")

}

func resetPasswordUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ctx := utils.GetContext(r)
	isAdmin := utils.CheckRole(ctx.User.Roles, "User Admin")
	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	db := db.UseDB()

	err = userModel.ResetPassword(db, id, r.FormValue("Password"))

	if err != nil {
		log.Fatal(err)
	}

	// Redirect to user view
	w.Header().Set("hx-redirect", fmt.Sprintf("/users/%d", id))
}
