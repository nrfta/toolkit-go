package duration

import (
	"fmt"
	"strings"
	"time"

	"github.com/sosodev/duration"
)

// ParseToAbsoluteTime parses an ISO 8601 duration string relative to the provided time.
// Returns (parsedTime, isDuration, error).
//
// Supports:
//   - Positive durations: "P2W" (2 weeks from now), "P30D" (30 days from now)
//   - Negative durations: "-P1M" (1 month ago), "-P2W" (2 weeks ago)
//   - Complex durations: "P1Y2M3DT4H5M6S"
//
// The isDuration return value is true if the input was a valid duration,
// false otherwise. This allows callers to distinguish between duration parsing
// and passthrough of absolute dates.
func ParseToAbsoluteTime(value string, now time.Time) (time.Time, bool, error) {
	if !IsDuration(value) {
		return time.Time{}, false, fmt.Errorf("not a duration format: %s", value)
	}

	// Handle negative durations by detecting "-" prefix
	isNegative := strings.HasPrefix(value, "-")
	durationStr := value
	if isNegative {
		durationStr = value[1:] // Remove "-" prefix for parsing
	}

	// Parse the duration using sosodev/duration
	d, err := duration.Parse(durationStr)
	if err != nil {
		return time.Time{}, true, fmt.Errorf("invalid ISO 8601 duration: %w", err)
	}

	// Convert duration to time.Time by adding to now
	// The library returns a Duration struct with Years, Months, Weeks, Days, Hours, Minutes, Seconds
	var result time.Time
	if isNegative {
		// Subtract the duration components
		result = now.AddDate(-int(d.Years), -int(d.Months), -int(d.Weeks*7+d.Days))
		result = result.Add(-time.Duration(d.Hours) * time.Hour)
		result = result.Add(-time.Duration(d.Minutes) * time.Minute)
		result = result.Add(-time.Duration(d.Seconds) * time.Second)
	} else {
		// Add the duration components
		result = now.AddDate(int(d.Years), int(d.Months), int(d.Weeks*7+d.Days))
		result = result.Add(time.Duration(d.Hours) * time.Hour)
		result = result.Add(time.Duration(d.Minutes) * time.Minute)
		result = result.Add(time.Duration(d.Seconds) * time.Second)
	}

	return result, true, nil
}

// ParseOrPassthrough attempts to parse the value as an ISO 8601 duration.
// If it's a duration, it converts it to an absolute RFC3339 timestamp.
// If it's not a duration, it returns the value unchanged (assuming it's already RFC3339).
//
// This is the primary function used by SQLBoiler mods to handle both
// absolute dates and relative durations transparently.
func ParseOrPassthrough(value string, now time.Time) (string, error) {
	if !IsDuration(value) {
		// Not a duration, assume it's already an RFC3339 date
		return value, nil
	}

	// Parse as duration and convert to RFC3339
	absoluteTime, _, err := ParseToAbsoluteTime(value, now)
	if err != nil {
		return "", err
	}

	return absoluteTime.Format(time.RFC3339), nil
}

// IsDuration returns true if the string appears to be an ISO 8601 duration.
// Checks for the "P" prefix which is required for all ISO 8601 durations,
// including support for negative durations with "-P" prefix.
func IsDuration(value string) bool {
	if value == "" {
		return false
	}

	// Check for duration format: starts with "P" or "-P"
	return strings.HasPrefix(value, "P") || strings.HasPrefix(value, "-P")
}
