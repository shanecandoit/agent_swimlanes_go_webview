# 🌊 Agent Swimlanes (Go + WebView)

![Dashboard Mockup](./assets/dashboard_mockup.png)

**Agent Swimlanes** is a high-performance **Workflow Logic Engine** designed to orchestrate and visualize complex multi-agent processes. Instead of simple task lists, it models AI workflows as a **polymorphic graph** of nodes and edges, allowing for sophisticated branching, decision-making, and human-in-the-loop oversight.

Built with **Go** for orchestration and a **web-based canvas UI** via `webview_go`, Agent Swimlanes provides a transparent, documented path for agentic reasoning from initial concept to finalized output.

---

## 🧩 Core Concepts

To enable true "Workflow Logic," the application uses a structured graph model:

### 🔄 Workflows
The top-level container for an end-to-end process. A workflow is a directed graph where data or state flows between different logic nodes.

### 🏊 Swimlanes
Visual and logical groupings by **Persona**. Swimlanes help organize the workflow visually, ensuring that each step is clearly attributed to a specific agent role (e.g., "Technical Architect", "Legal Reviewer").

### 🤖 Agents
The "brains" behind each swimlane. Each agent is defined by a specific **persona** and **system prompt**. Agents are responsible for both performing tasks and evaluating decisions based on their specialized expertise.

### 📝 Tasks (Action Nodes)
Deterministic nodes where an agent performs a specific action. Each task link to a Markdown document representing the context, instructions, and output of that specific step.

### 💎 Decisions (Branching Nodes)
Specialized "fork in the road" nodes. Decisions use LLM logic to evaluate current workspace evidence against the node's system prompt to determine the next path.
*   **Reasoning-First**: Every decision execution generates an **Execution Log**, capturing the agent's "internal monologue" before choosing a path.

### 🛤️ Edges
The connections defining the flow between nodes. Edges for decision nodes have descriptive labels (e.g., "Approved", "Needs Revision") which the agent selects based on its evaluation.

---

## ✨ Key Features

*   **Visual Orchestrator**: A canvas-based diagramming interface for dragging-and-dropping tasks and connecting decisions.
*   **Live Reasoning Feed**: Watch agents evaluate evidence in real-time as they navigate through decision points.
*   **Manual Override**: Human-in-the-loop capability to pause execution and manually choose a different path than the one recommended by an agent.
*   **History & Auditing**: Color-coded flow paths highlight the "Path Taken," while grayed-out edges keep the discarded reasoning trails visible for full transparency.
*   **Native Performance**: Low-memory footprint Go backend with a responsive glassmorphic UI.

---

## 🏗️ Technical Architecture

The application follows a **Decoupled Engine/UI** model:

1.  **Go Backend (The Engine)**:
    *   Maintains the internal state of the graph (Nodes, Edges, Execution Logs).
    *   Orchestrates LLM API calls and parses decision logic.
    *   Handles local file system access for Markdown persistence.
2.  **Web-Based Frontend (The Orchestrator)**:
    *   A high-end UI built with modern HTML/CSS/JS.
    *   Utilizes canvas/graph libraries for intuitive workflow mapping.
    *   Real-time status updates via JS bridges to the Go runtime.

---

## 🛠️ Project Structure

```bash
.
├── assets/             # Brand assets and mockups
├── docs/               # Markdown documents for tasks/execution logs
├── internal/           # Private Go packages
│   ├── engine/         # Graph orchestration and decision logic
│   ├── agent/          # Persona and LLM provider integration
│   └── ui/             # WebView glue code
├── ui/                 # Frontend source
│   ├── index.html      # Main application entry
│   ├── css/            # Style definitions
│   └── js/             # Graph and UI logic
├── main.go             # Application entry point
├── workflow.go         # Data models and YAML persistence
└── README.md
```

---

## 🚦 Getting Started

### Prerequisites
*   [Go 1.21+](https://go.dev/dl/)
*   **Windows**: WebView2 runtime (Standard on Win 10/11).
*   **macOS/Linux**: See [webview_go requirements](https://github.com/webview/webview_go).

### Installation
1.  **Clone the repository**:
    ```bash
    git clone https://github.com/shanecandoit/agent_swimlanes_go_webview.git
    cd agent_swimlanes_go_webview
    ```
2.  **Run the application**:
    ```bash
    go run .
    ```

---

## 📄 License

Distributed under the GNU Affero General Public License v3. See `LICENSE` for more information.

---

## 🗺️ Roadmap

- [ ] Interactive Canvas UI (Diagram-style editing).
- [ ] DECISION node evaluation logic (LLM integration).
- [ ] Support for local models (Ollama/LM Studio).
- [ ] Export full workflow audit logs to PDF/JSON.

---

*Transforming static boards into dynamic agentic logic engines.* 🌊
