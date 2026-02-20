// Package uid provides short UUID encoding/decoding for compact MCP responses.
//
// UUIDs are encoded to ~22-char base57 strings in all outgoing responses.
// Incoming parameters accept both short and full UUID formats transparently.
package uid

import (
	"regexp"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

var uuidRE = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// Encode converts a full UUID string to a short base57 string.
func Encode(uuidStr string) string {
	if uuidStr == "" {
		return ""
	}
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuidStr
	}
	return shortuuid.DefaultEncoder.Encode(u)
}

// Decode accepts a short or full UUID string; always returns a full UUID string.
func Decode(s string) string {
	if uuidRE.MatchString(s) {
		return s
	}
	u, err := shortuuid.DefaultEncoder.Decode(s)
	if err != nil {
		return s
	}
	return u.String()
}

// IsUUID returns true if the string looks like a full UUID.
func IsUUID(s string) bool {
	return uuidRE.MatchString(s)
}

// IsUUIDLike returns true if s looks like a full UUID or a shortUUID (base57, ~22 chars).
func IsUUIDLike(s string) bool {
	if IsUUID(s) {
		return true
	}
	if len(s) >= 20 {
		for _, c := range s {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
				return false
			}
		}
		return true
	}
	return false
}

// EncodeValue recursively encodes any UUID strings found in maps/slices/scalars.
func EncodeValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, vv := range val {
			out[k] = EncodeValue(vv)
		}
		return out
	case []any:
		out := make([]any, len(val))
		for i, vv := range val {
			out[i] = EncodeValue(vv)
		}
		return out
	case string:
		if uuidRE.MatchString(val) {
			return Encode(val)
		}
		return val
	default:
		return v
	}
}
