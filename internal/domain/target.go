package domain

import "time"

// TargetType represents the type of monitoring target
type TargetType string

const (
	TargetTypeHTTP   TargetType = "http"
	TargetTypeTCP    TargetType = "tcp"
	TargetTypeDNS    TargetType = "dns"
	TargetTypeICMP   TargetType = "icmp"
	TargetTypeScript TargetType = "script"
)

// Target represents a monitoring target
type Target struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        TargetType             `json:"type"`
	Config      map[string]interface{} `json:"config"`
	Interval    time.Duration          `json:"interval"`
	Timeout     time.Duration          `json:"timeout"`
	Enabled     bool                   `json:"enabled"`
	Tags        []string               `json:"tags,omitempty"`
	Description string                 `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// HTTPConfig represents HTTP-specific configuration
type HTTPConfig struct {
	URL                string            `json:"url"`
	Method             string            `json:"method"`
	Headers            map[string]string `json:"headers,omitempty"`
	Body               string            `json:"body,omitempty"`
	ExpectedStatusCode int               `json:"expected_status_code"`
	MaxResponseTimeMs  int               `json:"max_response_time_ms,omitempty"`
	FollowRedirects    bool              `json:"follow_redirects"`
	ValidateSSL        bool              `json:"validate_ssl"`
	CheckSSLExpiry     bool              `json:"check_ssl_expiry"`
	SSLExpiryDays      int               `json:"ssl_expiry_days,omitempty"`
}

// TCPConfig represents TCP-specific configuration
type TCPConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// DNSConfig represents DNS-specific configuration
type DNSConfig struct {
	Domain       string   `json:"domain"`
	RecordType   string   `json:"record_type"`
	ExpectedIPs  []string `json:"expected_ips,omitempty"`
	DNSServer    string   `json:"dns_server,omitempty"`
	ExpectValues []string `json:"expect_values,omitempty"`
}

// ICMPConfig represents ICMP/Ping-specific configuration
type ICMPConfig struct {
	Host        string `json:"host"`
	PacketCount int    `json:"packet_count"`
	MaxRTTMs    int    `json:"max_rtt_ms,omitempty"`
}
