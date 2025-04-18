# Fetch Kit

A Go-based MCP (Model Context Protocol) server that enables AI assistants like Claude to interact with web services. This tool provides a seamless interface for AI models to perform various web-related operations.

## Features

* Interact with web services through AI models
* Retrieve web content using Jina AI
* Leverage Google AI (Gemini) services
* Support for various web service integrations
* Configurable tool groups

## Installation

**Requirements:** Go 1.23.2+ (for building from source)

There are several ways to install Fetch Kit:

### Option 1: Go Install

```bash
go install github.com/nguyenvanduocit/fetch-kit@latest
```

### Option 2: Docker

#### Using Docker directly

1. Pull the pre-built image from GitHub Container Registry:
```bash
docker pull ghcr.io/nguyenvanduocit/fetch-kit:latest
```

2. Or build the Docker image locally:
```bash
docker build -t fetch-kit .
```

## Configuration

### Environment Variables

The following environment variables are used for configuration:

```
# Required for AI services
GOOGLE_AI_API_KEY=    # Required: API key for Google AI (Gemini) service
JINA_API_KEY=         # Required: API key for Jina AI service

# Optional configurations
ENABLE_TOOLS=         # Optional: Comma-separated list of tool groups to enable (empty = all enabled)
PROXY_URL=            # Optional: HTTP/HTTPS proxy URL if needed
```

You can set these:

1. Directly in the Docker run command
2. Through a `.env` file (use the `-env` flag)
3. Directly in your shell environment

## Using with Claude and Cursor

To make Fetch Kit work with Claude and Cursor, you need to add configuration to your Cursor settings.

### Step 1: Install Fetch Kit

Choose one of the installation methods above (Docker recommended).

### Step 2: Configure Cursor

1. Open Cursor
2. Go to Settings > MCP > Add MCP Server
3. Add the following configuration:

#### Option A: Using Docker (Recommended)

```json
{
  "mcpServers": {
    "fetch_kit": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "-e", "GOOGLE_AI_API_KEY=your_google_ai_key",
        "-e", "JINA_API_KEY=your_jina_api_key",
        "ghcr.io/nguyenvanduocit/fetch-kit:latest"
      ]
    }
  }
}
```

#### Option B: Using Local Binary

```json
{
  "mcpServers": {
    "fetch_kit": {
      "command": "fetch-kit",
      "args": ["-env", "/path/to/.env"]
    }
  }
}
```

### Enable Tools

Fetch Kit supports various tool groups that can be enabled or disabled using the `ENABLE_TOOLS` environment variable. If not specified, all tools are enabled by default.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

For a list of recent changes, see CHANGELOG.md.