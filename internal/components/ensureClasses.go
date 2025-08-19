package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
)

// Helper function to find or add classes
func ensureClasses(children []g.Node, additionalClasses c.Classes) []g.Node {
	var foundClasses *c.Classes

	// Iterate through the children to find an instance of c.Classes
	for _, child := range children {
		if class, ok := child.(c.Classes); ok {
			foundClasses = &class
			break
		}
	}

	if foundClasses != nil {
		// Merge additional classes into the found classes
		for k, v := range additionalClasses {
			(*foundClasses)[k] = v
		}
	} else {
		// Append new classes if none found
		children = append(children, additionalClasses)
	}

	return children
}
