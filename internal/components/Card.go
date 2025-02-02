package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func Card(children ...g.Node) g.Node {

	var foundClasses *c.Classes

	// Iterate through the children to find an instance of c.Classes
	for _, child := range children {
		if class, ok := child.(c.Classes); ok {
			foundClasses = &class
			break
		}
	}

	if foundClasses != nil {
		(*foundClasses)["card"] = true // ensure card class is set
	} else {
		children = append(children, c.Classes{
			"card": true,
		})
	}

	return h.Div(children...)
}
