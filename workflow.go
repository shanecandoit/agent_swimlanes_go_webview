package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type NodeType string

const (
	NodeTypeTask     NodeType = "TASK"
	NodeTypeDecision NodeType = "DECISION"
)

type Node struct {
	ID           string   `yaml:"id" json:"id"`
	Type         NodeType `yaml:"type" json:"type"`
	SwimlaneID   string   `yaml:"swimlane_id" json:"swimlane_id"`
	Label        string   `yaml:"label" json:"label"`
	SystemPrompt string   `yaml:"system_prompt" json:"system_prompt"`
	Filename     string   `yaml:"filename,omitempty" json:"filename,omitempty"`
}

type Edge struct {
	ID         string `yaml:"id" json:"id"`
	FromNodeID string `yaml:"from_node_id" json:"from_node_id"`
	ToNodeID   string `yaml:"to_node_id" json:"to_node_id"`
	Label      string `yaml:"label" json:"label"`
}

type ExecutionLog struct {
	ID             string `yaml:"id" json:"id"`
	NodeID         string `yaml:"node_id" json:"node_id"`
	Reasoning      string `yaml:"reasoning" json:"reasoning"`
	SelectedEdgeID string `yaml:"selected_edge_id" json:"selected_edge_id"`
}

type WorkflowSpec struct {
	ID        string `yaml:"id" json:"id"`
	Name      string `yaml:"name" json:"name"`
	DocFolder string `yaml:"doc_folder" json:"doc_folder"`
	Nodes     []Node `yaml:"nodes" json:"nodes"`
	Edges     []Edge `yaml:"edges" json:"edges"`
}

// LoadWorkflow reads the YAML configuration from the specified path.
func LoadWorkflow(path string) (*WorkflowSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var workflow WorkflowSpec
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return &workflow, nil
}

// SaveWorkflow writes the workflow spec back to a YAML file.
func SaveWorkflow(path string, workflow *WorkflowSpec) error {
	data, err := yaml.Marshal(workflow)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// GetDocContent reads the content of a document from the workflow's document folder.
func GetDocContent(docFolder, filename string) (string, error) {
	path := filepath.Join(docFolder, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read doc: %w", err)
	}

	return string(data), nil
}

// SaveDocContent writes markdown content to the specified file in the document folder.
func SaveDocContent(docFolder, filename, content string) error {
	// Ensure directory exists
	if err := os.MkdirAll(docFolder, 0755); err != nil {
		return fmt.Errorf("mkdir doc-folder: %w", err)
	}
	path := filepath.Join(docFolder, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write doc: %w", err)
	}

	return nil
}
