package model

// This file defines the permissions for each module in the application.
type UserAdminPermissions struct {
	Access bool `description:"Able to manage users and permissions"`
}

type ProductionPermissions struct {
	Admin      bool `description:"Administrative tasks relating to production"`
	TeamMember bool `description:"General tasks relating to production"`
}

type SupplyChainPermissions struct {
	Admin      bool `description:"Administrative supply chain tasks such as configuring KANBAN"`
	TeamMember bool `description:"General supply chain tasks such as picking and moving KANBAN"`
}

// Finally, group under the UserPermissions struct
type UserPermissions struct {
	UserAdmin   UserAdminPermissions
	Production  ProductionPermissions
	SupplyChain SupplyChainPermissions
}
