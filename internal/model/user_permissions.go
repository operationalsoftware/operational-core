package model

// This file defines the permissions for each module in the application.
type UserAdminPermissions struct {
	Access bool `description:"Able to manage users, permissions and teams"`
}

type AndonPermissions struct {
	Admin bool `description:"Able to manage andon issues and structure (groups)"`
}

type StockItemPermissions struct {
	Admin bool `description:"Able to manage stock items"`
}

type SupplyChainPermissions struct {
	Admin      bool `description:"Administrative supply chain tasks"`
	TeamMember bool `description:"General supply chain tasks"`
}

// Finally, group under the UserPermissions struct
type UserPermissions struct {
	Andon       AndonPermissions       `description:"Andon"`
	Stock       StockItemPermissions   `description:"Stock"`
	SupplyChain SupplyChainPermissions `description:"Supply Chain"`
	UserAdmin   UserAdminPermissions   `description:"User Admin"`
}
