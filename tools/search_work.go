package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"bachstreet-classical-mcp/client"
	"bachstreet-classical-mcp/models"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func SearchWorkTool(imslpClient *client.Client) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.Tool{
		Name:        "search_work",
		Description: "Search for classical music works by title, composer name, or keywords. Returns a list of matching compositions from IMSLP.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query (e.g., 'Bach Prelude C major', 'Moonlight Sonata', 'Mozart K545')",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of results to return (default: 10, max: 50)",
					"default":     10,
				},
			},
			Required: []string{"query"},
		},
		},
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args struct {
				Query string `json:"query"`
				Limit int    `json:"limit"`
			}

			argsBytes, _ := json.Marshal(request.Params.Arguments)
			if err := json.Unmarshal(argsBytes, &args); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
			}

			// Set default limit
			if args.Limit == 0 {
				args.Limit = 10
			}
			if args.Limit > 50 {
				args.Limit = 50
			}

			works, err := imslpClient.SearchWorks(args.Query, args.Limit)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
			}

			result := models.SearchResult{
				Query:      args.Query,
				TotalCount: len(works),
				Works:      works,
			}

			resultJSON, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to format results: %v", err)), nil
			}

			return mcp.NewToolResultText(string(resultJSON)), nil
		},
	}
}
