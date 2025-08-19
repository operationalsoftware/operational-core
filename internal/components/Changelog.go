package components

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ChangelogEntry struct {
	ChangedAt         time.Time
	ChangedByUsername string
	IsCreation        bool
	Changes           map[string]interface{}
}

type ChangelogFieldDefinition struct {
	Name  string
	Label string
}

func Changelog(entries []ChangelogEntry, fieldDefs []ChangelogFieldDefinition) g.Node {
	return h.Div(
		h.Class("changelog"),
		h.H3(g.Text("Changelog")),
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
							g.Text(fmt.Sprintf(" by %s at ", entry.ChangedByUsername)),
							h.Span(h.Class("local-datetime"), g.Text(entry.ChangedAt.Format(time.RFC3339))),
						),
					),
					renderChanges(entry.Changes, fieldDefs),
				)
			})),
		),
	)
}

func renderChanges(changes map[string]interface{}, fieldDefs []ChangelogFieldDefinition) g.Node {

	return h.Ul(
		g.Group(g.Map(fieldDefs, func(fd ChangelogFieldDefinition) g.Node {
			if value, exists := changes[fd.Name]; exists {
				return renderChange(fd.Label, value)
			}
			return nil
		})),
	)
}

func renderChange(fieldLabel string, value interface{}) g.Node {
	formattedValue := formatValue(value)
	if formattedValue != "" {
		return h.Li(g.Text(fmt.Sprintf("%s: %s", fieldLabel, formattedValue)))
	}
	return nil
}

func formatValue(value interface{}) string {
	v := reflect.ValueOf(value)

	switch v.Type() {
	case reflect.TypeOf(pgtype.Text{}):
		if v.Interface().(pgtype.Text).Valid {
			return v.Interface().(pgtype.Text).String
		}
	case reflect.TypeOf(sql.NullInt16{}):
		if v.Interface().(sql.NullInt16).Valid {
			return fmt.Sprintf("%d", v.Interface().(sql.NullInt16).Int16)
		}
	case reflect.TypeOf(sql.NullInt32{}):
		if v.Interface().(sql.NullInt32).Valid {
			return fmt.Sprintf("%d", v.Interface().(sql.NullInt32).Int32)
		}
	case reflect.TypeOf(sql.NullInt64{}):
		if v.Interface().(sql.NullInt64).Valid {
			return fmt.Sprintf("%d", v.Interface().(sql.NullInt64).Int64)
		}
	case reflect.TypeOf(sql.NullFloat64{}):
		if v.Interface().(sql.NullFloat64).Valid {
			return fmt.Sprintf("%f", v.Interface().(sql.NullFloat64).Float64)
		}
	case reflect.TypeOf(sql.NullBool{}):
		if v.Interface().(sql.NullBool).Valid {
			return fmt.Sprintf("%t", v.Interface().(sql.NullBool).Bool)
		}
	case reflect.TypeOf(decimal.NullDecimal{}):
		if v.Interface().(decimal.NullDecimal).Valid {
			return v.Interface().(decimal.NullDecimal).Decimal.String()
		}
	case reflect.TypeOf(pgtype.Timestamptz{}):
		if v.Interface().(pgtype.Timestamptz).Valid {
			return v.Interface().(pgtype.Timestamptz).Time.Format(time.RFC3339)
		}
	}

	return ""
}
