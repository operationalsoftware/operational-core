package format

import (
	"fmt"

	"github.com/dromara/carbon/v2"
)

// FormatSecondsIntoDuration converts seconds into "00d 00h 00m" using carbon.
func FormatSecondsIntoDuration(seconds int) string {
	if seconds <= 0 {
		return ""
	}

	base := carbon.CreateFromTimestamp(0)
	target := carbon.CreateFromTimestamp(int64(seconds))

	days := target.DiffAbsInDays(base)
	hours := target.DiffAbsInHours(base) - days*24
	minutes := target.DiffAbsInMinutes(base) - (days*24+hours)*60
	remainderSeconds := target.DiffAbsInSeconds(base) - ((days*24+hours)*60+minutes)*60

	return fmt.Sprintf("%02dd %02dh %02dm %02ds", days, hours, minutes, remainderSeconds)
}

// FormatSecondsIntoMinutes converts seconds into whole minutes.
func FormatSecondsIntoMinutes(seconds int) string {
	if seconds <= 0 {
		return ""
	}

	base := carbon.CreateFromTimestamp(0)
	target := carbon.CreateFromTimestamp(int64(seconds))
	minutes := target.DiffAbsInMinutes(base)

	return fmt.Sprintf("%d", minutes)
}

// FormatOptionalSecondsIntoMinutes returns minutes and tooltip strings for nullable durations.
func FormatOptionalSecondsIntoMinutes(seconds *int) (string, string) {
	if seconds == nil {
		return "\u2013", "\u2013"
	}

	return FormatSecondsIntoMinutes(*seconds), FormatSecondsIntoDuration(*seconds)
}
