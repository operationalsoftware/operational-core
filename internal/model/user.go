package model

import (
	"app/pkg/appsort"
	"app/pkg/pgconv"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserDB struct {
	UserID                 int
	Username               string
	IsAPIUser              bool
	Email                  pgtype.Text
	FirstName              pgtype.Text
	LastName               pgtype.Text
	Created                time.Time
	LastLogin              pgtype.Timestamptz
	Permissions            json.RawMessage
	SessionDurationMinutes *int
}

type User struct {
	UserID                 int
	Username               string `sortable:"true"`
	IsAPIUser              bool
	Email                  *string    `sortable:"true"`
	FirstName              *string    `sortable:"true"`
	LastName               *string    `sortable:"true"`
	Created                time.Time  `sortable:"true"`
	LastLogin              *time.Time `sortable:"true"`
	Permissions            UserPermissions
	SessionDurationMinutes *int
	Teams                  []UserTeam `json:"teams"`
}

func (u UserDB) ToDomain() User {

	var permissions UserPermissions
	if u.Permissions != nil {
		err := json.Unmarshal(u.Permissions, &permissions)
		if err != nil {
			permissions = UserPermissions{}
		}
	}

	return User{
		UserID:                 u.UserID,
		Username:               u.Username,
		IsAPIUser:              u.IsAPIUser,
		Email:                  pgconv.PGTextToStringPtr(u.Email),
		FirstName:              pgconv.PGTextToStringPtr(u.FirstName),
		LastName:               pgconv.PGTextToStringPtr(u.LastName),
		Created:                u.Created,
		LastLogin:              pgconv.PGTimestamptzToTimePtr(u.LastLogin),
		Permissions:            permissions,
		SessionDurationMinutes: u.SessionDurationMinutes,
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
	Password    string
	Permissions UserPermissions
}

type UserUpdate struct {
	Username               string
	Email                  *string
	FirstName              string
	LastName               string
	Permissions            UserPermissions
	SessionDurationMinutes *int
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

type UserTeam struct {
	TeamID   int    `json:"team_id"`
	UserID   int    `json:"user_id"`
	TeamName string `sortable:"true" json:"team_name"`
	Role     string `sortable:"true" json:"role"`
}

type ListUserTeamsQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int
}
