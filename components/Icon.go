package components

import (
	"io"

	g "github.com/maragudk/gomponents"
)

func Icon(identifier string) g.Node {
	file, err := Assets.Open("/icon-svgs/" + identifier + ".svg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Convert the byte slice to a string
	svgString := string(fileBytes)

	return g.Raw(svgString)
}
