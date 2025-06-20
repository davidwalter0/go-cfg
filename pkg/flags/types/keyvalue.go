package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// KeyValueValue implements pflag.Value for key=value pairs
type KeyValueValue struct {
	pairs map[string]string
}

func NewKeyValueValue() *KeyValueValue {
	return &KeyValueValue{
		pairs: make(map[string]string),
	}
}

func (kvv *KeyValueValue) String() string {
	if len(kvv.pairs) == 0 {
		return ""
	}
	
	var pairs []string
	for k, v := range kvv.pairs {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(pairs, ",")
}

func (kvv *KeyValueValue) Set(s string) error {
	if s == "" {
		return nil
	}
	
	// Split by comma for multiple pairs
	parts := strings.Split(s, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		// Split by = for key-value
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid key=value format: %s (expected key=value)", part)
		}
		
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		
		if key == "" {
			return fmt.Errorf("empty key in key=value pair: %s", part)
		}
		
		kvv.pairs[key] = value
	}
	
	return nil
}

func (kvv *KeyValueValue) Type() string {
	return "key=value"
}

func (kvv *KeyValueValue) Get() interface{} {
	// Return a copy to prevent external modification
	result := make(map[string]string)
	for k, v := range kvv.pairs {
		result[k] = v
	}
	return result
}

func (kvv *KeyValueValue) GetMap() map[string]string {
	result := make(map[string]string)
	for k, v := range kvv.pairs {
		result[k] = v
	}
	return result
}

func (kvv *KeyValueValue) GetValue(key string) (string, bool) {
	value, exists := kvv.pairs[key]
	return value, exists
}

func (kvv *KeyValueValue) SetValue(key, value string) {
	if kvv.pairs == nil {
		kvv.pairs = make(map[string]string)
	}
	kvv.pairs[key] = value
}

func (kvv *KeyValueValue) DeleteValue(key string) {
	delete(kvv.pairs, key)
}

func (kvv *KeyValueValue) Keys() []string {
	var keys []string
	for k := range kvv.pairs {
		keys = append(keys, k)
	}
	return keys
}

func (kvv *KeyValueValue) Len() int {
	return len(kvv.pairs)
}

// JSON marshaling
func (kvv KeyValueValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(kvv.pairs)
}

func (kvv *KeyValueValue) UnmarshalJSON(data []byte) error {
	var pairs map[string]string
	if err := json.Unmarshal(data, &pairs); err != nil {
		return err
	}
	
	kvv.pairs = pairs
	return nil
}

// YAML marshaling
func (kvv KeyValueValue) MarshalYAML() (interface{}, error) {
	return kvv.pairs, nil
}

func (kvv *KeyValueValue) UnmarshalYAML(node *yaml.Node) error {
	var pairs map[string]string
	if err := node.Decode(&pairs); err != nil {
		return err
	}
	
	kvv.pairs = pairs
	return nil
}

// TypedKeyValueValue implements pflag.Value for typed key=value pairs
type TypedKeyValueValue[T any] struct {
	pairs  map[string]T
	parser func(string) (T, error)
}

func NewTypedKeyValueValue[T any](parser func(string) (T, error)) *TypedKeyValueValue[T] {
	return &TypedKeyValueValue[T]{
		pairs:  make(map[string]T),
		parser: parser,
	}
}

func (tkvv *TypedKeyValueValue[T]) String() string {
	if len(tkvv.pairs) == 0 {
		return ""
	}
	
	var pairs []string
	for k, v := range tkvv.pairs {
		pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(pairs, ",")
}

func (tkvv *TypedKeyValueValue[T]) Set(s string) error {
	if s == "" {
		return nil
	}
	
	parts := strings.Split(s, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid key=value format: %s", part)
		}
		
		key := strings.TrimSpace(kv[0])
		valueStr := strings.TrimSpace(kv[1])
		
		if key == "" {
			return fmt.Errorf("empty key in key=value pair: %s", part)
		}
		
		value, err := tkvv.parser(valueStr)
		if err != nil {
			return fmt.Errorf("invalid value for key '%s': %w", key, err)
		}
		
		tkvv.pairs[key] = value
	}
	
	return nil
}

func (tkvv *TypedKeyValueValue[T]) Type() string {
	return "key=typedValue"
}

func (tkvv *TypedKeyValueValue[T]) Get() interface{} {
	result := make(map[string]T)
	for k, v := range tkvv.pairs {
		result[k] = v
	}
	return result
}

func (tkvv *TypedKeyValueValue[T]) GetMap() map[string]T {
	result := make(map[string]T)
	for k, v := range tkvv.pairs {
		result[k] = v
	}
	return result
}

func (tkvv *TypedKeyValueValue[T]) GetValue(key string) (T, bool) {
	value, exists := tkvv.pairs[key]
	return value, exists
}

// LabelsValue is a specialized key-value type for Kubernetes-style labels
type LabelsValue struct {
	*KeyValueValue
}

func NewLabelsValue() *LabelsValue {
	return &LabelsValue{
		KeyValueValue: NewKeyValueValue(),
	}
}

func (lv *LabelsValue) Set(s string) error {
	if s == "" {
		return nil
	}
	
	// Validate label format for Kubernetes compatibility
	parts := strings.Split(s, ",")
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid label format: %s (expected key=value)", part)
		}
		
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		
		// Validate key format (basic validation)
		if err := lv.validateLabelKey(key); err != nil {
			return fmt.Errorf("invalid label key '%s': %w", key, err)
		}
		
		// Validate value format
		if err := lv.validateLabelValue(value); err != nil {
			return fmt.Errorf("invalid label value '%s': %w", value, err)
		}
		
		lv.pairs[key] = value
	}
	
	return nil
}

func (lv *LabelsValue) validateLabelKey(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	
	if len(key) > 63 {
		return fmt.Errorf("key too long (max 63 characters)")
	}
	
	// Basic validation - could be more strict for full Kubernetes compliance
	if strings.Contains(key, " ") {
		return fmt.Errorf("key cannot contain spaces")
	}
	
	return nil
}

func (lv *LabelsValue) validateLabelValue(value string) error {
	if len(value) > 63 {
		return fmt.Errorf("value too long (max 63 characters)")
	}
	
	return nil
}

func (lv *LabelsValue) Type() string {
	return "labels"
}