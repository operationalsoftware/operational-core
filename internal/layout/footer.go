package layout

import (
	"app/assets"
	"fmt"
	"io"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func footer() g.Node {
	assetPath := "/internal/layout/build_version.txt"
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

	return Footer(
		Div(
			Class("footer-content"),
			Div(Class("footer-left")), // flexible spacer
			Div(
				Class("footer-center"),
				g.Text("An OperationalPlatform"),
				Sup(g.Text("TM")),
				g.Text(" by "),
				A(
					Href("https://operationalsoftware.co"),
					Target("_blank"),
					g.Text("Operational Software"),
				),
			),
			Div(
				Class("footer-right"),
				Code(
					Class("commit-hash"),
					g.Text(fmt.Sprintf("build: %s", fileContent)),
				),
			),
		),
	)

}
