package types

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// HostPortValue implements pflag.Value for host:port combinations
type HostPortValue struct {
	Host string
	Port int
}

func NewHostPortValue(defaultHost string, defaultPort int) *HostPortValue {
	return &HostPortValue{
		Host: defaultHost,
		Port: defaultPort,
	}
}

func (hpv *HostPortValue) String() string {
	if hpv.Host == "" && hpv.Port == 0 {
		return ""
	}
	if hpv.Host == "" {
		return fmt.Sprintf(":%d", hpv.Port)
	}
	return fmt.Sprintf("%s:%d", hpv.Host, hpv.Port)
}

func (hpv *HostPortValue) Set(s string) error {
	if s == "" {
		hpv.Host = ""
		hpv.Port = 0
		return nil
	}

	// Handle IPv6 addresses in brackets
	if strings.HasPrefix(s, "[") {
		closeBracket := strings.Index(s, "]")
		if closeBracket == -1 {
			return fmt.Errorf("invalid IPv6 address format: %s", s)
		}

		hpv.Host = s[1:closeBracket]
		if closeBracket+1 < len(s) {
			if s[closeBracket+1] != ':' {
				return fmt.Errorf("expected ':' after IPv6 address: %s", s)
			}
			portStr := s[closeBracket+2:]
			if portStr != "" {
				port, err := strconv.Atoi(portStr)
				if err != nil {
					return fmt.Errorf("invalid port number: %s", portStr)
				}
				if port < 1 || port > 65535 {
					return fmt.Errorf("port number out of range (1-65535): %d", port)
				}
				hpv.Port = port
			}
		}
		return nil
	}

	// Handle regular host:port format
	parts := strings.SplitN(s, ":", 2)

	if len(parts) == 1 {
		// Only host provided
		hpv.Host = parts[0]
		hpv.Port = 0
	} else {
		// Host and port provided
		hpv.Host = parts[0]
		if parts[1] != "" {
			port, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("invalid port number: %s", parts[1])
			}
			if port < 1 || port > 65535 {
				return fmt.Errorf("port number out of range (1-65535): %d", port)
			}
			hpv.Port = port
		}
	}

	// Validate host if provided
	if hpv.Host != "" {
		if net.ParseIP(hpv.Host) == nil {
			// Not an IP, check if it's a valid hostname
			if !isValidHostname(hpv.Host) {
				return fmt.Errorf("invalid host: %s", hpv.Host)
			}
		}
	}

	return nil
}

func (hpv *HostPortValue) Type() string {
	return "host:port"
}

func (hpv *HostPortValue) Get() interface{} {
	return struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}{
		Host: hpv.Host,
		Port: hpv.Port,
	}
}

func (hpv *HostPortValue) GetHost() string {
	return hpv.Host
}

func (hpv *HostPortValue) GetPort() int {
	return hpv.Port
}

func (hpv *HostPortValue) GetAddress() string {
	return net.JoinHostPort(hpv.Host, strconv.Itoa(hpv.Port))
}

func isValidHostname(hostname string) bool {
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}

	// Simple hostname validation
	for _, r := range hostname {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '.') {
			return false
		}
	}

	return true
}

// IPValue implements pflag.Value for IP addresses
type IPValue struct {
	IP net.IP
}

func NewIPValue(defaultIP string) *IPValue {
	ipv := &IPValue{}
	if defaultIP != "" {
		_ = ipv.Set(defaultIP)
	}
	return ipv
}

func (ipv *IPValue) String() string {
	if ipv.IP == nil {
		return ""
	}
	return ipv.IP.String()
}

func (ipv *IPValue) Set(s string) error {
	if s == "" {
		ipv.IP = nil
		return nil
	}

	ip := net.ParseIP(s)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", s)
	}

	ipv.IP = ip
	return nil
}

func (ipv *IPValue) Type() string {
	return "ip"
}

func (ipv *IPValue) Get() interface{} {
	return ipv.IP
}

func (ipv *IPValue) GetIP() net.IP {
	return ipv.IP
}

func (ipv *IPValue) IsIPv4() bool {
	return ipv.IP != nil && ipv.IP.To4() != nil
}

func (ipv *IPValue) IsIPv6() bool {
	return ipv.IP != nil && ipv.IP.To4() == nil
}

// IsSet returns true if IP address has been set
func (ipv *IPValue) IsSet() bool {
	return ipv.IP != nil
}

// CIDRValue implements pflag.Value for CIDR notation
type CIDRValue struct {
	IPNet *net.IPNet
}

func NewCIDRValue(defaultCIDR string) *CIDRValue {
	cv := &CIDRValue{}
	if defaultCIDR != "" {
		_ = cv.Set(defaultCIDR)
	}
	return cv
}

func (cv *CIDRValue) String() string {
	if cv.IPNet == nil {
		return ""
	}
	return cv.IPNet.String()
}

func (cv *CIDRValue) Set(s string) error {
	if s == "" {
		cv.IPNet = nil
		return nil
	}

	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return fmt.Errorf("invalid CIDR notation: %s", s)
	}

	cv.IPNet = ipnet
	return nil
}

func (cv *CIDRValue) Type() string {
	return "cidr"
}

func (cv *CIDRValue) Get() interface{} {
	return cv.IPNet
}

func (cv *CIDRValue) GetIPNet() *net.IPNet {
	return cv.IPNet
}

func (cv *CIDRValue) Contains(ip net.IP) bool {
	if cv.IPNet == nil {
		return false
	}
	return cv.IPNet.Contains(ip)
}

// PortRangeValue implements pflag.Value for port ranges
type PortRangeValue struct {
	Start int
	End   int
}

func NewPortRangeValue(defaultStart, defaultEnd int) *PortRangeValue {
	return &PortRangeValue{
		Start: defaultStart,
		End:   defaultEnd,
	}
}

func (prv *PortRangeValue) String() string {
	if prv.Start == 0 && prv.End == 0 {
		return ""
	}
	if prv.Start == prv.End {
		return strconv.Itoa(prv.Start)
	}
	return fmt.Sprintf("%d-%d", prv.Start, prv.End)
}

func (prv *PortRangeValue) Set(s string) error {
	if s == "" {
		prv.Start = 0
		prv.End = 0
		return nil
	}

	if !strings.Contains(s, "-") {
		// Single port
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid port number: %s", s)
		}

		if port < 1 || port > 65535 {
			return fmt.Errorf("port number out of range (1-65535): %d", port)
		}

		prv.Start = port
		prv.End = port
		return nil
	}

	// Port range
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid port range format: %s (expected start-end)", s)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid start port: %s", parts[0])
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid end port: %s", parts[1])
	}

	if start < 1 || start > 65535 {
		return fmt.Errorf("start port out of range (1-65535): %d", start)
	}

	if end < 1 || end > 65535 {
		return fmt.Errorf("end port out of range (1-65535): %d", end)
	}

	if start > end {
		return fmt.Errorf("start port must be less than or equal to end port")
	}

	prv.Start = start
	prv.End = end
	return nil
}

func (prv *PortRangeValue) Type() string {
	return "port-range"
}

func (prv *PortRangeValue) Get() interface{} {
	return struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}{
		Start: prv.Start,
		End:   prv.End,
	}
}

func (prv *PortRangeValue) GetStart() int {
	return prv.Start
}

func (prv *PortRangeValue) GetEnd() int {
	return prv.End
}

func (prv *PortRangeValue) Contains(port int) bool {
	return port >= prv.Start && port <= prv.End
}

func (prv *PortRangeValue) Size() int {
	if prv.Start == 0 && prv.End == 0 {
		return 0
	}
	return prv.End - prv.Start + 1
}
