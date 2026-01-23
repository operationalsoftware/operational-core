package pdfview

import (
	"app/pkg/printnode"
	"bytes"
	"encoding/json"
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

func printNodeStatusBox(status printnode.Status) g.Node {
	message := status.Message
	statusClass := "warning"

	if !status.Configured {
		message = "PrintNode API key is not configured. Set PRINTNODE_API_KEY to enable automated printing."
	} else if status.Reachable {
		statusClass = "success"
		switch {
		case status.AccountName != "" && status.AccountEmail != "":
			message = fmt.Sprintf("Connected to PrintNode as %s (%s).", status.AccountName, status.AccountEmail)
		case status.AccountEmail != "":
			message = fmt.Sprintf("Connected to PrintNode as %s.", status.AccountEmail)
		case status.AccountName != "":
			message = fmt.Sprintf("Connected to PrintNode as %s.", status.AccountName)
		default:
			message = "PrintNode connection is working."
		}
	} else {
		statusClass = "error"
		if message == "" {
			message = "Unable to reach PrintNode. Check the API key and network connectivity."
		}
	}

	return h.Div(
		c.Classes{
			"printnode-status-box": true,
			statusClass:            true,
		},
		h.Div(
			h.Class("printnode-status-title"),
			g.Text("PrintNode"),
		),
		h.P(g.Text(message)),
	)
}

func documentLinkCell(title string, fileURL *string) g.Node {
	if fileURL == nil || *fileURL == "" {
		return g.Text(title)
	}
	return h.A(
		h.Href(*fileURL),
		h.Target("_blank"),
		h.Rel("noopener noreferrer"),
		g.Text(title),
	)
}

func prettyJSON(input json.RawMessage) string {
	if len(input) == 0 {
		return ""
	}
	var out bytes.Buffer
	if err := json.Indent(&out, input, "", "  "); err != nil {
		return string(input)
	}
	return out.String()
}
