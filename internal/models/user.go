package models

import (
	"app/pkg/appsort"
	"app/pkg/pgconv"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserDB struct {
	UserID      int
	Username    string
	IsAPIUser   bool
	Email       pgtype.Text
	FirstName   pgtype.Text
	LastName    pgtype.Text
	Created     time.Time
	LastLogin   pgtype.Timestamptz
	Permissions json.RawMessage
}

type User struct {
	UserID      int
	Username    string
	IsAPIUser   bool
	Email       *string
	FirstName   *string
	LastName    *string
	Created     time.Time
	LastLogin   *time.Time
	Permissions UserPermissions
}

func (u UserDB) ToDomain() User {

	var permissions UserPermissions
	if u.Permissions != nil {
		err := json.Unmarshal(u.Permissions, &permissions)
		if err != nil {
			// If unmarshalling fails, leave the permissions empty or set a default
			permissions = UserPermissions{}
		}
	}

	return User{
		UserID:      u.UserID,
		Username:    u.Username,
		IsAPIUser:   u.IsAPIUser,
		Email:       pgconv.PGTextToStringPtr(u.Email),
		FirstName:   pgconv.PGTextToStringPtr(u.FirstName),
		LastName:    pgconv.PGTextToStringPtr(u.LastName),
		Created:     u.Created,
		LastLogin:   pgconv.PGTimestamptzToTimePtr(u.LastLogin),
		Permissions: permissions,
	}
}

type NewUser struct {
	Username        string
	Email           *string
	FirstName       string
	LastName        string
	Password        string
	ConfirmPassword string
	Permissions     UserPermissions
}

type NewAPIUser struct {
	Username    string
	Permissions UserPermissions
}

type UserUpdate struct {
	Username    string
	Email       *string
	FirstName   string
	LastName    string
	Permissions UserPermissions
}

type PasswordReset struct {
	Password        string
	ConfirmPassword string
}

type GetUsersQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int
}

var GetUsersSortableKeys = []string{
	"Username",
	"Email",
	"FirstName",
	"LastName",
	"Created",
	"LastLogin",
}
