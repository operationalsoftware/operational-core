package userview

import (
	"app/internal/components"
	"app/internal/model"
	"reflect"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

func getPermissionDescription(module, permission string) string {
	// Use reflection to get the description of
	// Create an instance of UserPermissions
	permissions := model.UserPermissions{}
	permissionsType := reflect.TypeOf(permissions)
	moduleField, found := permissionsType.FieldByName(module)
	if found {
		permissionField, found := moduleField.Type.FieldByName(permission)
		if found {
			description := permissionField.Tag.Get("description")
			if description != "" {
				return description
			}
		}
	}

	return ""
}

func permissionsCheckboxesPartial(userPermissions model.UserPermissions) g.Node {

	return components.Fieldset(
		h.Label(g.Text("Permissions")),
		h.H4(h.Class("module-title"), g.Text("User Admin")),
		components.Checkbox(
			&components.CheckboxProps{
				Classes: c.Classes{
					"permission-checkbox": true,
				},
				Name:    "Permissions.UserAdmin.Access",
				Label:   getPermissionDescription("UserAdmin", "Access"),
				Checked: userPermissions.UserAdmin.Access,
				Value:   "true",
			},
		),

		components.InlineStyle("/internal/views/userview/permissions_checkboxes_partial.css"),
	)

}
