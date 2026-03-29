package main

import (
	"path/filepath"
	"testing"
)

func TestLoadSaveWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.yaml")

	spec := &WorkflowSpec{
		ID:        "test-wf",
		Name:      "Test Workflow",
		DocFolder: "test-docs",
		Nodes: []Node{
			{
				ID:         "lane1-task1",
				Type:       NodeTypeTask,
				SwimlaneID: "lane1",
				Label:      "Test Node 1",
				Filename:   "task1.md",
			},
		},
		Edges: []Edge{
			{
				ID:         "edge1",
				FromNodeID: "lane1-task1",
				ToNodeID:   "lane1-task2",
				Label:      "Success",
			},
		},
	}

	// Test Save
	if err := SaveWorkflow(path, spec); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Test Load
	loaded, err := LoadWorkflow(path)
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if loaded.Name != spec.Name {
		t.Errorf("Expected name %s, got %s", spec.Name, loaded.Name)
	}

	if len(loaded.Nodes) != 1 || loaded.Nodes[0].ID != "lane1-task1" {
		t.Errorf("Failed to load node data correctly")
	}

	if loaded.Nodes[0].Filename != "task1.md" {
		t.Errorf("Failed to load task data correctly")
	}

	if len(loaded.Edges) != 1 || loaded.Edges[0].ID != "edge1" {
		t.Errorf("Failed to load edge data correctly")
	}
}

func TestDocOperations(t *testing.T) {
	tmpDir := t.TempDir()
	filename := "test.md"
	content := "# Hello World\nThis is a test document."

	// Test SaveDocContent
	if err := SaveDocContent(tmpDir, filename, content); err != nil {
		t.Fatalf("Failed to save doc: %v", err)
	}

	// Test GetDocContent
	loaded, err := GetDocContent(tmpDir, filename)
	if err != nil {
		t.Fatalf("Failed to load doc: %v", err)
	}

	if loaded != content {
		t.Errorf("Expected content '%s', got '%s'", content, loaded)
	}
}
