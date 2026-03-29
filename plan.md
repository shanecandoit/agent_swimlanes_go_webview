# 🗺️ Project Plan - Agent Swimlanes

## ✅ Phase 1: Core Engine & Graph Model (COMPLETED)
- [x] Defined polymorphic Workflow Model (Nodes, Edges, Execution Logs).
- [x] Implemented Graph Persistence (YAML loading/saving).
- [x] Created "Modern Website Evaluation" scenario (Architect, Designer, Writer, Strategist).
- [x] Implemented "Decision Engine" API (Mock Logic).
- [x] Verified with Unit and Integration tests (`go test -v`).

## ✅ Phase 2: Frontend Migration & Visualization (COMPLETED)
- [x] **Data Sync**: Update `index.html` to consume `/api/workflow/graph` with absolute URLs (handling CORS 'null' origin).
- [x] **Node Layout**: Group nodes by `swimlane_id` in the UI.
- [x] **UX Polish**: Added Lane Borders and Agent Persona Tooltips.
- [x] **Visual Edges**: Implemented LeaderLine connections with scroll/resize persistence.
- [x] **Local Bundling**: Local bundling of dependencies (overcoming ERR_NETWORK_ACCESS_DENIED).
- [ ] **Rich Rendering**: Integrate a Markdown parser (e.g., `marked.js`) for task cards.

## 🧠 Phase 3: Agentic Evaluation (LLM Integration)
- [ ] **Agent Package**: Add `internal/agent` to handle LLM provider calls (OpenAI/Ollama).
- [ ] **Prompt Engineering**: Build the context aggregator (System Prompt + Evidence Markdown).
- [ ] **Decision Parsing**: Parse LLM reasoning and selected edge label from the response.

## 🚦 Phase 4: Run Mode & Orchestration
- [ ] **Active State Tracker**: Add "Current Active Node" pointer to the backend.
- [ ] **Step Traversal**: Automate the transition from one node to the next.
- [ ] **Manual Override**: UI/API for human "skip" or "reroute" actions.

---

## 🛠️ Design Philosophy
- **Local-First**: Orchestration happens on the user's machine.
- **Transparent AI**: Every decision is backed by an `ExecutionLog` (reasoning trail).
- **Modular Graph**: Highly flexible workflows that aren't restricted to linear paths.
