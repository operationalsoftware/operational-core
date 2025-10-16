package format

import "github.com/dromara/carbon/v2"

// formatSecondsIntoDuration converts seconds into a human-readable relative duration string.
// For non-positive values, it returns an en dash.
func FormatSecondsIntoDuration(seconds int) string {
	if seconds <= 0 {
		return "\u2013"
	}

	base := carbon.CreateFromTimestamp(0)
	target := carbon.CreateFromTimestamp(0).AddSeconds(seconds)
	return target.DiffAbsInString(base)
}
