package cfg // import "github.com/davidwalter0/go-cfg"

import (
	"testing"
)

// TestUnderScoreCamelCaseWords split on CamelCase words
func TestUnderScoreCamelCaseWords(t *testing.T) {
	var tag = &Field{}
	tag.UnderScoreCamelCaseWords()

}

// TestHyphenateCamelCaseWords converts camel case name string and
// hyphenates words for flags between words
func TestHyphenateCamelCaseWords(t *testing.T) {
	var tag = &Field{}
	tag.HyphenateCamelCaseWords()
}
