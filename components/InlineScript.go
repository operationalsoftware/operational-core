package components

import (
	"io"

	"github.com/jessevdk/go-assets"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func InlineScript(fs *assets.FileSystem, assetPath string) g.Node {
	file, err := fs.Open(assetPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Convert the byte slice to a string
	fileContent := string(fileBytes)

	return Script(
		Type("text/javascript"),
		g.Raw(fileContent),
	)
}
