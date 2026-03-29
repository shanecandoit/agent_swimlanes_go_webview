# 📡 Agent Swimlanes REST API Design

This document outlines the RESTful API surface for the **Agent Swimlanes Workflow Logic Engine**. It has been updated to support a graph-based data model with polymorphic nodes (Tasks/Decisions) and directed edges.

## Base URL
The local orchestration engine runs on: `http://localhost:8080`

---

## 1. Graph Structure Management

### `GET /api/workflows/:id/graph`
Returns a JSON representation of the entire workflow graph (nodes and edges). This is designed for compatibility with canvas libraries like `react-flow` or `d3-graphviz`.

### `POST /api/workflows/:id/nodes`
Adds a new node to the workflow.
*   **Payload**:
    ```json
    {
      "type": "TASK" | "DECISION",
      "swimlane_id": "UUID",
      "label": "String",
      "system_prompt": "String"
    }
    ```

### `POST /api/workflows/:id/edges`
Connects two nodes in the graph.
*   **Payload**:
    ```json
    {
      "from_node_id": "UUID",
      "to_node_id": "UUID",
      "label": "String" (Required for edges from DECISION nodes)
    }
    ```

---

## 2. Decision Engine Logic

### `POST /api/nodes/:id/evaluate`
**The Decision Logic Endpoint.** Triggers an agent to evaluate the current context and choose the next path in the workflow.

*   **Payload**:
    ```json
    {
      "current_input_markdown": "String"
    }
    ```
*   **Logic**:
    1.  The agent retrieves the node's `system_prompt` and the labels of all outgoing `edges`.
    2.  The agent evaluates the input evidence.
    3.  The agent selects the best-fit edge label and provides a rationale.
    4.  An `ExecutionLog` is saved on the backend.
*   **Response**: Returns the `target_node_id` for the next step in the flow.

---

## 3. Execution & Auditing

### `GET /api/execution/logs/:node_id`
Retrieves the most recent reasoning and decision outcome for a specific node.

**Response:**
```json
{
  "id": "UUID",
  "node_id": "UUID",
  "reasoning": "The input content focuses on security, so the 'Security Review' path was selected.",
  "selected_edge_id": "UUID"
}
```

---

## 4. Agent Personas

### `GET /api/agents`
Lists all active agent personas and their specialized instructions.

---

## 5. System & Health

### `GET /api/health`
Returns system status and heartbeat: `Agent Swimlanes Engine: Online`
