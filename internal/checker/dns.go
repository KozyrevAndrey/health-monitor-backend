package checker

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"health-monitor/internal/domain"
)

// dnsRecordTypes lists the DNS record types the checker can resolve.
var dnsRecordTypes = map[string]bool{
	"A": true, "AAAA": true, "CNAME": true, "MX": true, "TXT": true, "NS": true,
}

// DNSChecker resolves DNS records for a domain, measures resolution time and
// optionally validates the result against expected values.
type DNSChecker struct{}

func NewDNSChecker() *DNSChecker {
	return &DNSChecker{}
}

func (c *DNSChecker) Check(ctx context.Context, target *domain.Target) (*domain.CheckResult, error) {
	dnsConfig, err := parseDNSConfig(target.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid DNS config: %w", err)
	}

	result := &domain.CheckResult{
		TargetID:  target.ID,
		CheckedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
	result.Metadata["domain"] = dnsConfig.Domain
	result.Metadata["record_type"] = dnsConfig.RecordType
	if dnsConfig.DNSServer != "" {
		result.Metadata["dns_server"] = dnsConfig.DNSServer
	}

	checkCtx, cancel := context.WithTimeout(ctx, target.Timeout)
	defer cancel()

	resolver := buildResolver(dnsConfig.DNSServer)

	startTime := time.Now()
	values, err := resolveRecords(checkCtx, resolver, dnsConfig.RecordType, dnsConfig.Domain)
	resolveTime := time.Since(startTime)
	result.ResponseTimeMs = resolveTime.Milliseconds()

	if err != nil {
		result.Status = domain.CheckStatusFailure
		result.Error = err.Error()
		result.Message = fmt.Sprintf("DNS resolution of %s %s failed: %s", dnsConfig.RecordType, dnsConfig.Domain, err.Error())
		return result, nil
	}

	values = uniqueStrings(values)
	sort.Strings(values)
	result.Metadata["records"] = values
	result.Metadata["record_count"] = len(values)

	if len(values) == 0 {
		result.Status = domain.CheckStatusFailure
		result.Message = fmt.Sprintf("No %s records found for %s", dnsConfig.RecordType, dnsConfig.Domain)
		return result, nil
	}

	expected := append(append([]string{}, dnsConfig.ExpectedIPs...), dnsConfig.ExpectValues...)
	if missing := missingValues(expected, values); len(missing) > 0 {
		result.Status = domain.CheckStatusFailure
		result.Message = fmt.Sprintf("Expected %s record(s) not found: %s (got: %s)",
			dnsConfig.RecordType, strings.Join(missing, ", "), strings.Join(values, ", "))
		return result, nil
	}

	result.Status = domain.CheckStatusSuccess
	result.Message = fmt.Sprintf("DNS %s for %s resolved to %d record(s) in %dms",
		dnsConfig.RecordType, dnsConfig.Domain, len(values), resolveTime.Milliseconds())

	return result, nil
}

func (c *DNSChecker) Type() domain.TargetType {
	return domain.TargetTypeDNS
}

func (c *DNSChecker) Validate(config map[string]interface{}) error {
	_, err := parseDNSConfig(config)
	return err
}

// buildResolver returns the default resolver, or one that queries a custom DNS
// server when dnsServer is set (host or host:port; :53 is assumed by default).
func buildResolver(dnsServer string) *net.Resolver {
	if dnsServer == "" {
		return net.DefaultResolver
	}

	addr := dnsServer
	if _, _, err := net.SplitHostPort(addr); err != nil {
		addr = net.JoinHostPort(addr, "53")
	}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, network, addr)
		},
	}
}

func resolveRecords(ctx context.Context, r *net.Resolver, recordType, domainName string) ([]string, error) {
	switch recordType {
	case "A":
		return lookupIPs(ctx, r, "ip4", domainName)
	case "AAAA":
		return lookupIPs(ctx, r, "ip6", domainName)
	case "CNAME":
		cname, err := r.LookupCNAME(ctx, domainName)
		if err != nil {
			return nil, err
		}
		return []string{trimDot(cname)}, nil
	case "MX":
		mxs, err := r.LookupMX(ctx, domainName)
		if err != nil {
			return nil, err
		}
		values := make([]string, 0, len(mxs))
		for _, mx := range mxs {
			values = append(values, fmt.Sprintf("%s (pref %d)", trimDot(mx.Host), mx.Pref))
		}
		return values, nil
	case "TXT":
		return r.LookupTXT(ctx, domainName)
	case "NS":
		nss, err := r.LookupNS(ctx, domainName)
		if err != nil {
			return nil, err
		}
		values := make([]string, 0, len(nss))
		for _, ns := range nss {
			values = append(values, trimDot(ns.Host))
		}
		return values, nil
	default:
		return nil, fmt.Errorf("unsupported record type: %s", recordType)
	}
}

func lookupIPs(ctx context.Context, r *net.Resolver, network, domainName string) ([]string, error) {
	ips, err := r.LookupIP(ctx, network, domainName)
	if err != nil {
		return nil, err
	}
	values := make([]string, 0, len(ips))
	for _, ip := range ips {
		values = append(values, ip.String())
	}
	return values, nil
}

// missingValues returns the expected values that are not present in got.
// Comparison is case-insensitive and tolerant of trailing dots; an expected IP
// or hostname matches if it appears as a substring-free token in any record
// (records like "mx1.example.com (pref 10)" still match "mx1.example.com").
func missingValues(expected, got []string) []string {
	var missing []string
	for _, e := range expected {
		want := normalizeDNS(e)
		if want == "" {
			continue
		}
		found := false
		for _, g := range got {
			gn := normalizeDNS(g)
			if gn == want || strings.HasPrefix(gn, want+" ") {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, e)
		}
	}
	return missing
}

func normalizeDNS(s string) string {
	return strings.ToLower(trimDot(strings.TrimSpace(s)))
}

// uniqueStrings returns values with duplicates removed, preserving first-seen order.
func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func trimDot(s string) string {
	return strings.TrimSuffix(s, ".")
}

func parseDNSConfig(config map[string]interface{}) (*domain.DNSConfig, error) {
	dnsConfig := &domain.DNSConfig{RecordType: "A"}

	domainName, ok := config["domain"].(string)
	if !ok || domainName == "" {
		return nil, fmt.Errorf("domain is required")
	}
	dnsConfig.Domain = domainName

	if rt, ok := config["record_type"].(string); ok && rt != "" {
		dnsConfig.RecordType = strings.ToUpper(rt)
	}
	if !dnsRecordTypes[dnsConfig.RecordType] {
		return nil, fmt.Errorf("unsupported record_type %q (supported: A, AAAA, CNAME, MX, TXT, NS)", dnsConfig.RecordType)
	}

	if server, ok := config["dns_server"].(string); ok {
		dnsConfig.DNSServer = server
	}

	dnsConfig.ExpectedIPs = parseStringSlice(config["expected_ips"])
	dnsConfig.ExpectValues = parseStringSlice(config["expect_values"])

	return dnsConfig, nil
}

// parseStringSlice coerces a config value into a []string, accepting both
// []interface{} (from JSON) and []string.
func parseStringSlice(raw interface{}) []string {
	switch v := raw.(type) {
	case []interface{}:
		out := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return v
	default:
		return nil
	}
}
