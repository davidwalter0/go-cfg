package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// EnumValue implements pflag.Value for enum-like string constants with validation
type EnumValue struct {
	value       string
	validValues []string
	allowEmpty  bool
}

func NewEnumValue(validValues []string, defaultValue string, allowEmpty bool) *EnumValue {
	ev := &EnumValue{
		validValues: make([]string, len(validValues)),
		allowEmpty:  allowEmpty,
	}

	// Store valid values in lowercase for case-insensitive comparison
	for i, v := range validValues {
		ev.validValues[i] = strings.ToLower(v)
	}

	if defaultValue != "" {
		// If setting the default value fails, this is a programming error
		// We should panic or return an error, but for backward compatibility,
		// we'll ignore the error here
		_ = ev.Set(defaultValue)
	}

	return ev
}

func (ev *EnumValue) String() string {
	return ev.value
}

func (ev *EnumValue) Set(s string) error {
	if s == "" {
		if ev.allowEmpty {
			ev.value = ""
			return nil
		}
		return fmt.Errorf("empty value not allowed, valid options: %s",
			strings.Join(ev.validValues, ", "))
	}

	normalized := strings.ToLower(strings.TrimSpace(s))

	// Check if value is valid
	for _, valid := range ev.validValues {
		if normalized == valid {
			ev.value = s // Store original case
			return nil
		}
	}

	return fmt.Errorf("invalid value '%s', valid options: %s",
		s, strings.Join(ev.validValues, ", "))
}

func (ev *EnumValue) Type() string {
	return fmt.Sprintf("enum[%s]", strings.Join(ev.validValues, "|"))
}

func (ev *EnumValue) Get() interface{} {
	return ev.value
}

func (ev *EnumValue) GetValue() string {
	return ev.value
}

func (ev *EnumValue) IsValid(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	for _, valid := range ev.validValues {
		if normalized == valid {
			return true
		}
	}
	return false
}

func (ev *EnumValue) GetValidValues() []string {
	return append([]string(nil), ev.validValues...) // return copy
}

// JSON marshaling
func (ev EnumValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(ev.value)
}

func (ev *EnumValue) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return ev.Set(s)
}

// YAML marshaling
func (ev EnumValue) MarshalYAML() (interface{}, error) {
	return ev.value, nil
}

func (ev *EnumValue) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}
	return ev.Set(s)
}

// LogLevelValue is a specialized enum for log levels
type LogLevelValue struct {
	*EnumValue
}

func NewLogLevelValue(defaultLevel string) *LogLevelValue {
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	return &LogLevelValue{
		EnumValue: NewEnumValue(validLevels, defaultLevel, false),
	}
}

func (llv *LogLevelValue) IsDebugEnabled() bool {
	return strings.ToLower(llv.value) == "debug"
}

func (llv *LogLevelValue) IsInfoEnabled() bool {
	level := strings.ToLower(llv.value)
	return level == "debug" || level == "info"
}

func (llv *LogLevelValue) IsWarnEnabled() bool {
	level := strings.ToLower(llv.value)
	return level == "debug" || level == "info" || level == "warn"
}

func (llv *LogLevelValue) GetLevel() string {
	return llv.value
}

// EnvironmentValue is a specialized enum for deployment environments
type EnvironmentValue struct {
	*EnumValue
}

func NewEnvironmentValue(defaultEnv string) *EnvironmentValue {
	validEnvs := []string{"development", "dev", "staging", "stage", "production", "prod", "test", "testing"}
	return &EnvironmentValue{
		EnumValue: NewEnumValue(validEnvs, defaultEnv, false),
	}
}

func (env *EnvironmentValue) IsProduction() bool {
	level := strings.ToLower(env.value)
	return level == "production" || level == "prod"
}

func (env *EnvironmentValue) IsDevelopment() bool {
	level := strings.ToLower(env.value)
	return level == "development" || level == "dev"
}

func (env *EnvironmentValue) IsStaging() bool {
	level := strings.ToLower(env.value)
	return level == "staging" || level == "stage"
}

func (env *EnvironmentValue) IsTesting() bool {
	level := strings.ToLower(env.value)
	return level == "test" || level == "testing"
}

func (env *EnvironmentValue) GetEnvironment() string {
	return env.value
}

// MultiEnumValue allows multiple values from an enum set
type MultiEnumValue struct {
	values      []string
	validValues []string
	allowEmpty  bool
}

func NewMultiEnumValue(validValues []string, defaultValues []string, allowEmpty bool) *MultiEnumValue {
	mev := &MultiEnumValue{
		validValues: make([]string, len(validValues)),
		allowEmpty:  allowEmpty,
	}

	for i, v := range validValues {
		mev.validValues[i] = strings.ToLower(v)
	}

	for _, v := range defaultValues {
		if v != "" {
			_ = mev.Set(v)
		}
	}

	return mev
}

func (mev *MultiEnumValue) String() string {
	if len(mev.values) == 0 {
		return ""
	}
	return strings.Join(mev.values, ",")
}

func (mev *MultiEnumValue) Set(s string) error {
	if s == "" {
		if mev.allowEmpty {
			mev.values = nil
			return nil
		}
		return fmt.Errorf("empty value not allowed")
	}

	// Split by comma and validate each value
	parts := strings.Split(s, ",")
	var newValues []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		normalized := strings.ToLower(part)
		valid := false
		for _, validValue := range mev.validValues {
			if normalized == validValue {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("invalid value '%s', valid options: %s",
				part, strings.Join(mev.validValues, ", "))
		}

		// Check for duplicates
		duplicate := false
		for _, existing := range newValues {
			if strings.ToLower(existing) == normalized {
				duplicate = true
				break
			}
		}

		if !duplicate {
			newValues = append(newValues, part)
		}
	}

	mev.values = newValues
	return nil
}

func (mev *MultiEnumValue) Type() string {
	return fmt.Sprintf("multiEnum[%s]", strings.Join(mev.validValues, "|"))
}

func (mev *MultiEnumValue) Get() interface{} {
	return mev.values
}

func (mev *MultiEnumValue) GetValues() []string {
	return append([]string(nil), mev.values...) // return copy
}

func (mev *MultiEnumValue) Contains(value string) bool {
	normalized := strings.ToLower(value)
	for _, v := range mev.values {
		if strings.ToLower(v) == normalized {
			return true
		}
	}
	return false
}
