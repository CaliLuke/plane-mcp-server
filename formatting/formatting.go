// Package formatting provides compact single-line formatters for MCP list responses.
package formatting

import (
	"fmt"
	"strings"

	"github.com/CaliLuke/plane-mcp-server/uid"
	"golang.org/x/net/html"
)

// Project formats: "IDENTIFIER · Name [short_id]"
func Project(id, name, identifier string) string {
	short := uid.Encode(id)
	if short == "" {
		short = "?"
	}
	return fmt.Sprintf("%s · %s [%s]", identifier, name, short)
}

// WorkItem formats: "PROJ-seq (state) · Name [short_id]"
func WorkItem(id string, sequenceID int, name, projectIdentifier, stateName string) string {
	short := uid.Encode(id)
	if short == "" {
		short = "?"
	}
	var seq string
	if projectIdentifier != "" {
		seq = fmt.Sprintf("%s-%d", projectIdentifier, sequenceID)
	} else {
		seq = fmt.Sprintf("#%d", sequenceID)
	}
	statePart := ""
	if stateName != "" {
		statePart = fmt.Sprintf(" (%s)", stateName)
	}
	return fmt.Sprintf("%s%s · %s [%s]", seq, statePart, name, short)
}

// State formats: "Name (group) [short_id]" — marks default with *
func State(id, name, group string, isDefault bool) string {
	short := uid.Encode(id)
	if short == "" {
		short = "?"
	}
	def := ""
	if isDefault {
		def = " *"
	}
	return fmt.Sprintf("%s (%s)%s [%s]", name, group, def, short)
}

// StripHTML removes HTML tags and returns plain text, or empty string if empty.
func StripHTML(s string) string {
	if s == "" {
		return ""
	}
	tokenizer := html.NewTokenizer(strings.NewReader(s))
	var parts []string
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.TextToken {
			text := strings.TrimSpace(tokenizer.Token().Data)
			if text != "" {
				parts = append(parts, text)
			}
		}
	}
	return strings.Join(parts, " ")
}
