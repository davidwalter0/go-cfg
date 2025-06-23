package generate // import "github.com/davidwalter0/go-cfg/v3/pkg/generate"


TYPES

type DocumentationGenerator struct{}
    DocumentationGenerator creates documentation from struct tags

func NewDocumentationGenerator() *DocumentationGenerator
    NewDocumentationGenerator creates a new documentation generator

func (g *DocumentationGenerator) GenerateMarkdown(config interface{}) string
    GenerateMarkdown creates markdown documentation

type SampleGenerator struct{}
    SampleGenerator creates sample configuration files

func NewSampleGenerator() *SampleGenerator
    NewSampleGenerator creates a new sample generator

func (g *SampleGenerator) GenerateJSON(config interface{}) ([]byte, error)
    GenerateJSON creates a sample JSON configuration

func (g *SampleGenerator) GenerateYAML(config interface{}) ([]byte, error)
    GenerateYAML creates a sample YAML configuration

