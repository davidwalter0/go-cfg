package flags

import (
	"reflect"
	"testing"
	"time"

	"github.com/davidwalter0/go-cfg/pkg/flags/types"
)

// Test configuration with slices and nested structs
type SliceTestConfig struct {
	StringSlice []string `usage:"String slice" default:"a,b,c"`
	IntSlice    []int    `usage:"Int slice" default:"1,2,3"`
	Name        string   `usage:"Application name" default:"slice-test"`
}

type DatabaseConfig struct {
	Host     string        `usage:"Database host" default:"localhost"`
	Port     int           `usage:"Database port" default:"5432"`
	Name     string        `usage:"Database name" default:"mydb"`
	Timeout  time.Duration `usage:"Database timeout" default:"30s"`
	ReadOnly bool          `usage:"Database read-only mode" default:"false"`
}

type ServerConfig struct {
	Host string `usage:"Server host" default:"0.0.0.0"`
	Port int    `usage:"Server port" default:"8080"`
}

type NestedTestConfig struct {
	AppName  string         `usage:"Application name" default:"nested-test"`
	Debug    bool           `usage:"Debug mode" default:"false"`
	Database DatabaseConfig `usage:"Database configuration"`
	Server   ServerConfig   `usage:"Server configuration"`
	Tags     []string       `usage:"Application tags" default:"web,api"`
	Ports    []int          `usage:"Additional ports" default:"8081,8082,8083"`
}

func TestSliceParsing_StringSlice(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		envVars  map[string]string
		expected []string
	}{
		{
			name:     "default values",
			args:     []string{},
			envVars:  map[string]string{},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "flag override",
			args:     []string{"--string-slice", "x,y,z"},
			envVars:  map[string]string{},
			expected: []string{"x", "y", "z"},
		},
		{
			name:     "env var override",
			args:     []string{},
			envVars:  map[string]string{"STRING_SLICE": "env1,env2,env3"},
			expected: []string{"env1", "env2", "env3"},
		},
		{
			name:     "flag overrides env var",
			args:     []string{"--string-slice", "flag1,flag2"},
			envVars:  map[string]string{"STRING_SLICE": "env1,env2"},
			expected: []string{"flag1", "flag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &SliceTestConfig{}

			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Parse configuration
			fc := NewFlagConfig()
			err := fc.RegisterCustomTypes(config)
			if err != nil {
				t.Fatalf("RegisterFlags() error = %v", err)
			}

			err = fc.Parse(tt.args)
			if err != nil {
				t.Fatalf("ParseArgs() error = %v", err)
			}

			if !reflect.DeepEqual(config.StringSlice, tt.expected) {
				t.Errorf("StringSlice = %v, want %v", config.StringSlice, tt.expected)
			}
		})
	}
}

func TestSliceParsing_IntSlice(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		envVars  map[string]string
		expected []int
	}{
		{
			name:     "default values",
			args:     []string{},
			envVars:  map[string]string{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "flag override",
			args:     []string{"--int-slice", "10,20,30"},
			envVars:  map[string]string{},
			expected: []int{10, 20, 30},
		},
		{
			name:     "env var override",
			args:     []string{},
			envVars:  map[string]string{"INT_SLICE": "100,200,300"},
			expected: []int{100, 200, 300},
		},
		{
			name:     "single value",
			args:     []string{"--int-slice", "42"},
			envVars:  map[string]string{},
			expected: []int{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &SliceTestConfig{}

			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Parse configuration
			fc := NewFlagConfig()
			err := fc.RegisterCustomTypes(config)
			if err != nil {
				t.Fatalf("RegisterFlags() error = %v", err)
			}

			err = fc.Parse(tt.args)
			if err != nil {
				t.Fatalf("ParseArgs() error = %v", err)
			}

			if !reflect.DeepEqual(config.IntSlice, tt.expected) {
				t.Errorf("IntSlice = %v, want %v", config.IntSlice, tt.expected)
			}
		})
	}
}

func TestNestedStructParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		envVars     map[string]string
		expectedApp string
		expectedDB  DatabaseConfig
		expectedSrv ServerConfig
		expectedTag []string
	}{
		{
			name:        "default values",
			args:        []string{},
			envVars:     map[string]string{},
			expectedApp: "nested-test",
			expectedDB: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				Name:     "mydb",
				Timeout:  30 * time.Second,
				ReadOnly: false,
			},
			expectedSrv: ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
			},
			expectedTag: []string{"web", "api"},
		},
		{
			name: "flag overrides",
			args: []string{
				"--app-name", "custom-app",
				"--database-host", "db.example.com",
				"--database-port", "3306",
				"--server-port", "9000",
				"--tags", "production,backend",
			},
			envVars:     map[string]string{},
			expectedApp: "custom-app",
			expectedDB: DatabaseConfig{
				Host:     "db.example.com",
				Port:     3306,
				Name:     "mydb",
				Timeout:  30 * time.Second,
				ReadOnly: false,
			},
			expectedSrv: ServerConfig{
				Host: "0.0.0.0",
				Port: 9000,
			},
			expectedTag: []string{"production", "backend"},
		},
		{
			name: "env var overrides",
			args: []string{},
			envVars: map[string]string{
				"APP_NAME":      "env-app",
				"DATABASE_HOST": "env-db.com",
				"DATABASE_PORT": "5433",
				"SERVER_HOST":   "127.0.0.1",
				"SERVER_PORT":   "8090",
				"TAGS":          "env,test",
				"DEBUG":         "true",
			},
			expectedApp: "env-app",
			expectedDB: DatabaseConfig{
				Host:     "env-db.com",
				Port:     5433,
				Name:     "mydb",
				Timeout:  30 * time.Second,
				ReadOnly: false,
			},
			expectedSrv: ServerConfig{
				Host: "127.0.0.1",
				Port: 8090,
			},
			expectedTag: []string{"env", "test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &NestedTestConfig{}

			// Set environment variables
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Parse configuration
			fc := NewFlagConfig()
			err := fc.RegisterCustomTypes(config)
			if err != nil {
				t.Fatalf("RegisterFlags() error = %v", err)
			}

			err = fc.Parse(tt.args)
			if err != nil {
				t.Fatalf("ParseArgs() error = %v", err)
			}

			// Validate results
			if config.AppName != tt.expectedApp {
				t.Errorf("AppName = %s, want %s", config.AppName, tt.expectedApp)
			}

			if config.Database.Host != tt.expectedDB.Host {
				t.Errorf("Database.Host = %s, want %s", config.Database.Host, tt.expectedDB.Host)
			}

			if config.Database.Port != tt.expectedDB.Port {
				t.Errorf("Database.Port = %d, want %d", config.Database.Port, tt.expectedDB.Port)
			}

			if config.Server.Host != tt.expectedSrv.Host {
				t.Errorf("Server.Host = %s, want %s", config.Server.Host, tt.expectedSrv.Host)
			}

			if config.Server.Port != tt.expectedSrv.Port {
				t.Errorf("Server.Port = %d, want %d", config.Server.Port, tt.expectedSrv.Port)
			}

			if !reflect.DeepEqual(config.Tags, tt.expectedTag) {
				t.Errorf("Tags = %v, want %v", config.Tags, tt.expectedTag)
			}
		})
	}
}

func TestSliceAndNestedStruct_Integration(t *testing.T) {
	config := &NestedTestConfig{}

	// Set some environment variables
	t.Setenv("DATABASE_TIMEOUT", "1m30s")
	t.Setenv("PORTS", "9001,9002,9003")

	args := []string{
		"--app-name", "integration-test",
		"--database-read-only", "true",
		"--tags", "integration,test,suite",
	}

	fc := NewFlagConfig()
	err := fc.RegisterCustomTypes(config)
	if err != nil {
		t.Fatalf("RegisterFlags() error = %v", err)
	}

	err = fc.Parse(args)
	if err != nil {
		t.Fatalf("ParseArgs() error = %v", err)
	}

	// Validate all configurations
	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"AppName", config.AppName, "integration-test"},
		{"Database.ReadOnly", config.Database.ReadOnly, true},
		{"Database.Timeout", config.Database.Timeout, 90 * time.Second},
		{"Tags", config.Tags, []string{"integration", "test", "suite"}},
		{"Ports", config.Ports, []int{9001, 9002, 9003}},
		{"Database.Host (default)", config.Database.Host, "localhost"},
		{"Server.Port (default)", config.Server.Port, 8080},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(test.got, test.expected) {
				t.Errorf("%s = %v, want %v", test.name, test.got, test.expected)
			}
		})
	}
}

