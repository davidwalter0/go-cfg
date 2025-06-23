package types

import (
	"fmt"
	"net/url"
)

// URLValue implements pflag.Value for URL parsing with validation
type URLValue struct {
	URL *url.URL
}

func NewURLValue(defaultURL string) *URLValue {
	uv := &URLValue{}
	if defaultURL != "" {
		_ = uv.Set(defaultURL)
	}
	return uv
}

func (uv *URLValue) String() string {
	if uv.URL == nil {
		return ""
	}
	return uv.URL.String()
}

func (uv *URLValue) Set(s string) error {
	if s == "" {
		uv.URL = nil
		return nil
	}

	parsed, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid URL '%s': %w", s, err)
	}

	// Validation: ensure scheme is provided for absolute URLs
	if parsed.Scheme == "" && parsed.Host != "" {
		return fmt.Errorf("URL must include scheme (http/https): %s", s)
	}

	// Validation: ensure host for absolute URLs
	if parsed.Scheme != "" && parsed.Host == "" {
		return fmt.Errorf("URL with scheme must include host: %s", s)
	}

	uv.URL = parsed
	return nil
}

func (uv *URLValue) Type() string {
	return "url"
}

func (uv *URLValue) Get() interface{} {
	return uv.URL
}

// GetURL returns the parsed URL or nil
func (uv *URLValue) GetURL() *url.URL {
	return uv.URL
}

// IsSet returns true if URL has been set
func (uv *URLValue) IsSet() bool {
	return uv.URL != nil
}

// IsAbsolute returns true if URL is absolute (has scheme)
func (uv *URLValue) IsAbsolute() bool {
	return uv.URL != nil && uv.URL.Scheme != ""
}

// URLSliceValue implements pflag.Value for multiple URLs
type URLSliceValue struct {
	URLs []*url.URL
}

func NewURLSliceValue(defaultURLs ...string) *URLSliceValue {
	usv := &URLSliceValue{}
	for _, u := range defaultURLs {
		if u != "" {
			_ = usv.Set(u)
		}
	}
	return usv
}

func (usv *URLSliceValue) String() string {
	if len(usv.URLs) == 0 {
		return ""
	}

	var urls []string
	for _, u := range usv.URLs {
		if u != nil {
			urls = append(urls, u.String())
		}
	}
	return fmt.Sprintf("[%s]", fmt.Sprintf("%q", urls)[1:len(fmt.Sprintf("%q", urls))-1])
}

func (usv *URLSliceValue) Set(s string) error {
	uv := &URLValue{}
	if err := uv.Set(s); err != nil {
		return err
	}

	if uv.URL != nil {
		usv.URLs = append(usv.URLs, uv.URL)
	}
	return nil
}

func (usv *URLSliceValue) Type() string {
	return "urlSlice"
}

func (usv *URLSliceValue) Get() interface{} {
	return usv.URLs
}

// GetURLs returns the slice of parsed URLs
func (usv *URLSliceValue) GetURLs() []*url.URL {
	return usv.URLs
}

// Len returns the number of URLs in the slice
func (usv *URLSliceValue) Len() int {
	return len(usv.URLs)
}
