package cfg // import "github.com/davidwalter0/go-cfg"

import (
	"fmt"
	"testing"
)

type camelCaseWantGot struct {
	have     string
	wantKey  string
	wantFlag string
}

type camelCaseTestSet []camelCaseWantGot

var tests = camelCaseTestSet{
	camelCaseWantGot{"Abc", "ABC", "abc"},
	camelCaseWantGot{"AbcDef", "ABC_DEF", "abc-def"},
	camelCaseWantGot{"AbcDefGhi", "ABC_DEF_GHI", "abc-def-ghi"},
	camelCaseWantGot{"AbcDEFGhi", "ABC_DEF_GHI", "abc-def-ghi"},
	camelCaseWantGot{"QAMCoral", "QAM_CORAL", "qam-coral"},
	camelCaseWantGot{"QAMCoralQAMCoral", "QAM_CORAL_QAM_CORAL", "qam-coral-qam-coral"},
	camelCaseWantGot{"QAMCoral", "QAM_CORAL", "qam-coral"},
	camelCaseWantGot{"AbcDEFGhi", "ABC_DEF_GHI", "abc-def-ghi"},
	camelCaseWantGot{"AbcD_EFGhi", "ABC_D_EF_GHI", "abc-d-ef-ghi"},
}

var camelCaseFmt = "%20s %20s %20s %5v\n"
var camelCaseHeadline = fmt.Sprintf("%20s %20s %20s %5v", "have", "want", "got", "bool")
var dash = "--------------------"
var camelCaseSeparator = fmt.Sprintf("%20.20s %20.20s %20.20s %5.5v", dash, dash, dash, dash)

// TestUnderScoreCamelCaseWords split on CamelCase words
func TestUnderScoreCamelCaseWords(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range tests {
		var tag = &Field{KeyName: field.have}
		tag.KeyNameFromCamelCase()
		t.Logf(camelCaseFmt, field.have, field.wantKey, tag.KeyName, tag.KeyName == field.wantKey)
		if tag.KeyName != field.wantKey {
			t.Errorf(camelCaseFmt, field.have, field.wantKey, tag.KeyName, tag.KeyName == field.wantKey)
		}
	}
}

// TestHyphenateCamelCaseWords converts camel case name string and
// hyphenates words for flags between words
func TestHyphenateCamelCaseWords(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range tests {
		var tag = &Field{FlagName: field.have}
		tag.FlagNameFromCamelCase()
		t.Logf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		if tag.FlagName != field.wantFlag {
			t.Errorf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		}
	}
}
