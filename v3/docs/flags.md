package flags // import "github.com/davidwalter0/go-cfg/v3/pkg/flags"


TYPES

type FieldInfo struct {
	Name         string
	Type         string
	JSONTag      string
	DefaultValue string
	HelpText     string
	Required     bool
	FieldPath    []string
}
    FieldInfo contains metadata about a struct field

type FlatNestStrategy struct{}
    FlatNestStrategy implements clean flat naming (--api-host)

func (s *FlatNestStrategy) EnvVarName(fieldPath []string) string
    EnvVarName generates environment variable names

func (s *FlatNestStrategy) FlagName(fieldPath []string) string
    FlagName generates clean flat flag names

func (s *FlatNestStrategy) HelpText(field FieldInfo) string
    HelpText generates help text for a field

type NamingStrategy interface {
	FlagName(fieldPath []string) string
	EnvVarName(fieldPath []string) string
	HelpText(field FieldInfo) string
}
    NamingStrategy defines how struct fields map to flag names

func NewFlatNestStrategy() NamingStrategy
    NewFlatNestStrategy creates a new flat nest naming strategy

func NewNestedStrategy() NamingStrategy
    NewNestedStrategy creates a new nested naming strategy

func NewPrefixedStrategy() NamingStrategy
    NewPrefixedStrategy creates a new prefixed naming strategy

type NestedStrategy struct{}
    NestedStrategy implements hierarchical naming (--api.host)

func (s *NestedStrategy) EnvVarName(fieldPath []string) string
    EnvVarName generates environment variable names

func (s *NestedStrategy) FlagName(fieldPath []string) string
    FlagName generates nested flag names

func (s *NestedStrategy) HelpText(field FieldInfo) string
    HelpText generates help text for a field

type PrefixedStrategy struct {
	// Has unexported fields.
}
    PrefixedStrategy implements prefixed naming (--config-api-host)

func (s *PrefixedStrategy) EnvVarName(fieldPath []string) string
    EnvVarName generates environment variable names

func (s *PrefixedStrategy) FlagName(fieldPath []string) string
    FlagName generates prefixed flag names

func (s *PrefixedStrategy) HelpText(field FieldInfo) string
    HelpText generates help text for a field

type Processor struct {
	// Has unexported fields.
}
    Processor implements FlagProcessor interface

func NewProcessor() *Processor
    NewProcessor creates a new flag processor

func (p *Processor) ApplyFlagValues(config interface{}, flagSet *pflag.FlagSet) error
    ApplyFlagValues applies flag values back to the config struct

func (p *Processor) GenerateFlags(config interface{}) (*pflag.FlagSet, error)
    GenerateFlags generates a flag set from a struct

func (p *Processor) ProcessStruct(config interface{}, flagSet *pflag.FlagSet) error
    ProcessStruct processes a struct and generates flags

