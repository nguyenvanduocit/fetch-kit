package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nguyenvanduocit/fetch-kit/services"
	"github.com/nguyenvanduocit/fetch-kit/util"
)

func RegisterJinaTool(s *server.MCPServer) {
	tool := mcp.NewTool("get_web_content",
		mcp.WithDescription("Fetches content from a given HTTP/HTTPS URL. This tool allows you to retrieve text content from web pages, APIs, or any accessible HTTP endpoints. Returns the raw content as text."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The complete HTTP/HTTPS URL to fetch content from (e.g., https://example.com)"),
		),
	)

	s.AddTool(tool, util.ErrorGuard(fetchHandler))
}

func fetchHandler(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	url, ok := arguments["url"].(string)
	if !ok {
		return mcp.NewToolResultError("url must be a string"), nil
	}

	// Construct Jina proxy URL
	jinaProxyURL := fmt.Sprintf("https://r.jina.ai/%s", url)

	req, err := http.NewRequest("GET", jinaProxyURL, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %s", err)), nil
	}

	// Get Jina API key from environment
	jinaAPIKey := os.Getenv("JINA_API_KEY")
	if jinaAPIKey == "" {
		return mcp.NewToolResultError("JINA_API_KEY environment variable is not set"), nil
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jinaAPIKey))

	resp, err := services.DefaultHttpClient().Do(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to fetch URL: %s", err)), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read response body: %s", err)), nil
	}

	return mcp.NewToolResultText(string(body)), nil
}
