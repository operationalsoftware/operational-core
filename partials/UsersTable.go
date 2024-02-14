package partials

import (
	"operationalcore/components"

	g "github.com/maragudk/gomponents"
)

var columns = []components.TableColumn{
	{
		Name:     "Username",
		Key:      "Username",
		Sortable: true,
	},
	{
		Name:     "First Name",
		Key:      "FirstName",
		Sortable: true,
	},
	{
		Name: "Last Name",
		Key:  "LastName",
	},
	{
		Name: "Email",
		Key:  "Email",
	},
	{
		Name: "Created",
		Key:  "Created",
	},
	{
		Name: "Last Login",
		Key:  "LastLogin",
	},
}

func UsersTable() g.Node {
	return components.Table(&components.TableProps{
		Columns: columns,
		Data:    []components.TableRowRenderer{},
		HxGet:   "/users/table",
	})
}
