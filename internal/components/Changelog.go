package components

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ChangelogEntry struct {
	ChangedAt        time.Time
	ChangeByUsername string
	IsCreation       bool
	Changes          map[string]any
}

type ChangelogProperty struct {
	FieldKey string
	Label    g.Node
}

func Changelog(entries []ChangelogEntry, changelogProperties []ChangelogProperty) g.Node {
	return h.Div(
		h.Class("changelog"),
		h.H3(g.Text("Changelog")),
		g.If(len(entries) == 0,
			h.Div(
				h.Class("no-entries"),
				g.Text("No changes yet."),
			),
		),
		h.Ul(
			h.Class("main-list"),

			g.Group(g.Map(entries, func(entry ChangelogEntry) g.Node {
				return h.Li(
					h.Div(
						h.Class("heading"),

						h.Strong(
							g.If(
								entry.IsCreation,
								g.Text("Created"),
							),
							g.If(
								!entry.IsCreation,
								g.Text("Changed"),
							),
							g.Text(fmt.Sprintf(" by %s at ", entry.ChangeByUsername)),
							h.Span(h.Class("local-datetime"), g.Text(entry.ChangedAt.Format(time.RFC3339))),
						),
					),
					changesList(entry.Changes, changelogProperties),
				)
			})),
		),
	)
}

func changesList(changes map[string]any, changeLogProperties []ChangelogProperty) g.Node {

	return h.Ul(
		g.Group(g.Map(changeLogProperties, func(clp ChangelogProperty) g.Node {
			if value, exists := changes[clp.FieldKey]; exists {
				return change(clp.Label, value)
			}
			return nil
		})),
	)
}

func change(fieldLabel g.Node, value any) g.Node {
	if value == nil {
		return nil
	}

	formattedValue := formatValue(value)
	if formattedValue != "" {
		return h.Li(
			h.B(fieldLabel),
			g.Textf(": %s", formattedValue),
		)
	}
	return nil
}

func formatValue(value any) string {

	switch v := value.(type) {
	case *string:
		if v == nil {
			return ""
		}
		return *v
	case *int:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%d", *v)
	case *int8:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%d", *v)
	case *int16:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%d", *v)
	case *int32:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%d", *v)
	case *int64:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%d", *v)
	case *float32:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%f", *v)
	case *float64:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%f", *v)
	case *bool:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%t", *v)
	case *time.Time:
		if v == nil {
			return ""
		}
		return v.Format(time.RFC3339)
	case *decimal.Decimal:
		if v == nil {
			return ""
		}
		return v.String()
	default:
		return ""
	}
}
