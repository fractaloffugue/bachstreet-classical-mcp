package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"bachstreet-classical-mcp/client"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func GetScoreLinksTool(imslpClient *client.Client) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.Tool{
		Name:        "get_score_links",
		Description: "Get download links for PDF sheet music scores of a work. Returns all available editions and arrangements with direct download URLs.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"page_title": map[string]interface{}{
					"type":        "string",
					"description": "The IMSLP page title (e.g., 'Prelude and Fugue in C major, BWV 846 (Bach, Johann Sebastian)'). You can get this from search_work results.",
				},
			},
			Required: []string{"page_title"},
		},
		},
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			var args struct {
				PageTitle string `json:"page_title"`
			}

			argsBytes, _ := json.Marshal(request.Params.Arguments)
			if err := json.Unmarshal(argsBytes, &args); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
			}

			scores, err := imslpClient.GetScoreLinks(args.PageTitle)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get score links: %v", err)), nil
			}

			if len(scores) == 0 {
				return mcp.NewToolResultText("No PDF scores found for this work. The work page may not have uploaded scores yet."), nil
			}

			resultJSON, err := json.MarshalIndent(scores, "", "  ")
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to format results: %v", err)), nil
			}

			return mcp.NewToolResultText(string(resultJSON)), nil
		},
	}
}
