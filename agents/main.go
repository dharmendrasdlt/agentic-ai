package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/user/agentic-ai/tools"
)

var toolMgr *tools.Manager

func main() {
	port := flag.String("port", "8082", "Port to listen on")
	flag.Parse()

	staticDir := filepath.Join(".", "static")
	if _, err := os.Stat(staticDir); err != nil {
		log.Fatalf("static directory not found: %v", err)
	}

	// Pick LLM backend
	if cfg.AnthropicAPIKey != "" && cfg.AnthropicModel != "" && cfg.AnthropicHasCredits {
		llmCaller = &AnthropicCaller{apiKey: cfg.AnthropicAPIKey, model: cfg.AnthropicModel}
		log.Printf("[LLM] Backend: Anthropic (%s)", cfg.AnthropicModel)
	} else {
		llmCaller = &OllamaCaller{}
		log.Printf("[LLM] Backend: Ollama (%s)", cfg.OllamaModel)
	}

	// Initialize tool registry
	toolMgr = tools.NewManager()
	toolMgr.Register(tools.NewPDFSearchTool(cfg.SearchEndpoint, cfg.MaxRetries))
	toolMgr.Register(tools.NewWebSearchTool(cfg.TavilyAPIKey, cfg.MaxRetries))

	// Connect to MCP server (LLM will discover tools dynamically via action: "list_tools")
	mcpURL := cfg.MCPServerURL
	if mcpURL == "" {
		mcpURL = "http://localhost:8083"
	}
	mcpClient, err := tools.Connect(context.Background(), mcpURL)
	if err != nil {
		log.Printf("[MCP] WARNING: Could not connect to MCP server at %s: %v (MCP tool unavailable)", mcpURL, err)
	} else {
		toolMgr.Register(tools.NewMCPTool(mcpClient))
		log.Printf("[MCP] Registered MCP tool (LLM will discover resources dynamically)")
	}

	log.Printf("[MAIN] Tool registry: %d tools", len(toolMgr.ListTools()))

	http.HandleFunc("/api/agent/query", handleAgentQuery)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		} else {
			http.ServeFile(w, r, filepath.Join(staticDir, r.URL.Path))
		}
	})

	log.Printf("Agent server listening on http://localhost:%s", *port)
	log.Printf("Endpoint: POST http://localhost:%s/api/agent/query", *port)
	log.Printf("UI: http://localhost:%s", *port)

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
