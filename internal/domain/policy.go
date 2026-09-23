package domain

import (
	"fmt"
	"time"
)

// Default delivery policy values, used when the notifier config omits them.
const (
	DefaultMinFailures = 3
	DefaultMaxSends    = 1
)

// NotifierTypeSMSPilot is the SMSPILOT voice-call channel.
const NotifierTypeSMSPilot = "smspilot"

// IsPolicyRouted reports whether a notifier type is dispatched from incident
// state (delivery policy) rather than from the per-check alert fan-out.
// Expensive, real-world channels live here: a phone call must not ring on the
// first blip, on a slow response, or on every alert type.
func IsPolicyRouted(notifierType string) bool {
	return notifierType == NotifierTypeSMSPilot
}

// DeliveryPolicy controls when a notifier is allowed to fire for an ongoing
// incident. It is parsed from the notifier's own config JSON, so every channel
// carries its own rules and no extra table is needed.
//
// Notifiers with a policy are dispatched from the incident state (see
// alerting.Manager) instead of the per-check alert fan-out.
type DeliveryPolicy struct {
	// MinFailures is how many failed checks an incident must accumulate before
	// this notifier fires. 1 means "notify on the first failure".
	MinFailures int

	// MaxSends caps how many notifications a single incident may produce.
	MaxSends int

	// RepeatInterval is the delay between repeated notifications for the same
	// incident. Zero means "notify once".
	RepeatInterval time.Duration

	// TargetIDs limits the policy to specific targets. Empty means all targets.
	TargetIDs []string
}

// ParseDeliveryPolicy reads the policy fields from a notifier config.
func ParseDeliveryPolicy(config map[string]interface{}) (DeliveryPolicy, error) {
	policy := DeliveryPolicy{
		MinFailures: DefaultMinFailures,
		MaxSends:    DefaultMaxSends,
	}

	if raw, ok := config["min_failures"]; ok {
		n, err := policyInt(raw)
		if err != nil {
			return policy, fmt.Errorf("min_failures must be a number")
		}
		if n < 1 {
			return policy, fmt.Errorf("min_failures must be >= 1")
		}
		policy.MinFailures = n
	}

	if raw, ok := config["max_sends"]; ok {
		n, err := policyInt(raw)
		if err != nil {
			return policy, fmt.Errorf("max_sends must be a number")
		}
		if n < 1 {
			return policy, fmt.Errorf("max_sends must be >= 1")
		}
		policy.MaxSends = n
	}

	if raw, ok := config["repeat_interval"]; ok {
		value, ok := raw.(string)
		if !ok {
			return policy, fmt.Errorf("repeat_interval must be a string (e.g. \"30m\")")
		}
		if value != "" {
			d, err := time.ParseDuration(value)
			if err != nil {
				return policy, fmt.Errorf("invalid repeat_interval: %w", err)
			}
			if d < 0 {
				return policy, fmt.Errorf("repeat_interval must not be negative")
			}
			policy.RepeatInterval = d
		}
	}

	if raw, ok := config["target_ids"]; ok {
		ids, err := policyStrings(raw)
		if err != nil {
			return policy, fmt.Errorf("target_ids must be an array of strings")
		}
		policy.TargetIDs = ids
	}

	return policy, nil
}

// AppliesTo reports whether the policy covers the given target.
func (p DeliveryPolicy) AppliesTo(targetID string) bool {
	if len(p.TargetIDs) == 0 {
		return true
	}
	for _, id := range p.TargetIDs {
		if id == targetID {
			return true
		}
	}
	return false
}

func policyInt(raw interface{}) (int, error) {
	switch v := raw.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case int64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("not a number")
	}
}

func policyStrings(raw interface{}) ([]string, error) {
	items, ok := raw.([]interface{})
	if !ok {
		if values, ok := raw.([]string); ok {
			return values, nil
		}
		return nil, fmt.Errorf("not an array")
	}

	values := make([]string, 0, len(items))
	for _, item := range items {
		value, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("not a string")
		}
		values = append(values, value)
	}
	return values, nil
}
