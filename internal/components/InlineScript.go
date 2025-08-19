package components

import (
	"app/assets"
	"fmt"
	"io"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
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
