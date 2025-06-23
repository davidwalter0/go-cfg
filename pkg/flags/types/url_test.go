package types

import (
	"testing"
)

func TestURLValue_SetAndGet(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		wantHost  string
		wantPath  string
	}{
		{
			name:      "valid absolute URL",
			input:     "https://api.example.com/v1/users",
			wantError: false,
			wantHost:  "api.example.com",
			wantPath:  "/v1/users",
		},
		{
			name:      "valid URL with port",
			input:     "http://localhost:8080/api",
			wantError: false,
			wantHost:  "localhost:8080",
			wantPath:  "/api",
		},
		{
			name:      "valid URL with query params",
			input:     "https://api.example.com/search?q=test&limit=10",
			wantError: false,
			wantHost:  "api.example.com",
			wantPath:  "/search",
		},
		{
			name:      "relative URL",
			input:     "/api/v1/users",
			wantError: false,
			wantHost:  "",
			wantPath:  "/api/v1/users",
		},
		{
			name:      "empty string",
			input:     "",
			wantError: false,
			wantHost:  "",
			wantPath:  "",
		},
		{
			name:      "invalid URL",
			input:     "not-a-valid-url",
			wantError: false, // url.Parse is lenient
			wantHost:  "",
			wantPath:  "not-a-valid-url",
		},
		{
			name:      "scheme without host",
			input:     "https://",
			wantError: true,
			wantHost:  "",
			wantPath:  "",
		},
		{
			name:      "host without scheme for absolute URL",
			input:     "example.com/path",
			wantError: false, // This is actually valid as a relative URL
			wantHost:  "",
			wantPath:  "example.com/path",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urlVal := NewURLValue("")
			err := urlVal.Set(test.input)

			if (err != nil) != test.wantError {
				t.Errorf("URLValue.Set(%s) error = %v, wantError = %v", test.input, err, test.wantError)
				return
			}

			if !test.wantError && urlVal.IsSet() {
				url := urlVal.GetURL()
				if url.Host != test.wantHost {
					t.Errorf("URLValue.Set(%s) host = %s, want %s", test.input, url.Host, test.wantHost)
				}
				if url.Path != test.wantPath {
					t.Errorf("URLValue.Set(%s) path = %s, want %s", test.input, url.Path, test.wantPath)
				}
			}
		})
	}
}

func TestURLValue_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute URL",
			input:    "https://api.example.com/v1",
			expected: "https://api.example.com/v1",
		},
		{
			name:     "relative URL",
			input:    "/api/v1",
			expected: "/api/v1",
		},
		{
			name:     "empty URL",
			input:    "",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urlVal := NewURLValue("")
			if test.input != "" {
				_ = urlVal.Set(test.input)
			}

			result := urlVal.String()
			if result != test.expected {
				t.Errorf("URLValue.String() = %s, want %s", result, test.expected)
			}
		})
	}
}

func TestURLValue_IsAbsolute(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "https URL",
			input:    "https://example.com",
			expected: true,
		},
		{
			name:     "http URL",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "relative URL",
			input:    "/path/to/resource",
			expected: false,
		},
		{
			name:     "empty URL",
			input:    "",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urlVal := NewURLValue("")
			if test.input != "" {
				_ = urlVal.Set(test.input)
			}

			result := urlVal.IsAbsolute()
			if result != test.expected {
				t.Errorf("URLValue.IsAbsolute() = %v, want %v", result, test.expected)
			}
		})
	}
}

func TestURLValue_Type(t *testing.T) {
	urlVal := NewURLValue("")
	if urlVal.Type() != "url" {
		t.Errorf("URLValue.Type() = %s, want url", urlVal.Type())
	}
}

func TestURLSliceValue_SetAndGet(t *testing.T) {
	tests := []struct {
		name      string
		inputs    []string
		wantError bool
		wantLen   int
	}{
		{
			name:      "multiple valid URLs",
			inputs:    []string{"https://api1.example.com", "https://api2.example.com", "http://backup.example.com"},
			wantError: false,
			wantLen:   3,
		},
		{
			name:      "mixed valid and relative URLs",
			inputs:    []string{"https://api.example.com", "/api/v1", "http://localhost:8080"},
			wantError: false,
			wantLen:   3,
		},
		{
			name:      "empty strings",
			inputs:    []string{"", "https://api.example.com", ""},
			wantError: false,
			wantLen:   1, // Empty strings should be ignored
		},
		{
			name:      "all empty",
			inputs:    []string{"", "", ""},
			wantError: false,
			wantLen:   0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			urlSlice := NewURLSliceValue()

			var err error
			for _, input := range test.inputs {
				if setErr := urlSlice.Set(input); setErr != nil {
					err = setErr
					break
				}
			}

			if (err != nil) != test.wantError {
				t.Errorf("URLSliceValue.Set() error = %v, wantError = %v", err, test.wantError)
				return
			}

			if !test.wantError {
				if urlSlice.Len() != test.wantLen {
					t.Errorf("URLSliceValue.Len() = %d, want %d", urlSlice.Len(), test.wantLen)
				}

				urls := urlSlice.GetURLs()
				if len(urls) != test.wantLen {
					t.Errorf("URLSliceValue.GetURLs() length = %d, want %d", len(urls), test.wantLen)
				}
			}
		})
	}
}

func TestURLSliceValue_String(t *testing.T) {
	urlSlice := NewURLSliceValue()

	// Test empty slice
	if urlSlice.String() != "" {
		t.Errorf("Empty URLSliceValue.String() = %s, want empty string", urlSlice.String())
	}

	// Add URLs
	_ = urlSlice.Set("https://api1.example.com")
	_ = urlSlice.Set("https://api2.example.com")

	result := urlSlice.String()
	expected := `["https://api1.example.com" "https://api2.example.com"]`
	if result != expected {
		t.Errorf("URLSliceValue.String() = %s, want %s", result, expected)
	}
}

func TestURLSliceValue_Type(t *testing.T) {
	urlSlice := NewURLSliceValue()
	if urlSlice.Type() != "urlSlice" {
		t.Errorf("URLSliceValue.Type() = %s, want urlSlice", urlSlice.Type())
	}
}

// Benchmark tests
func BenchmarkURLValue_Set(b *testing.B) {
	urlVal := NewURLValue("")
	testURL := "https://api.example.com/v1/users?limit=100&offset=0"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = urlVal.Set(testURL)
	}
}

func BenchmarkURLSliceValue_Set(b *testing.B) {
	urlSlice := NewURLSliceValue()
	testURLs := []string{
		"https://api1.example.com",
		"https://api2.example.com",
		"http://backup.example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, url := range testURLs {
			_ = urlSlice.Set(url)
		}
	}
}
