package main

import (
	"fmt"
	"net"

	"github.com/davidwalter0/go-cfg/pkg/flags/types"
)

func main() {
	fmt.Println("=== Testing Custom Flag Types ===\n")
	
	// Test URL Value
	fmt.Println("1. Testing URLValue:")
	urlVal := types.NewURLValue("")
	testURL := "https://api.example.com:8080/v1/users?limit=10"
	if err := urlVal.Set(testURL); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Input: %s\n", testURL)
		fmt.Printf("  String(): %s\n", urlVal.String())
		fmt.Printf("  IsSet(): %t\n", urlVal.IsSet())
		fmt.Printf("  IsAbsolute(): %t\n", urlVal.IsAbsolute())
		if u := urlVal.GetURL(); u != nil {
			fmt.Printf("  Host: %s, Path: %s\n", u.Host, u.Path)
		}
	}
	
	// Test invalid URL
	fmt.Println("\n  Testing invalid URL:")
	if err := urlVal.Set("not-a-valid-url"); err != nil {
		fmt.Printf("  Expected error: %v\n", err)
	}
	
	// Test Duration Value
	fmt.Println("\n2. Testing DurationValue:")
	durVal := types.NewDurationValue(0)
	testDurations := []string{"30s", "1m30s", "90", "90.5"}
	for _, dur := range testDurations {
		if err := durVal.Set(dur); err != nil {
			fmt.Printf("  %s -> Error: %v\n", dur, err)
		} else {
			fmt.Printf("  %s -> %v\n", dur, durVal.GetDuration())
		}
	}
	
	// Test Human Duration Value
	fmt.Println("\n3. Testing HumanDurationValue:")
	humanDur := types.NewHumanDurationValue(0)
	humanFormats := []string{"1 hour", "30 minutes", "2 hours 15 minutes", "1 day", "1w"}
	for _, dur := range humanFormats {
		if err := humanDur.Set(dur); err != nil {
			fmt.Printf("  %s -> Error: %v\n", dur, err)
		} else {
			fmt.Printf("  %s -> %v (displayed as: %s)\n", 
				dur, humanDur.GetDuration(), humanDur.String())
		}
	}
	
	// Test Enum Value
	fmt.Println("\n4. Testing EnumValue (LogLevel):")
	logLevel := types.NewLogLevelValue("info")
	testLevels := []string{"debug", "DEBUG", "warn", "invalid"}
	for _, level := range testLevels {
		if err := logLevel.Set(level); err != nil {
			fmt.Printf("  %s -> Error: %v\n", level, err)
		} else {
			fmt.Printf("  %s -> %s (debug enabled: %t)\n", 
				level, logLevel.GetLevel(), logLevel.IsDebugEnabled())
		}
	}
	
	// Test Environment Value
	fmt.Println("\n5. Testing EnvironmentValue:")
	env := types.NewEnvironmentValue("development")
	testEnvs := []string{"prod", "production", "staging", "dev", "invalid"}
	for _, e := range testEnvs {
		if err := env.Set(e); err != nil {
			fmt.Printf("  %s -> Error: %v\n", e, err)
		} else {
			fmt.Printf("  %s -> %s (production: %t, development: %t)\n", 
				e, env.GetEnvironment(), env.IsProduction(), env.IsDevelopment())
		}
	}
	
	// Test KeyValue
	fmt.Println("\n6. Testing KeyValueValue:")
	kv := types.NewKeyValueValue()
	testKV := "service=api,version=v1.2.3,env=production"
	if err := kv.Set(testKV); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Input: %s\n", testKV)
		fmt.Printf("  String(): %s\n", kv.String())
		fmt.Printf("  Keys: %v\n", kv.Keys())
		for _, key := range kv.Keys() {
			if value, exists := kv.GetValue(key); exists {
				fmt.Printf("    %s = %s\n", key, value)
			}
		}
	}
	
	// Test HostPort Value
	fmt.Println("\n7. Testing HostPortValue:")
	hostPort := types.NewHostPortValue("", 0)
	testHostPorts := []string{"localhost:8080", "192.168.1.1:3306", "[::1]:9090", "example.com"}
	for _, hp := range testHostPorts {
		if err := hostPort.Set(hp); err != nil {
			fmt.Printf("  %s -> Error: %v\n", hp, err)
		} else {
			fmt.Printf("  %s -> Host: %s, Port: %d, Address: %s\n", 
				hp, hostPort.GetHost(), hostPort.GetPort(), hostPort.GetAddress())
		}
	}
	
	// Test IP Value
	fmt.Println("\n8. Testing IPValue:")
	ipVal := types.NewIPValue("")
	testIPs := []string{"127.0.0.1", "::1", "192.168.1.100", "2001:db8::1", "invalid"}
	for _, ip := range testIPs {
		if err := ipVal.Set(ip); err != nil {
			fmt.Printf("  %s -> Error: %v\n", ip, err)
		} else {
			fmt.Printf("  %s -> IPv4: %t, IPv6: %t\n", 
				ip, ipVal.IsIPv4(), ipVal.IsIPv6())
		}
	}
	
	// Test CIDR Value
	fmt.Println("\n9. Testing CIDRValue:")
	cidr := types.NewCIDRValue("")
	testCIDRs := []string{"192.168.1.0/24", "10.0.0.0/8", "2001:db8::/32", "invalid"}
	for _, c := range testCIDRs {
		if err := cidr.Set(c); err != nil {
			fmt.Printf("  %s -> Error: %v\n", c, err)
		} else {
			fmt.Printf("  %s -> Network: %s\n", c, cidr.String())
			// Test if an IP is in the network
			if cidr.GetIPNet() != nil {
				testIP := net.ParseIP("192.168.1.100")
				if testIP != nil {
					fmt.Printf("    Contains 192.168.1.100: %t\n", cidr.Contains(testIP))
				}
			}
		}
	}
	
	// Test Port Range Value
	fmt.Println("\n10. Testing PortRangeValue:")
	portRange := types.NewPortRangeValue(0, 0)
	testRanges := []string{"8000-9000", "80", "1-65535", "70000-80000"}
	for _, pr := range testRanges {
		if err := portRange.Set(pr); err != nil {
			fmt.Printf("  %s -> Error: %v\n", pr, err)
		} else {
			fmt.Printf("  %s -> Start: %d, End: %d, Size: %d\n", 
				pr, portRange.GetStart(), portRange.GetEnd(), portRange.Size())
			fmt.Printf("    Contains 8080: %t\n", portRange.Contains(8080))
		}
	}
	
	// Test URL Slice
	fmt.Println("\n11. Testing URLSliceValue:")
	urlSlice := types.NewURLSliceValue()
	urls := []string{
		"https://api1.example.com",
		"https://api2.example.com:8080",
		"http://backup.example.com/api",
	}
	for _, u := range urls {
		if err := urlSlice.Set(u); err != nil {
			fmt.Printf("  %s -> Error: %v\n", u, err)
		} else {
			fmt.Printf("  Added: %s\n", u)
		}
	}
	fmt.Printf("  Total URLs: %d\n", urlSlice.Len())
	fmt.Printf("  String representation: %s\n", urlSlice.String())
	
	// Test Labels Value (specialized KeyValue)
	fmt.Println("\n12. Testing LabelsValue:")
	labels := types.NewLabelsValue()
	testLabels := "app=myservice,version=1.0.0,environment=production"
	if err := labels.Set(testLabels); err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Labels: %s\n", labels.String())
	}
	
	// Test naming conventions
	fmt.Println("\n13. Testing NamingConventions:")
	nc := &types.NamingConventions{}
	testNames := []string{"APIEndpoint", "DatabaseURL", "HTTPSProxy", "LogLevel"}
	for _, name := range testNames {
		snake := nc.ToUpperSnakeCase(name)
		kebab := nc.ToKebabCase(name)
		fmt.Printf("  %s -> ENV: %s, Flag: %s\n", name, snake, kebab)
	}
	
	// Test environment variable retrieval
	fmt.Println("\n14. Testing Environment Variable Defaults:")
	testFields := []string{"APIEndpoint", "DatabaseURL", "LogLevel"}
	for _, field := range testFields {
		envVal := types.GetEnvDefault(field)
		if envVal != "" {
			fmt.Printf("  %s: %s\n", field, envVal)
		} else {
			fmt.Printf("  %s: (not set)\n", field)
		}
	}
	
	fmt.Println("\n=== All Tests Completed ===")
}