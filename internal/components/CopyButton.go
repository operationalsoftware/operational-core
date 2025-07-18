package components

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type CopyButtonProps struct {
	TextToCopy string // The actual text that will be copied
	ButtonID   string // A unique ID to identify the button and clipboard target
	Label      string // Optional: button label
}

// CopyButton returns a copy-to-clipboard button
func CopyButton(p CopyButtonProps) g.Node {
	buttonID := fmt.Sprintf("copy-btn-%s", p.ButtonID)
	statusID := fmt.Sprintf("copy-status-%s", p.ButtonID)

	return h.Div(
		h.Class("clipboard-container"),
		Button(&ButtonProps{
			Classes: c.Classes{
				"clipboard-btn": true,
			},
		},
			h.ID(buttonID),
			Icon(&IconProps{
				Identifier: "content-copy",
				Classes: c.Classes{
					"icon": true,
				},
			}),
		),

		h.Span(
			h.ID(statusID),
			h.Class("clipboard-status hidden"),
			g.Text("Copied!"),
		),

		h.Script(g.Raw(fmt.Sprintf(`
			document.addEventListener("DOMContentLoaded", function () {
				const btn = document.getElementById("%s");
				const status = document.getElementById("%s");

				if (btn && status) {
					btn.addEventListener("click", function () {
						navigator.clipboard.writeText(%q)
							.then(() => {
								status.classList.remove("hidden");
								setTimeout(() => {
									status.classList.add("hidden");
								}, 2000);
							})
							.catch(err => {
								console.error("Clipboard copy failed", err);
							});
					});
				}
			});
		`, buttonID, statusID, p.TextToCopy))),
	)
}
