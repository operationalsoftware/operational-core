package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/userview"
	"app/pkg/appsort"
	"app/pkg/appurl"
	"app/pkg/encryptcredentials"
	"app/pkg/reqcontext"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UsersHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var uv usersHomePageURLVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	sort := appsort.Sort{}
	err = sort.ParseQueryParam(model.User{}, uv.Sort)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing sort: %v", err), http.StatusBadRequest)
		return
	}

	users, count, err := h.userService.GetUsers(r.Context(), model.GetUsersQuery{
		Sort:     sort,
		Page:     uv.Page,
		PageSize: uv.PageSize,
	})
	if err != nil {
		log.Panicln(err)
		http.Error(w, "Error listing users", http.StatusInternalServerError)
		return
	}

	_ = userview.UsersHomePage(&userview.UsersHomePageProps{
		Ctx:       ctx,
		Users:     users,
		UserCount: count,
		Sort:      sort,
		Page:      uv.Page,
		PageSize:  uv.PageSize,
	}).Render(w)
}

type usersHomePageURLVals struct {
	Sort     string
	Page     int
	PageSize int
}

func (uv *usersHomePageURLVals) normalise() {
	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}
}

func (h *UserHandler) UserPage(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_ = userview.UserPage(&userview.UserPageProps{
		Ctx:  ctx,
		User: *user,
	}).Render(w)
}

func (h *UserHandler) AddUserPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	_ = userview.AddUserPage(&userview.AddUserPageProps{
		Ctx: ctx,
	}).Render(w)
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
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

	var newUser model.NewUser
	err = appurl.Unmarshal(r.Form, &newUser)

	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	validationErrors, err := h.userService.CreateUser(r.Context(), newUser)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding user", http.StatusInternalServerError)
		return
	}

	if len(validationErrors) > 0 {
		_ = userview.AddUserPage(&userview.AddUserPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *UserHandler) AddAPIUserPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = userview.AddAPIUserPage(&userview.AddAPIUserPageProps{
		Ctx: ctx,
	}).Render(w)
}

func (h *UserHandler) AddAPIUser(w http.ResponseWriter, r *http.Request) {
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

	var newAPIUser model.NewAPIUser
	err = appurl.Unmarshal(r.Form, &newAPIUser)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	validationErrors, password, err := h.userService.CreateAPIUser(r.Context(), newAPIUser)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding API user", http.StatusInternalServerError)
	}

	if len(validationErrors) > 0 {
		_ = userview.AddAPIUserPage(&userview.AddAPIUserPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	_ = userview.APIUserCredentialsPage(&userview.APIUserCredentialsPageProps{
		Ctx:      ctx,
		Username: newAPIUser.Username,
		Password: password,
	}).Render(w)
}

func (h *UserHandler) EditUserPage(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	_ = userview.EditUserPage(&userview.EditUserPageProps{
		User: *user,
		Ctx:  ctx,
	}).Render(w)
}

func (h *UserHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var userUpdate model.UserUpdate
	err = appurl.Unmarshal(r.Form, &userUpdate)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	validationErrors, err := h.userService.UpdateUser(r.Context(), userID, userUpdate)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
	}
	if len(validationErrors) > 0 {
		_ = userview.EditUserPage(&userview.EditUserPageProps{
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (h *UserHandler) ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusBadRequest)
		return
	}

	_ = userview.ResetPasswordPage(&userview.ResetPasswordPageProps{
		User: *user,
		Ctx:  ctx,
	}).Render(w)
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var passwordReset model.PasswordReset
	err = appurl.Unmarshal(r.Form, &passwordReset)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	user, validationErrors, err := h.userService.ResetPassword(
		r.Context(),
		userID,
		passwordReset,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error resetting password", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User doesn't exist", http.StatusBadRequest)
		return
	}

	if len(validationErrors) > 0 {
		_ = userview.ResetPasswordPage(&userview.ResetPasswordPageProps{
			Ctx:              ctx,
			User:             *user,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	credentials := model.VerifyPasswordLoginInput{
		Username: user.Username,
		Password: passwordReset.Password,
	}

	encoded, err := encryptcredentials.Encrypt(credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d/reset-password?EncryptedCredentials=%s", userID, encoded), http.StatusSeeOther)
}
