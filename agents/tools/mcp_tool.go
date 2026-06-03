package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// MCPTool is the agent's single interface to the MCP server.
// The LLM calls this tool with an "action" field to either:
// - List available tools/resources: {"action": "list_tools"}
// - Execute a specific tool: {"action": "query_documents", "collection": "...", ...}
type MCPTool struct {
	client *MCPClient
}

func NewMCPTool(client *MCPClient) *MCPTool {
	return &MCPTool{client: client}
}

func (m *MCPTool) Name() string {
	return "mcp"
}

func (m *MCPTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "mcp",
		Description: "Access MongoDB database operations via the Model Context Protocol. Use this tool to discover available resources (action: list_tools) or execute operations (action: <operation_name>, then include operation-specific parameters).",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"action": map[string]any{
					"type":        "string",
					"description": "Either 'list_tools' to discover available operations, or the name of an MCP tool to execute (e.g., 'query_documents', 'insert_document', 'update_document', 'delete_document', 'list_collections')",
				},
			},
			"required": []string{"action"},
		},
	}
}

// Execute handles the MCP tool invocation.
// Input is JSON with "action" field (and operation-specific parameters for non-list calls).
func (m *MCPTool) Execute(input string) (string, error) {
	var req map[string]any
	if err := json.Unmarshal([]byte(input), &req); err != nil {
		return "", fmt.Errorf("invalid MCP tool input: %w", err)
	}

	action, ok := req["action"].(string)
	if !ok {
		return "", fmt.Errorf("action field required and must be a string")
	}

	ctx := context.Background()

	// Handle list_tools action — LLM discovers what's available
	if action == "list_tools" {
		return m.listTools(ctx)
	}

	// Handle specific tool invocation
	return m.callTool(ctx, action, req)
}

// listTools returns a formatted list of available MCP tools and their schemas.
func (m *MCPTool) listTools(ctx context.Context) (string, error) {
	tools, err := m.client.ListToolsRaw(ctx)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("Available MCP tools:\n\n")
	for _, t := range tools {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", t.Name, t.Description))
		// Convert InputSchema to map[string]any for formatting
		data, _ := json.Marshal(t.InputSchema)
		var schemaMap map[string]any
		json.Unmarshal(data, &schemaMap)
		if props := FormatInputSchema(schemaMap); props != "" {
			sb.WriteString(fmt.Sprintf("  Required: %s\n", props))
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

// callTool invokes a specific MCP tool with the given action and parameters.
func (m *MCPTool) callTool(ctx context.Context, action string, req map[string]any) (string, error) {
	// Extract tool parameters (everything except "action")
	params := make(map[string]any)
	for k, v := range req {
		if k != "action" {
			params[k] = v
		}
	}

	// If no params, use empty object
	if len(params) == 0 {
		params = make(map[string]any)
	}

	// Serialize params to JSON for MCP call
	paramsJSON, _ := json.Marshal(params)

	// Call the MCP server tool
	result, err := m.client.Call(ctx, action, string(paramsJSON))
	if err != nil {
		return "", err
	}

	return result, nil
}
