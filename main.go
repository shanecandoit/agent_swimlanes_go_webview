package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/abemedia/go-webview"
	_ "github.com/abemedia/go-webview/embedded"
)

var (
	Port        = "8080"
	DefaultYaml = "example1.yaml"
)

func main() {
	// Start the REST API server in a background goroutine
	go startServer()

	// Create a new webview instance
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("Agent Swimlanes")
	w.SetSize(1280, 800, webview.HintNone)

	// Load index.html content
	content, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Printf("Error: could not read index.html: %v\n", err)
		w.SetHtml("<html><body><h1>Error: index.html not found</h1></body></html>")
	} else {
		w.SetHtml(string(content))
	}

	w.Run()
}

func startServer() {
	handler := createRouter()
	fmt.Printf("REST API server starting on :%s\n", Port)
	if err := http.ListenAndServe(":"+Port, handler); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}

func createRouter() http.Handler {
	mux := http.NewServeMux()

	// 1. Graph State: GET /api/workflow/graph
	mux.HandleFunc("/api/workflow/graph", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			workflow, err := LoadWorkflow(DefaultYaml)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(workflow)
		}
	})

	// 2. Add Nodes: POST /api/workflow/nodes
	mux.HandleFunc("/api/workflow/nodes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var node Node
			if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			workflow, _ := LoadWorkflow(DefaultYaml)
			workflow.Nodes = append(workflow.Nodes, node)
			SaveWorkflow(DefaultYaml, workflow)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(node)
		}
	})

	// 3. Add Edges: POST /api/workflow/edges
	mux.HandleFunc("/api/workflow/edges", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var edge Edge
			if err := json.NewDecoder(r.Body).Decode(&edge); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			workflow, _ := LoadWorkflow(DefaultYaml)
			workflow.Edges = append(workflow.Edges, edge)
			SaveWorkflow(DefaultYaml, workflow)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(edge)
		}
	})

	// 4. Decision Logic: POST /api/nodes/{id}/evaluate
	mux.HandleFunc("/api/nodes/evaluate/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			nodeID := r.URL.Path[len("/api/nodes/evaluate/"):]
			var payload struct {
				InputMarkdown string `json:"current_input_markdown"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Mocking the result for now: choose the first outgoing edge
			workflow, _ := LoadWorkflow(DefaultYaml)
			var targetNodeID string
			var selectedEdgeID string
			for _, edge := range workflow.Edges {
				if edge.FromNodeID == nodeID {
					targetNodeID = edge.ToNodeID
					selectedEdgeID = edge.ID
					break
				}
			}

			if targetNodeID == "" {
				http.Error(w, "No outgoing edges found for this decision node", http.StatusNotFound)
				return
			}

			// Save execution log
			log := ExecutionLog{
				ID:             "log-" + nodeID,
				NodeID:         nodeID,
				Reasoning:      "Mock reasoning: Selecting the first available path based on input evidence.",
				SelectedEdgeID: selectedEdgeID,
			}
			fmt.Printf("Decision log saved: %+v\n", log)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"target_node_id": targetNodeID,
				"reasoning":      log.Reasoning,
			})
		}
	})

	// 5. Execution Logs: GET /api/execution/logs/:node_id
	mux.HandleFunc("/api/execution/logs/", func(w http.ResponseWriter, r *http.Request) {
		nodeID := r.URL.Path[len("/api/execution/logs/"):]
		// Returns mock log for now or can be extended with real storage
		log := ExecutionLog{
			ID:             "log-" + nodeID,
			NodeID:         nodeID,
			Reasoning:      "Mock reasoning: The agent analyzed the context and selected the optimal path.",
			SelectedEdgeID: "edge-x",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(log)
	})

	// 6. Agent Personas: GET /api/agents
	mux.HandleFunc("/api/agents", func(w http.ResponseWriter, r *http.Request) {
		workflow, _ := LoadWorkflow(DefaultYaml)
		agents := make(map[string]string)
		for _, node := range workflow.Nodes {
			if _, exists := agents[node.SwimlaneID]; !exists {
				agents[node.SwimlaneID] = node.SystemPrompt
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(agents)
	})

	// 7. Static Assets (JS/CSS)
	mux.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))

	// 8. Document Content
	mux.HandleFunc("/api/docs/", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Path[len("/api/docs/"):]
		workflow, err := LoadWorkflow(DefaultYaml)
		if err != nil {
			http.Error(w, "Workflow not loaded", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodGet {
			content, err := GetDocContent(workflow.DocFolder, filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/markdown")
			w.Write([]byte(content))
		}
	})

	// Health check
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Agent Swimlanes Engine: Online")
	})

	// Return a middleware that adds CORS headers to allow the WebView ('null' origin) to call the local API
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		mux.ServeHTTP(w, r)
	})
}
