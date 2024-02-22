package userModel

// This file defines the roles for each module in the application.

type UserAdminRoles struct {
	Access bool
}

// Finally, group under the UserRoles struct
type UserRoles struct {
	UserAdmin UserAdminRoles
}
