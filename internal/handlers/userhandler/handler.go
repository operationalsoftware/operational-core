package userhandler

import (
	"app/internal/model"
	"app/internal/services/userservice"
	"app/internal/views/userview"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"app/pkg/urlvalues"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService userservice.UserService
}

func NewUserHandler(userService userservice.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) UsersHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	type urlVals struct {
		Sort     string
		Page     int
		PageSize int
	}

	var uv urlVals

	err := urlvalues.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	sort := appsort.Sort{}
	sort.ParseQueryParam(uv.Sort, model.GetUsersSortableKeys)

	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}

	users, err := h.userService.GetUsers(r.Context(), model.GetUsersQuery{
		Sort:     sort,
		Page:     uv.Page,
		PageSize: uv.PageSize,
	})
	if err != nil {
		http.Error(w, "Error listing users", http.StatusInternalServerError)
		return
	}
	count, err := h.userService.GetUserCount(r.Context())
	if err != nil {
		http.Error(w, "Error counting users", http.StatusInternalServerError)
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
		Id:   id,
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
	err = urlvalues.Unmarshal(r.Form, &newUser)

	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := h.userService.ValidateNewUser(r.Context(), newUser)
	if !valid {
		_ = userview.AddUserPage(&userview.AddUserPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	err = h.userService.CreateUser(r.Context(), newUser)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding user", http.StatusInternalServerError)
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
	err = urlvalues.Unmarshal(r.Form, &newAPIUser)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := h.userService.ValidateNewAPIUser(r.Context(), newAPIUser)

	if !valid {
		fmt.Println(valid, validationErrors)
		_ = userview.AddAPIUserPage(&userview.AddAPIUserPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	password, err := h.userService.CreateAPIUser(r.Context(), newAPIUser)

	if err != nil {
		http.Error(w, "Error adding API user", http.StatusInternalServerError)
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
		http.Error(w, "Error getting user", http.StatusBadRequest)
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

	var userUpdate model.UserUpdate
	err = urlvalues.Unmarshal(r.Form, &userUpdate)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := h.userService.ValidateUserUpdate(r.Context(), userUpdate)
	if !valid {
		_ = userview.EditUserPage(&userview.EditUserPageProps{
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	err = h.userService.UpdateUser(r.Context(), id, userUpdate)

	if err != nil {
		log.Panic(err)
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
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

	var passwordReset model.PasswordReset
	err = urlvalues.Unmarshal(r.Form, &passwordReset)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	valid, validationErrors := h.userService.ValidatePasswordReset(passwordReset)
	if !valid {
		user, err := h.userService.GetUserByID(r.Context(), id)
		if err != nil {
			http.Error(w, "Error getting user", http.StatusBadRequest)
			return
		}

		_ = userview.ResetPasswordPage(&userview.ResetPasswordPageProps{
			Ctx:              ctx,
			User:             *user,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	err = h.userService.ResetPassword(r.Context(), id, passwordReset)

	if err != nil {
		http.Error(w, "Error resetting password", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d", id), http.StatusSeeOther)
}
