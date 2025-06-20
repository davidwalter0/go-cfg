package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DurationValue implements pflag.Value for time.Duration with custom formats
type DurationValue struct {
	Duration *time.Duration
}

func NewDurationValue(defaultDuration time.Duration) *DurationValue {
	return &DurationValue{Duration: &defaultDuration}
}

func (dv *DurationValue) String() string {
	if dv.Duration == nil {
		return ""
	}
	return dv.Duration.String()
}

func (dv *DurationValue) Set(s string) error {
	if s == "" {
		dv.Duration = nil
		return nil
	}
	
	// Try standard Go duration parsing first
	if d, err := time.ParseDuration(s); err == nil {
		dv.Duration = &d
		return nil
	}
	
	// Try parsing as integer seconds
	if seconds, err := strconv.Atoi(s); err == nil {
		d := time.Duration(seconds) * time.Second
		dv.Duration = &d
		return nil
	}
	
	// Try parsing as float seconds
	if seconds, err := strconv.ParseFloat(s, 64); err == nil {
		d := time.Duration(seconds * float64(time.Second))
		dv.Duration = &d
		return nil
	}
	
	return fmt.Errorf("invalid duration format: %s (expected: 1h30m, 90s, 90, or 90.5)", s)
}

func (dv *DurationValue) Type() string {
	return "duration"
}

func (dv *DurationValue) Get() interface{} {
	if dv.Duration == nil {
		return nil
	}
	return *dv.Duration
}

func (dv *DurationValue) GetDuration() time.Duration {
	if dv.Duration == nil {
		return 0
	}
	return *dv.Duration
}

func (dv *DurationValue) IsSet() bool {
	return dv.Duration != nil
}

// JSON marshaling as seconds (float)
func (dv DurationValue) MarshalJSON() ([]byte, error) {
	if dv.Duration == nil {
		return json.Marshal(nil)
	}
	seconds := dv.Duration.Seconds()
	return json.Marshal(seconds)
}

func (dv *DurationValue) UnmarshalJSON(data []byte) error {
	var seconds *float64
	if err := json.Unmarshal(data, &seconds); err != nil {
		return err
	}
	
	if seconds == nil {
		dv.Duration = nil
		return nil
	}
	
	d := time.Duration(*seconds * float64(time.Second))
	dv.Duration = &d
	return nil
}

// YAML marshaling as human-readable string
func (dv DurationValue) MarshalYAML() (interface{}, error) {
	if dv.Duration == nil {
		return nil, nil
	}
	return dv.Duration.String(), nil
}

func (dv *DurationValue) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}
	return dv.Set(s)
}

// HumanDurationValue implements pflag.Value with human-friendly parsing
type HumanDurationValue struct {
	Duration *time.Duration
}

func NewHumanDurationValue(defaultDuration time.Duration) *HumanDurationValue {
	return &HumanDurationValue{Duration: &defaultDuration}
}

func (hdv *HumanDurationValue) String() string {
	if hdv.Duration == nil {
		return ""
	}
	return hdv.humanFormat(*hdv.Duration)
}

func (hdv *HumanDurationValue) humanFormat(d time.Duration) string {
	if d == 0 {
		return "0s"
	}
	
	var parts []string
	
	// Days
	days := d / (24 * time.Hour)
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
		d -= days * 24 * time.Hour
	}
	
	// Hours
	hours := d / time.Hour
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
		d -= hours * time.Hour
	}
	
	// Minutes
	minutes := d / time.Minute
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
		d -= minutes * time.Minute
	}
	
	// Seconds
	if d > 0 || len(parts) == 0 {
		if d >= time.Second {
			seconds := d / time.Second
			parts = append(parts, fmt.Sprintf("%ds", seconds))
		} else {
			// Milliseconds or less
			ms := d / time.Millisecond
			if ms > 0 {
				parts = append(parts, fmt.Sprintf("%dms", ms))
			} else {
				parts = append(parts, d.String())
			}
		}
	}
	
	return strings.Join(parts, "")
}

func (hdv *HumanDurationValue) Set(s string) error {
	if s == "" {
		hdv.Duration = nil
		return nil
	}
	
	// Try standard duration parsing first
	if d, err := time.ParseDuration(s); err == nil {
		hdv.Duration = &d
		return nil
	}
	
	// Try parsing human formats like "1 day", "2 hours 30 minutes"
	d, err := hdv.parseHumanFormat(s)
	if err != nil {
		return fmt.Errorf("invalid duration format: %s", s)
	}
	
	hdv.Duration = &d
	return nil
}

func (hdv *HumanDurationValue) parseHumanFormat(s string) (time.Duration, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}
	
	var total time.Duration
	
	// Split by whitespace and process each part
	parts := strings.Fields(s)
	
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		
		// Try to extract number and unit
		var num float64
		
		// Look for number at start of part
		numStr := ""
		unitStr := part
		
		for j, r := range part {
			if (r >= '0' && r <= '9') || r == '.' {
				numStr += string(r)
			} else {
				unitStr = part[j:]
				break
			}
		}
		
		if numStr == "" {
			// Maybe number is separate (e.g., "1 hour")
			if i+1 < len(parts) {
				var err error
				num, err = strconv.ParseFloat(parts[i], 64)
				if err == nil {
					i++ // consume next part as unit
					unitStr = parts[i]
				} else {
					return 0, fmt.Errorf("invalid number: %s", parts[i])
				}
			} else {
				return 0, fmt.Errorf("no number found in: %s", part)
			}
		} else {
			var err error
			num, err = strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", numStr)
			}
		}
		
		// Parse unit
		var multiplier time.Duration
		switch unitStr {
		case "ns", "nanosecond", "nanoseconds":
			multiplier = time.Nanosecond
		case "us", "μs", "microsecond", "microseconds":
			multiplier = time.Microsecond
		case "ms", "millisecond", "milliseconds":
			multiplier = time.Millisecond
		case "s", "sec", "second", "seconds":
			multiplier = time.Second
		case "m", "min", "minute", "minutes":
			multiplier = time.Minute
		case "h", "hr", "hour", "hours":
			multiplier = time.Hour
		case "d", "day", "days":
			multiplier = 24 * time.Hour
		case "w", "week", "weeks":
			multiplier = 7 * 24 * time.Hour
		default:
			return 0, fmt.Errorf("unknown unit: %s", unitStr)
		}
		
		total += time.Duration(num * float64(multiplier))
	}
	
	return total, nil
}

func (hdv *HumanDurationValue) Type() string {
	return "humanDuration"
}

func (hdv *HumanDurationValue) Get() interface{} {
	if hdv.Duration == nil {
		return nil
	}
	return *hdv.Duration
}

func (hdv *HumanDurationValue) GetDuration() time.Duration {
	if hdv.Duration == nil {
		return 0
	}
	return *hdv.Duration
}