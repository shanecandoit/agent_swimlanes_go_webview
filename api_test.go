package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAPIRoutes(t *testing.T) {
	// Setup a temporary workflow file for testing
	testYaml := "test_api.yaml"
	originalYaml := DefaultYaml
	DefaultYaml = testYaml
	defer func() { DefaultYaml = originalYaml }()

	initialWorkflow := WorkflowSpec{
		ID:        "test-wf",
		Name:      "Test API Workflow",
		DocFolder: "test_docs",
		Nodes: []Node{
			{ID: "node-1", Type: NodeTypeTask, Label: "Start"},
			{ID: "node-2", Type: NodeTypeTask, Label: "End"},
		},
		Edges: []Edge{
			{ID: "edge-1", FromNodeID: "node-1", ToNodeID: "node-2", Label: "Next"},
		},
	}
	SaveWorkflow(testYaml, &initialWorkflow)
	defer os.Remove(testYaml)

	router := createRouter()

	// 1. Test GET /api/workflow/graph
	req, _ := http.NewRequest("GET", "/api/workflow/graph", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/workflow/graph returned %v", rr.Code)
	}

	var loaded WorkflowSpec
	json.Unmarshal(rr.Body.Bytes(), &loaded)
	if loaded.Name != initialWorkflow.Name {
		t.Errorf("Expected workflow name %s, got %s", initialWorkflow.Name, loaded.Name)
	}

	// 2. Test POST /api/workflow/nodes
	newNode := Node{ID: "node-3", Type: NodeTypeTask, Label: "Extra Node"}
	body, _ := json.Marshal(newNode)
	req, _ = http.NewRequest("POST", "/api/workflow/nodes", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /api/workflow/nodes returned %v", rr.Code)
	}

	// Verify node added
	workflow, _ := LoadWorkflow(testYaml)
	if len(workflow.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(workflow.Nodes))
	}

	// 3. Test POST /api/nodes/evaluate/node-1
	evaluatePayload := map[string]string{"current_input_markdown": "# Evidence"}
	body, _ = json.Marshal(evaluatePayload)
	req, _ = http.NewRequest("POST", "/api/nodes/evaluate/node-1", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("POST /api/nodes/evaluate/node-1 returned %v", rr.Code)
	}

	var result map[string]string
	json.Unmarshal(rr.Body.Bytes(), &result)
	if result["target_node_id"] != "node-2" {
		t.Errorf("Expected target node id node-2, got %s", result["target_node_id"])
	}

	// 4. Test POST /api/workflow/edges (NEW)
	newEdge := Edge{ID: "edge-x", FromNodeID: "node-2", ToNodeID: "node-3", Label: "Final"}
	body, _ = json.Marshal(newEdge)
	req, _ = http.NewRequest("POST", "/api/workflow/edges", bytes.NewBuffer(body))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("POST /api/workflow/edges returned %v", rr.Code)
	}

	// 5. Test GET /api/agents (NEW)
	req, _ = http.NewRequest("GET", "/api/agents", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/agents returned %v", rr.Code)
	}

	// 6. Test GET /api/execution/logs/:id (NEW)
	req, _ = http.NewRequest("GET", "/api/execution/logs/node-1", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/execution/logs/node-1 returned %v", rr.Code)
	}

	// 7. Test CORS Headers (NEW)
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin: *, got %s", rr.Header().Get("Access-Control-Allow-Origin"))
	}

	// 8. Test Health Check (NEW)
	req, _ = http.NewRequest("GET", "/api/health", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GET /api/health returned %v", rr.Code)
	}
}
