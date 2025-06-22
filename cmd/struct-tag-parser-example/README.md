# Struct Tag Parser Example

This example demonstrates how to use `go:generate` to extract comments from Go struct fields and use them with go-cfg's flag system for automatic help text generation.

## How It Works

1. **Comment Extraction**: The `go:generate` directive runs the struct-tag-parser to extract comments from struct fields
2. **Tag Generation**: Comments are converted to struct tags (usage, default, env, etc.)
3. **Flag Registration**: go-cfg's `RegisterCustomTypes` uses the generated tags
4. **Help Display**: `fc.Usage()` automatically displays help text from the extracted comments

## Files

- `types.go` - Nested struct definitions with descriptive comments
- `main.go` - Main application demonstrating the integration
- `go.mod` - Module dependencies

## Usage

### Step 1: Generate Tags from Comments

```bash
# Run go generate to extract comments and create struct tags
go generate
```

This will process the comment extraction and modify the source files to include proper struct tags.

### Step 2: Run the Application

```bash
# Show help (displays comment-derived usage text)
go run . --help

# Run with some configuration
go run . --service-name "my-api" --debug --worker-count 20

# Use nested struct flags
go run . --database-host "prod-db.com" --database-port 5432 --server-port 8080

# Use environment variables (automatically mapped)
export SERVICE_NAME="my-api"
export DATABASE_HOST="prod-db.company.com" 
export MONITORING_LOG_LEVEL="warn"
go run .
```

## Comment Format

The struct-tag-parser recognizes these comment patterns:

```go
type Config struct {
    // Field description, default: value, env: CUSTOM_ENV_NAME
    FieldName string
    
    // Examples:
    ServiceName string // Service identifier, default: my-service
    Port        int    // Server port, default: 8080
    APIKey      string // API authentication key, env: CUSTOM_API_KEY
    Debug       bool   // Enable debug mode
}
```

## Generated Tags

After running `go generate`, the comments become struct tags:

```go
type Config struct {
    ServiceName string `json:"service-name" yaml:"service-name" doc:"Service identifier" default:"my-service" env:"SERVICE_NAME" usage:"Service identifier"`
    Port        int    `json:"port" yaml:"port" doc:"Server port" default:"8080" env:"PORT" usage:"Server port"`
    APIKey      string `json:"api-key" yaml:"api-key" doc:"API authentication key" env:"CUSTOM_API_KEY" usage:"API authentication key"`
    Debug       bool   `json:"debug" yaml:"debug" doc:"Enable debug mode" env:"DEBUG" usage:"Enable debug mode"`
}
```

## Benefits

1. **Single Source of Truth**: Write documentation once in comments
2. **Automatic Help Text**: go-cfg generates help from comment-derived usage tags
3. **Consistent Naming**: Automatic kebab-case flags, SCREAMING_SNAKE_CASE env vars
4. **Nested Struct Support**: Works with complex nested configurations
5. **IDE Friendly**: Natural comments work with Go tooling

## Integration with go-cfg

The generated `usage` tags are automatically used by:

- `fc.Usage()` - Displays comprehensive help text
- Flag descriptions in `--help` output
- Error messages when validation fails
- Documentation generation tools

This approach combines the convenience of comment-based documentation with the power of go-cfg's type-safe configuration system.