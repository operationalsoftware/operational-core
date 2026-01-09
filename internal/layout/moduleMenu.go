package layout

import (
	"app/internal/components"
	"app/internal/model"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

var AppMenu = []components.GridMenuGroup{
	{
		GroupName: "Admin",
		Items: []components.GridMenuItem{
			{
				Icon: "account-multiple",
				Name: "Users",
				Link: "/users",
				Show: func(permissions model.UserPermissions) bool {
					return permissions.UserAdmin.Access
				},
			},
			{
				Icon: "account-group",
				Name: "Teams",
				Link: "/teams",
				Show: func(permissions model.UserPermissions) bool {
					return permissions.UserAdmin.Access
				},
			},
		},
	},
	{
		GroupName: "Tools",
		Items: []components.GridMenuItem{
			{
				Icon: "alert-octagon-outline",
				Name: "Andons",
				Link: "/andons",
			},
		},
	},
	{
		GroupName: "Automation",
		Items: []components.GridMenuItem{
			{
				Icon: "text-box-outline",
				Name: "PDFs",
				Link: "/pdf",
				Show: func(permissions model.UserPermissions) bool {
					return permissions.Automation.AutomationAdmin
				},
			},
			{
				Icon: "printer-settings",
				Name: "Printing",
				Link: "/printing",
				Show: func(permissions model.UserPermissions) bool {
					return permissions.Automation.PrinterAssignmentsEditor
				},
			},
		},
	},
	{
		GroupName: "Stock",
		Items: []components.GridMenuItem{
			{
				Icon: "package-variant-closed",
				Name: "Stock",
				Link: "/stock",
				Show: func(permissions model.UserPermissions) bool {
					return true
				},
			},
			{
				Icon: "cube-outline",
				Name: "Stock Items",
				Link: "/stock-items",
				Show: func(permissions model.UserPermissions) bool {
					return true
				},
			},
		},
	},
	{
		GroupName: "Resource Management",
		Items: []components.GridMenuItem{
			{
				Icon: "cube-scan",
				Name: "Resources",
				Link: "/resources",
			},
			{
				Icon: "account-wrench",
				Name: "Servicing",
				Link: "/services",
			},
		},
	},
}

type moduleMenuProps struct {
	Ctx reqcontext.ReqContext
}

func moduleMenu(p *moduleMenuProps) g.Node {

	return h.Div(
		h.Button(
			h.ID("navbar-module-menu-button"),
			h.Class("menu-button"),
			components.Icon(&components.IconProps{
				Identifier: "dots",
			}),
		),
		// styled with position: fixed
		h.Div(
			h.Class("dropdown-panel"),
			h.ID("navbar-module-menu"),

			components.GridMenu(&components.GridMenuProps{
				Groups:          AppMenu,
				UserPermissions: p.Ctx.User.Permissions,
			}),
		),
		components.InlineScript("/internal/layout/moduleMenu.js"),
	)
}
