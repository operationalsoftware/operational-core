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

type PrintingPermissions struct {
	// Printing Operator
	Operator bool `description:"Able to view print jobs and reprint PDFs"`

	// Printing Admin
	Admin bool `description:"Able to manage printer assignments"`
}

type AutomationPermissions struct {
	// Automation Admin
	AutomationAdmin bool `description:"Able to manage SQL Actions and use test functions"`

	PrinterAssignmentsEditor bool `description:"Able to edit printer assignments"`
}

// Finally, group under the UserPermissions struct
type UserPermissions struct {
	Andon       AndonPermissions       `description:"Andon"`
	Stock       StockItemPermissions   `description:"Stock"`
	SupplyChain SupplyChainPermissions `description:"Supply Chain"`
	UserAdmin   UserAdminPermissions   `description:"User Admin"`
	Automation  AutomationPermissions  `description:"Automation"`
	Printing    PrintingPermissions    `description:"Printing"`
}
