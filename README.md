# agentic-ai

A local agentic AI system built in Go. An LLM-powered ReAct agent that reasons over PDF documents, the web, and a MongoDB database — all wired together via the Model Context Protocol (MCP).

---

## Repository Structure

```
agentic-ai/
├── agents/   — ReAct agent (Ollama + PDF search + web search + MCP client)
└── mcp/      — MongoDB MCP server (exposes agentic_mcps DB as MCP tools)
```

---

## Architecture

```mermaid
graph TD
    subgraph agentic-ai
        subgraph agents
            A[ReAct Agent<br/>:8082]
            T1[search_pdf tool]
            T2[web_search tool]
            T3[MCP client tools<br/>discovered via tools/list]
        end

        subgraph mcp
            M[MCP Server<br/>:8083 / SSE]
            DB[(MongoDB<br/>agentic_mcps)]
            C1[learning_todo]
            C2[links_tracker]
            C3[job_portals]
        end
    end

    LLM[Ollama LLM] -->|ReAct loop| A
    A --> T1 --> PDF[PDF Search<br/>:8081]
    A --> T2 --> Web[Tavily Web Search]
    A --> T3 -->|MCP over SSE| M
    M -->|CRUD| DB
    DB --> C1
    DB --> C2
    DB --> C3
```

---

## How it works

```mermaid
sequenceDiagram
    participant U as User
    participant A as Agent :8082
    participant L as Ollama LLM
    participant M as MCP Server :8083
    participant DB as MongoDB

    U->>A: POST /api/agent/query
    A->>M: initialize + tools/list
    M-->>A: tool definitions

    loop ReAct steps
        A->>L: Thought prompt
        L-->>A: Action + Action Input
        alt search_pdf / web_search
            A->>A: call local tool
        else db_query / db_insert / ...
            A->>M: tools/call
            M->>DB: CRUD
            DB-->>M: result
            M-->>A: CallToolResult
        end
        A->>L: Observation
    end

    L-->>A: Final Answer
    A-->>U: JSON response
```

---

## Modules

### [`agents/`](./agents/)
ReAct loop agent powered by a local Ollama model. On startup it connects to the MCP server, calls `tools/list` to discover available DB tools, and injects their real descriptions into the system prompt.

| Tool | Source |
|---|---|
| `search_pdf` | Local PDF vector search endpoint |
| `web_search` | Tavily API |
| `list_collections`, `query_documents`, `insert_document`, `update_document`, `delete_document` | Discovered from MCP server at runtime |

### [`mcp/`](./mcp/)
Standalone MCP server exposing a MongoDB database over HTTP/SSE. Any MCP-compatible client (Claude Desktop, Cursor, or a custom agent) can connect to it — no agent-specific coupling.

SSE endpoint: `http://localhost:8083/sse`

---

## Quick Start

```bash
# 1. Start MongoDB
mongosh --eval "db.adminCommand({ping:1})"

# 2. Start the MCP server
cd mcp && go run .

# 3. Start the agent
cd agents && go run .

# 4. Query the agent
curl -X POST http://localhost:8082/api/agent/query \
  -H "Content-Type: application/json" \
  -d '{"query": "Show me all free job portals"}'
```

---

## Prerequisites

| Dependency | Purpose |
|---|---|
| Go 1.21+ | Build both modules |
| MongoDB | Backing store for MCP server |
| Ollama | Local LLM inference |
| Tavily API key | Web search fallback |
| PDF search endpoint | Optional — `http://localhost:8081/api/search` |
