package discovery // import "github.com/davidwalter0/go-cfg/v3/pkg/discovery"


TYPES

type DiscoveryMode int
    DiscoveryMode defines how configuration files are discovered and merged

const (
	Union    DiscoveryMode = 1 << iota // Merge all found configs
	First                              // Stop on first found
	Direct                             // Load config files directly
	Indirect                           // Support AutoCfg pointer files
)
type DiscoveryOptions struct {
	Mode        DiscoveryMode
	SearchPaths []string
	EnvOverride string
	ProgramName string
}
    DiscoveryOptions provides configuration for file discovery

type Engine struct{}
    Engine implements DiscoveryEngine interface

func NewEngine() *Engine
    NewEngine creates a new discovery engine

func (e *Engine) FindConfigFiles(opts DiscoveryOptions) ([]string, error)
    FindConfigFiles discovers configuration files based on options

func (e *Engine) LoadConfigFile(filename string, config interface{}) error
    LoadConfigFile loads a single configuration file

func (e *Engine) MergeConfigs(files []string, config interface{}) error
    MergeConfigs loads and merges multiple configuration files

