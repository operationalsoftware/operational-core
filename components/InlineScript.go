package components

import (
	"fmt"
	"io"
	"operationalcore/assets"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func InlineScript(assetPath string) g.Node {
	file, err := assets.Assets.Open(assetPath)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", assetPath, err)
		panic(err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Convert the byte slice to a string
	fileContent := string(fileBytes)

	return h.Script(
		h.Type("text/javascript"),
		g.Raw(fileContent),
	)
}
