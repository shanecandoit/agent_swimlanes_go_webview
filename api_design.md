# 📡 Agent Swimlanes REST API Design

This document outlines the RESTful API surface for the **Agent Swimlanes Workflow Logic Engine**. 

## Base URL
The local orchestration engine runs on: `http://localhost:8080`

---

## 1. Graph Structure Management

### `GET /api/workflow/graph`
Returns a JSON representation of the current workflow graph (nodes and edges).
*   **Response**: `Workflow` object.

### `POST /api/workflow/nodes`
Adds a new node to the active workflow.
*   **Payload**: `Node` object.
*   **Response**: 201 Created.

### `POST /api/workflow/edges`
Connects two nodes in the graph.
*   **Payload**: `Edge` object.
*   **Response**: 201 Created.

---

## 2. Decision Engine & Agent Logic

### `POST /api/nodes/evaluate/:id`
Triggers an agent to evaluate the node context and choose the next path.
*   **Path**: `:id` is the `NodeID`.
*   **Payload**: `{ "current_input_markdown": "string" }`
*   **Logic**: Returns the `target_node_id` and the agent's `reasoning`.

### `GET /api/docs/:filename`
Retrieves the raw Markdown content for a document in the workflow's `doc_folder`.
*   **Response**: `text/markdown`.

---

## 3. Auditing & Personas

### `GET /api/execution/logs/:node_id`
Retrieves the most recent reasoning and decision outcome for a node.

### `GET /api/agents`
Lists all unique agent personas (swimlanes) and their representative system prompts.

---

## 4. System & Assets

### `GET /api/health`
Status check: `Agent Swimlanes Engine: Online`

### `GET /ui/*`
Serves static assets (JS, CSS) for the dashboard.
