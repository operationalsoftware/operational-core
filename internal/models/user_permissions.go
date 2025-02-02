package models

// This file defines the permissions for each module in the application.
type UserAdminPermissions struct {
	Access bool `description:"Able to manage users and permissions"`
}

// Finally, group under the UserPermissions struct
type UserPermissions struct {
	UserAdmin UserAdminPermissions
}
