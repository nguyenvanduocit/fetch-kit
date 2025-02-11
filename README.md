# Fetch Kit - Model Context Protocol (MCP) Server

The Model Context Protocol (MCP) implementation in Fetch Kit enables AI models to interact with web services.

## Prerequisites

- Go 1.23.2 or higher
- Google AI API key (for Gemini services)
- Jina AI API key (for web content retrieval)
- Various API keys and tokens for the services you want to use

## Installation

### Installing via Go

1. Install the server:

```bash
go install github.com/nguyenvanduocit/fetch-kit@latest
```

2. Create a `.env` file with your configuration:

```env
# Required for AI services
GOOGLE_AI_API_KEY=    # Required: API key for Google AI (Gemini) service
JINA_API_KEY=         # Required: API key for Jina AI service

# Optional configurations
ENABLE_TOOLS=         # Optional: Comma-separated list of tool groups to enable (empty = all enabled)
PROXY_URL=           # Optional: HTTP/HTTPS proxy URL if needed
```

3. Config your claude's config:

```json{claude_desktop_config.json}
{
  "mcpServers": {
    "fetch_kit": {
      "command": "fetch-kit",
      "args": ["-env", "/path/to/.env"],
    }
  }
}
```

## Enable Tools

There are a hidden variable `ENABLE_TOOLS` in the environment variable. It is a comma separated list of tools group to enable. If not set, all tools will be enabled. Leave it empty to enable all tools.

## Available Tools

### Group: gemini

#### ai_web_search

search the web by using Google AI Search. Best tool to update realtime information

### Group: jina

#### get_web_content

Fetches content from a given HTTP/HTTPS URL. This tool allows you to retrieve text content from web pages, APIs, or any accessible HTTP endpoints. Returns the raw content as text.

### Group: youtube

#### youtube_transcript

Get YouTube video transcript
