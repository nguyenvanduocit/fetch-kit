# Fetch Kit

A Go-based MCP (Model Control Protocol) connector that enables AI assistants like Claude to interact with web services. This tool provides a seamless interface for AI models to perform various web-related operations.

## Features

* Retrieve web content using Jina AI
* Leverage Google AI (Gemini) services
* Support for various web service integrations
* Configurable tool groups
* Interact with web services through AI models

## Installation

**Requirements:** Go 1.23.2+ (for building from source)

There are several ways to install Fetch Kit:

### Option 1: Download from GitHub Releases

1. Visit the GitHub Releases page
2. Download the binary for your platform:  
   * `fetch-kit_linux_amd64` for Linux  
   * `fetch-kit_darwin_amd64` for macOS  
   * `fetch-kit_windows_amd64.exe` for Windows
3. Make the binary executable (Linux/macOS):  
   ```
   chmod +x fetch-kit_*
   ```
4. Move it to your PATH (Linux/macOS):  
   ```
   sudo mv fetch-kit_* /usr/local/bin/fetch-kit
   ```

### Option 2: Go Install

```bash
go install github.com/nguyenvanduocit/fetch-kit@latest
```

### Option 3: Docker

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

1. Directly in the Docker run command (recommended)
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

### Step 3: Test Connection

You can test if the connection is working by asking Claude in Cursor:

```
@fetch_kit get web content from https://example.com
```

### Enable Tools

Fetch Kit supports various tool groups that can be enabled or disabled using the `ENABLE_TOOLS` environment variable. If not specified, all tools are enabled by default.

## CLI Usage

In addition to the MCP server, `fetch-kit` ships a standalone CLI binary (`fetch-cli`) for direct terminal use — no MCP client needed.

### Installation

```bash
# Build locally
just install-cli

# Or with Go
go install github.com/nguyenvanduocit/fetch-kit/cmd/fetch-cli@latest
```

### Quick Start

```bash
# Set credentials (only needed for the tools you use)
export JINA_API_KEY=your-jina-key
export GOOGLE_AI_API_KEY=your-gemini-key

# Or use a .env file
fetch-cli --env .env <command> [flags]
```

### Commands

| Command | Description |
|---------|-------------|
| `jina-fetch` | Fetch and extract web content via Jina AI |
| `gemini-fetch` | Fetch and analyze content using Gemini AI |
| `youtube` | Get YouTube video transcripts |

### Examples

```bash
# Fetch a web page
fetch-cli jina-fetch --url https://example.com

# Ask Gemini about a URL
fetch-cli gemini-fetch --question "Summarize this" --context "https://example.com"

# Get YouTube transcript
fetch-cli youtube --video-id dQw4w9WgXcQ

# JSON output
fetch-cli jina-fetch --url https://example.com --output json
```

### Flags

Every command accepts:
- `--env string` — Path to `.env` file (default: `.env`)
- `--output string` — Output format: `text` (default) or `json`

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
## Installation

### Homebrew (macOS/Linux)

```bash
brew install nguyenvanduocit/tap/fetch-kit
```
