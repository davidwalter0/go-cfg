package legacy

import (
	"testing"

	"github.com/davidwalter0/go-cfg/v3/pkg/config"
)

type TestConfig struct {
	ServiceName string `json:"service_name" default:"test-service"`
	Version     string `json:"version" default:"1.0.0"`
	API         struct {
		Host string `json:"host" default:"localhost"`
		Port int    `json:"port" default:"8080"`
	} `json:"api"`
	Features map[string]bool   `json:"features"`
	Services map[string]string `json:"services"`
}

func TestLegacyInit(t *testing.T) {
	var cfg TestConfig
	
	err := Init(&cfg)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}

	if !Frozen() {
		t.Error("Expected configuration to be frozen after Init()")
	}
}

func TestLegacyEval(t *testing.T) {
	var cfg TestConfig
	
	err := Eval(&cfg)
	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestLegacySimple(t *testing.T) {
	var cfg TestConfig
	
	err := Simple(&cfg)
	if err != nil {
		t.Fatalf("Simple failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestLegacySimpleFlags(t *testing.T) {
	var cfg TestConfig
	
	err := SimpleFlags(&cfg)
	if err != nil {
		t.Fatalf("SimpleFlags failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}

	if !Frozen() {
		t.Error("Expected configuration to be frozen after SimpleFlags()")
	}
}

func TestLegacyNest(t *testing.T) {
	var cfg TestConfig
	
	err := Nest(&cfg)
	if err != nil {
		t.Fatalf("Nest failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestLegacyWrap(t *testing.T) {
	var cfg TestConfig
	
	err := Wrap("PREFIX", &cfg)
	if err != nil {
		t.Fatalf("Wrap failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestLegacyUnwrap(t *testing.T) {
	var cfg TestConfig
	
	err := Unwrap(&cfg)
	if err != nil {
		t.Fatalf("Unwrap failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestLegacyBare(t *testing.T) {
	var cfg TestConfig
	
	err := Bare(&cfg)
	if err != nil {
		t.Fatalf("Bare failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestLegacyAdd(t *testing.T) {
	var cfg TestConfig
	
	err := Add(&cfg)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestLegacyFlags(t *testing.T) {
	var cfg TestConfig
	
	err := Flags(&cfg)
	if err != nil {
		t.Fatalf("Flags failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}

	if !Frozen() {
		t.Error("Expected configuration to be frozen after Flags()")
	}
}

func TestFreezeThaw(t *testing.T) {
	// Test freeze/thaw functionality
	Thaw()
	if Frozen() {
		t.Error("Expected configuration to be thawed")
	}

	Freeze()
	if !Frozen() {
		t.Error("Expected configuration to be frozen")
	}

	Final()
	if !Frozen() {
		t.Error("Expected configuration to be frozen after Final()")
	}
}

func TestTestMode(t *testing.T) {
	// Test test mode functionality
	ProductionMode()
	if InTestMode() {
		t.Error("Expected production mode")
	}

	TestMode()
	if !InTestMode() {
		t.Error("Expected test mode")
	}

	if Ok() {
		t.Error("Expected Ok() to return false in test mode")
	}
}

func TestUtilityFunctions(t *testing.T) {
	var cfg TestConfig

	// Test type checking functions
	if !IsStructPtr(&cfg) {
		t.Error("Expected IsStructPtr to return true for struct pointer")
	}

	if IsStructPtr(cfg) {
		t.Error("Expected IsStructPtr to return false for struct value")
	}

	testMap := make(map[string]string)
	if !IsMap(testMap) {
		t.Error("Expected IsMap to return true for map")
	}

	// Test CheckArgs
	err := CheckArgs(&cfg)
	if err != nil {
		t.Errorf("CheckArgs failed: %v", err)
	}

	err = CheckArgs(cfg) // Not a pointer
	if err == nil {
		t.Error("Expected CheckArgs to fail for non-pointer")
	}
}

func TestEnhancedFlatNest(t *testing.T) {
	var cfg TestConfig

	// Test enhanced FlatNest function
	err := FlatNest(&cfg)
	if err != nil {
		t.Fatalf("FlatNest failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}

	if !Frozen() {
		t.Error("Expected configuration to be frozen after FlatNest()")
	}
}

func TestFlatNestWithDiscovery(t *testing.T) {
	var cfg TestConfig

	// Test FlatNest with discovery modes
	err := FlatNestWithDiscovery(config.Union|config.Direct, &cfg)
	if err != nil {
		t.Fatalf("FlatNestWithDiscovery failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestQuickAPI(t *testing.T) {
	var cfg TestConfig

	// Test Quick API (simple one-liner)
	err := Quick(&cfg)
	if err != nil {
		t.Fatalf("Quick failed: %v", err)
	}

	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected ServiceName 'test-service', got '%s'", cfg.ServiceName)
	}
}

func TestV3Enhancements(t *testing.T) {
	// Test v3-specific enhancements available through legacy API

	// Test help generation
	help := GenerateHelp()
	if help == "" {
		t.Error("Expected GenerateHelp to return non-empty string")
	}

	// Test sample generation
	sample, err := GenerateSample("json")
	if err != nil {
		t.Fatalf("GenerateSample failed: %v", err)
	}
	if len(sample) == 0 {
		t.Error("Expected GenerateSample to return non-empty data")
	}

	// Test config saving
	var cfg TestConfig
	cfg.ServiceName = "saved-service"
	cfg.Version = "2.0.0"

	tempFile := t.TempDir() + "/test-config.json"
	err = SaveCurrent(&cfg, tempFile)
	if err != nil {
		t.Fatalf("SaveCurrent failed: %v", err)
	}
}