func TestSliceParsing_EdgeCases(t *testing.T) {
	t.Run("empty_slice_values", func(t *testing.T) {
		config := &SliceTestConfig{}

		fc := NewFlagConfig()
		err := fc.RegisterCustomTypes(config)
		if err != nil {
			t.Fatalf("RegisterFlags() error = %v", err)
		}

		// Test empty slice via flag
		err = fc.Parse([]string{"--string-slice", ""})
		if err != nil {
			t.Fatalf("ParseArgs() error = %v", err)
		}

		// Should result in empty slice, not default
		if len(config.StringSlice) != 1 || config.StringSlice[0] != "" {
			t.Errorf("Empty string slice = %v, want [\"\"]", config.StringSlice)
		}
	})

	t.Run("whitespace_handling", func(t *testing.T) {
		config := &SliceTestConfig{}

		fc := NewFlagConfig()
		err := fc.RegisterCustomTypes(config)
		if err != nil {
			t.Fatalf("RegisterFlags() error = %v", err)
		}

		// Test slice with whitespace
		err = fc.Parse([]string{"--string-slice", " a , b , c "})
		if err != nil {
			t.Fatalf("ParseArgs() error = %v", err)
		}

		expected := []string{" a ", " b ", " c "}
		if !reflect.DeepEqual(config.StringSlice, expected) {
			t.Errorf("Whitespace slice = %v, want %v", config.StringSlice, expected)
		}
	})
}

func TestNestedStruct_FlagNaming(t *testing.T) {
	config := &NestedTestConfig{}

	fc := NewFlagConfig()
	err := fc.RegisterCustomTypes(config)
	if err != nil {
		t.Fatalf("RegisterFlags() error = %v", err)
	}

	// Test that all expected flags are registered
	expectedFlags := []string{
		"app-name",
		"debug",
		"database-host",
		"database-port",
		"database-name",
		"database-timeout",
		"database-read-only",
		"server-host",
		"server-port",
		"tags",
		"ports",
	}

	flagSet := fc.FlagSet
	for _, flagName := range expectedFlags {
		flag := flagSet.Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag %s not found", flagName)
		}
	}
}

func TestCustomTypes_InNestedStructs(t *testing.T) {
	type AdvancedConfig struct {
		APIEndpoint *types.URLValue         `usage:"API endpoint"`
		LogLevel    *types.LogLevelValue    `usage:"Log level"`
		Environment *types.EnvironmentValue `usage:"Environment"`
		Labels      *types.KeyValueValue    `usage:"Labels"`
		ListenAddr  *types.HostPortValue    `usage:"Listen address"`
		Timeouts    []time.Duration         `usage:"Timeout values" default:"1s,5s,30s"`
	}

	config := &AdvancedConfig{}

	fc := NewFlagConfig()
	err := fc.RegisterCustomTypes(config)
	if err != nil {
		t.Fatalf("RegisterFlags() error = %v", err)
	}

	args := []string{
		"--timeouts", "2s,10s,60s",
	}

	err = fc.Parse(args)
	if err != nil {
		t.Fatalf("ParseArgs() error = %v", err)
	}

	expected := []time.Duration{2 * time.Second, 10 * time.Second, 60 * time.Second}
	if !reflect.DeepEqual(config.Timeouts, expected) {
		t.Errorf("Timeouts = %v, want %v", config.Timeouts, expected)
	}
}
