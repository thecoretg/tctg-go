package automate

import (
	"strings"
	"time"
)

// Time is a timestamp returned by the ConnectWise Automate API. Automate emits
// timestamps without a timezone offset (e.g. "2026-06-26T15:52:53"), which the
// standard library's RFC3339-based time.Time JSON unmarshalling rejects with
// `cannot parse "" as "Z07:00"`. Time accepts the zone-less layout (with or
// without fractional seconds) and falls back to RFC3339; zone-less values are
// interpreted as UTC. It embeds time.Time, so all the usual methods are
// available, and its promoted IsZero satisfies the `omitzero` JSON option.
type Time struct {
	time.Time
}

// automateLayout matches Automate's zone-less timestamps; the trailing
// fractional-second component is optional when parsing.
const automateLayout = "2006-01-02T15:04:05.999999999"

func (t *Time) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		t.Time = time.Time{}
		return nil
	}
	// Zone-bearing values (RFC3339) take precedence; fall back to Automate's
	// zone-less layout interpreted as UTC.
	if parsed, err := time.Parse(time.RFC3339, s); err == nil {
		t.Time = parsed
		return nil
	}
	parsed, err := time.Parse(automateLayout, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.Format(automateLayout) + `"`), nil
}
