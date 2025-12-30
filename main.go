package main

import (
	"log"

	"bachstreet-classical-mcp/client"
	"bachstreet-classical-mcp/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	imslpClient := client.NewClient()
	s := server.NewMCPServer(
		"Classical Music MCP Server",
		"0.1.0",
		server.WithRecovery(),
	)

	s.AddTools(
		tools.SearchWorkTool(imslpClient),
		tools.GetWorkDetailsTool(imslpClient),
		tools.GetScoreLinksTool(imslpClient),
	)

	// serve via stdio for MCP protocol
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}