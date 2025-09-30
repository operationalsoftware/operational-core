package userview

import (
	"app/internal/components"
	"app/internal/model"
	"fmt"
	"reflect"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// Data structure for permissions
type permissionModule struct {
	name        string
	description string
	permissions []permission
}

type permission struct {
	name        string
	description string
	value       bool
	inputName   string // For form binding
}

// For add/edit user pages - shows all permissions with checkboxes
func permissionsCheckboxesPartial(userPermissions model.UserPermissions) g.Node {
	modules := extractPermissionStructure(userPermissions)

	return h.Div(
		components.InlineStyle("/internal/views/userview/permissions_partials.css"),
		g.Text("Permissions"),
		h.Div(
			h.Class("fieldset"),
			g.Group(renderPermissionsAsCheckboxes(modules)),
		),
	)
}

// For user view page - shows only granted permissions
func permissionsDisplayPartial(userPermissions model.UserPermissions) g.Node {
	modules := extractPermissionStructure(userPermissions)
	nodes := renderPermissionsAsDisplay(modules)

	if len(nodes) == 0 {
		nodes = []g.Node{h.P(g.Text("No permissions granted"))}
	}

	return g.Group([]g.Node{
		components.InlineStyle("/internal/views/userview/permissions_partials.css"),
		h.Div(
			h.Class("permissions-display"),
			g.Group(nodes),
		),
	})
}

// Extract permissions into a structured format
func extractPermissionStructure(userPermissions model.UserPermissions) []permissionModule {
	var modules []permissionModule
	permissionsVal := reflect.ValueOf(userPermissions)
	permissionsType := reflect.TypeOf(userPermissions)

	for i := 0; i < permissionsVal.NumField(); i++ {
		module := permissionsVal.Field(i)
		moduleType := permissionsType.Field(i)
		moduleName := moduleType.Name
		moduleDescription := moduleType.Tag.Get("description")

		var permissions []permission
		for j := 0; j < module.NumField(); j++ {
			permField := module.Field(j)
			permFieldType := module.Type().Field(j)
			permName := permFieldType.Name
			permDescription := permFieldType.Tag.Get("description")
			permValue := permField.Bool()
			inputName := fmt.Sprintf("Permissions.%s.%s", moduleName, permName)

			permissions = append(permissions, permission{
				name:        permName,
				description: permDescription,
				value:       permValue,
				inputName:   inputName,
			})
		}

		modules = append(modules, permissionModule{
			name:        moduleName,
			description: moduleDescription,
			permissions: permissions,
		})
	}

	return modules
}

// Render as editable checkboxes
func renderPermissionsAsCheckboxes(modules []permissionModule) []g.Node {
	var moduleNodes []g.Node

	for _, module := range modules {
		var permissionNodes []g.Node
		for _, perm := range module.permissions {
			permissionNode := h.Div(
				h.Class("permission"),
				h.Label(
					h.Class("label-checkbox"),
					g.Text(perm.name),
					h.Input(
						h.Class("checkbox"),
						h.Type("checkbox"),
						h.Name(perm.inputName),
						h.Value("true"),
						g.If(perm.value, h.Checked()),
					),
				),
				h.Div(
					h.Class("description"),
					g.Text(perm.description),
				),
			)

			permissionNodes = append(permissionNodes, permissionNode)
		}

		moduleNode := h.Div(
			h.Class("module"),
			h.H4(h.Class("title"), g.Text(module.description)),
			g.Group(permissionNodes),
		)

		moduleNodes = append(moduleNodes, moduleNode)
	}

	return moduleNodes
}

// Render as read-only display (only granted permissions)
func renderPermissionsAsDisplay(modules []permissionModule) []g.Node {
	var moduleNodes []g.Node

	for _, module := range modules {
		var grantedPerms []g.Node

		for _, perm := range module.permissions {
			if perm.value {
				grantedPerms = append(grantedPerms,
					h.Div(
						h.Class("permission"),
						h.Div(g.Text(perm.name)),
						h.Div(h.Class("description"), g.Text(perm.description)),
					),
				)
			}
		}

		// Only show module if it has granted permissions
		if len(grantedPerms) > 0 {
			moduleNode := h.Div(
				h.Class("module"),
				h.H4(
					h.Class("title"),
					g.Text(module.description),
				),
				g.Group(grantedPerms),
			)
			moduleNodes = append(moduleNodes, moduleNode)
		}
	}

	return moduleNodes
}
