package model

// This file defines the permissions for each module in the application.
type UserAdminPermissions struct {
	Access bool `description:"Able to manage users and permissions"`
}

type SupplyChainPermissions struct {
	Admin      bool `description:"Administrative supply chain tasks"`
	TeamMember bool `description:"General supply chain tasks"`
}

// Finally, group under the UserPermissions struct
type UserPermissions struct {
	UserAdmin   UserAdminPermissions
	SupplyChain SupplyChainPermissions
}
