// Package models provides slim response models for Plane MCP tools,
// using godantic for validation and schema generation.
package models

import (
	"encoding/json"

	"github.com/CaliLuke/plane-mcp-server/uid"
	"github.com/deepankarm/godantic/pkg/godantic"
)

// ProjectSummary — id, name, identifier only.
type ProjectSummary struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

func (p *ProjectSummary) FieldName() godantic.FieldOptions[string] {
	return godantic.Field(godantic.Required[string]())
}

func (p *ProjectSummary) FieldIdentifier() godantic.FieldOptions[string] {
	return godantic.Field(godantic.Required[string]())
}

// AssigneeSummary for work item assignees.
type AssigneeSummary struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// LabelSummary for work item labels.
type LabelSummary struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}

// WorkItemSummary — list-level work item, no description.
type WorkItemSummary struct {
	ID         string   `json:"id,omitempty"`
	SequenceID int      `json:"sequence_id,omitempty"`
	Name       string   `json:"name"`
	Priority   string   `json:"priority,omitempty"`
	State      string   `json:"state,omitempty"`
	Assignees  []string `json:"assignees,omitempty"`
	Labels     []string `json:"labels,omitempty"`
}

func (w *WorkItemSummary) FieldName() godantic.FieldOptions[string] {
	return godantic.Field(godantic.Required[string]())
}

// WorkItemFull — detail view with description, expanded assignees/labels, dates.
type WorkItemFull struct {
	ID          string            `json:"id,omitempty"`
	SequenceID  int               `json:"sequence_id,omitempty"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Priority    string            `json:"priority,omitempty"`
	State       string            `json:"state,omitempty"`
	Assignees   []AssigneeSummary `json:"assignees,omitempty"`
	Labels      []LabelSummary    `json:"labels,omitempty"`
	Parent      string            `json:"parent,omitempty"`
	StartDate   string            `json:"start_date,omitempty"`
	TargetDate  string            `json:"target_date,omitempty"`
}

func (w *WorkItemFull) FieldName() godantic.FieldOptions[string] {
	return godantic.Field(godantic.Required[string]())
}

// StateSummary — id, name, group, default.
type StateSummary struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Group   string `json:"group,omitempty"`
	Default bool   `json:"default,omitempty"`
}

func (s *StateSummary) FieldName() godantic.FieldOptions[string] {
	return godantic.Field(godantic.Required[string]())
}

// Slim serializes a model to a map, drops zero/empty values, encodes UUIDs to short form.
func Slim(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	encoded := uid.EncodeValue(m)
	return encoded.(map[string]any), nil
}

// SlimJSON serializes a model to JSON string with short UUIDs.
func SlimJSON(v any) (string, error) {
	m, err := Slim(v)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
