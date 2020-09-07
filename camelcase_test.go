package cfg_test // import "github.com/davidwalter0/go-cfg"

import (
	"fmt"
	"strings"
	"testing"

	"github.com/davidwalter0/go-cfg"
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

var camelTests = camelCaseTestSet{
	camelCaseWantGot{"Abc", "ABC", "Abc"},
	camelCaseWantGot{"AbcDef", "AbcDef", "AbcDef"},
	camelCaseWantGot{"abc-def", "AbcDef", "AbcDef"},
	camelCaseWantGot{"ABC-DEF", "AbcDef", "AbcDef"},
	camelCaseWantGot{"abc_def", "AbcDef", "AbcDef"},
	camelCaseWantGot{"ABC_DEF", "AbcDef", "AbcDef"},
	camelCaseWantGot{"AbcDefGhi", "AbcDef Ghi", "AbcDefGhi"},
	camelCaseWantGot{"AbcDEFGhi", "AbcDefGhi", "AbcDefGhi"},
	camelCaseWantGot{"QAMCoral", "QamCoral", "QamCoral"},
	camelCaseWantGot{"QAMCoralQAMCoral", "QamCoralQamCoral", "QamCoralQamCoral"},
	camelCaseWantGot{"QAM-Coral_QAM_Coral", "QamCoralQamCoral", "QamCoralQamCoral"},
	camelCaseWantGot{"QAM_Coral_QAM-Coral", "QamCoralQamCoral", "QamCoralQamCoral"},
	camelCaseWantGot{"QAM-Coral-QAM_Coral", "QamCoralQamCoral", "QamCoralQamCoral"},
	camelCaseWantGot{"QAMCoral", "QamCoral", "QamCoral"},
	camelCaseWantGot{"AbcDEFGhi", "AbcDefGhi", "AbcDefGhi"},
	camelCaseWantGot{"AbcD_EFGhi", "AbcDEfGhi", "AbcDEfGhi"},
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
		var tag = &cfg.Field{KeyName: field.have}
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
		var tag = &cfg.Field{FlagName: field.have}
		tag.FlagNameFromCamelCase()
		t.Logf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		if tag.FlagName != field.wantFlag {
			t.Errorf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		}
	}
}

func TestToLowerSnakeCase(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range tests {
		var tag = &cfg.Field{KeyName: field.have}
		tag.KeyName = cfg.ToLowerSnakeCase(tag.KeyName)
		t.Logf(camelCaseFmt, field.have, strings.ToLower(field.wantKey), tag.KeyName, tag.KeyName == strings.ToLower(field.wantKey))
		if tag.KeyName != strings.ToLower(field.wantKey) {
			t.Errorf(camelCaseFmt, field.have, strings.ToLower(field.wantKey), tag.KeyName, tag.KeyName == strings.ToLower(field.wantKey))
		}
	}
}

func TestToUpperSnakeCase(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range tests {
		var tag = &cfg.Field{KeyName: field.have}
		tag.KeyName = cfg.ToUpperSnakeCase(tag.KeyName)
		t.Logf(camelCaseFmt, field.have, strings.ToUpper(field.wantKey), tag.KeyName, tag.KeyName == strings.ToUpper(field.wantKey))
		if tag.KeyName != strings.ToUpper(field.wantKey) {
			t.Errorf(camelCaseFmt, field.have, strings.ToUpper(field.wantKey), tag.KeyName, tag.KeyName == strings.ToUpper(field.wantKey))
		}
	}
}

func TestToCamelCase(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range camelTests {
		var tag = &cfg.Field{FlagName: cfg.ToLowerSnakeCase(field.have)}
		tag.FlagName = cfg.ToCamelCase(tag.FlagName)
		t.Logf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		if tag.FlagName != field.wantFlag {
			t.Errorf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		}
	}
}
func TestToLowerKebabCase(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range tests {
		var tag = &cfg.Field{FlagName: field.have}
		tag.FlagName = cfg.ToLowerKebabCase(tag.FlagName)
		t.Logf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		if tag.FlagName != field.wantFlag {
			t.Errorf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		}
	}
}

func TestToUpperKebabCase(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range tests {
		var tag = &cfg.Field{FlagName: field.have}
		tag.FlagName = cfg.ToUpperKebabCase(tag.FlagName)
		t.Logf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == strings.ToUpper(field.wantFlag))
		if tag.FlagName != strings.ToUpper(field.wantFlag) {
			t.Errorf(camelCaseFmt, field.have, strings.ToUpper(field.wantFlag), tag.FlagName, tag.FlagName == strings.ToUpper(field.wantFlag))
		}
	}
}

func TestField_KeyNameFromCamelCase(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range tests {
		var tag = &cfg.Field{KeyName: field.have}
		tag.KeyNameFromCamelCase()
		t.Logf(camelCaseFmt, field.have, field.wantKey, tag.KeyName, tag.KeyName == field.wantKey)
		if tag.KeyName != field.wantKey {
			t.Errorf(camelCaseFmt, field.have, field.wantKey, tag.KeyName, tag.KeyName == field.wantKey)
		}
	}
}

func TestField_FlagNameFromCamelCase(t *testing.T) {
	t.Log(camelCaseHeadline)
	t.Log(camelCaseSeparator)
	for _, field := range tests {
		var tag = &cfg.Field{FlagName: field.have}
		tag.FlagNameFromCamelCase()
		t.Logf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		if tag.FlagName != field.wantFlag {
			t.Errorf(camelCaseFmt, field.have, field.wantFlag, tag.FlagName, tag.FlagName == field.wantFlag)
		}
	}
}

func TestCapitalize(t *testing.T) {
	want := []string{"Abc", "Default", "ABC", "DEFAULT", "ABc", "DefAult"}
	have := []string{"abc", "default", "aBC", "dEFAULT", "aBc", "defAult"}
	for i := range want {
		want, got, have := want[i], cfg.Capitalize(have[i]), have[i]
		t.Logf("have %-20s want %-20s got %-20s\n", have, want, got)
		if want != got {
			t.Errorf("have %-20s want %-20s got %-20s\n", have, want, got)
		}
	}
}

func TestDowncase(t *testing.T) {
	want := []string{"abc", "default", "abc", "default", "abc", "default"}
	have := []string{"abc", "default", "aBC", "dEFAULT", "aBc", "defAult"}
	for i := range want {
		want, got, have := want[i], cfg.Downcase(have[i]), have[i]
		t.Logf("have %-20s want %-20s got %-20s\n", have, want, got)
		if want != got {
			t.Errorf("have %-20s want %-20s got %-20s\n", have, want, got)
		}
	}
}
