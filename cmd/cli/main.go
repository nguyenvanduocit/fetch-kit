package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nguyenvanduocit/fetch-kit/services"
	"github.com/nguyenvanduocit/fetch-kit/tools"
	"google.golang.org/genai"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "jina-fetch":
		runJinaFetch(os.Args[2:])
	case "gemini-fetch":
		runGeminiFetch(os.Args[2:])
	case "youtube-transcript":
		runYouTubeTranscript(os.Args[2:])
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `fetch-cli - fetch content from various sources

Usage:
  fetch-cli <command> [flags]

Commands:
  jina-fetch          Fetch web content via Jina AI reader proxy
  gemini-fetch        Search the web using Google Gemini AI
  youtube-transcript  Get transcript for a YouTube video

Flags (all commands):
  --env string    Path to environment file (default ".env")
  --output string Output format: text or json (default "text")

Examples:
  fetch-cli jina-fetch --url https://example.com
  fetch-cli gemini-fetch --question "What is Go?" --context "programming language"
  fetch-cli youtube-transcript --video-id dQw4w9WgXcQ
`)
}

func loadEnv(envFile string) {
	if err := godotenv.Load(envFile); err != nil {
		// non-fatal: env vars may already be set
	}
}

// runJinaFetch fetches web content through the Jina AI reader proxy.
func runJinaFetch(args []string) {
	fs := flag.NewFlagSet("jina-fetch", flag.ExitOnError)
	envFile := fs.String("env", ".env", "Path to environment file")
	output := fs.String("output", "text", "Output format: text or json")
	url := fs.String("url", "", "The HTTP/HTTPS URL to fetch content from (required)")
	fs.Parse(args)

	loadEnv(*envFile)

	if *url == "" {
		fmt.Fprintln(os.Stderr, "error: --url is required")
		fs.Usage()
		os.Exit(1)
	}

	jinaAPIKey := os.Getenv("JINA_API_KEY")
	if jinaAPIKey == "" {
		fmt.Fprintln(os.Stderr, "error: JINA_API_KEY environment variable is not set")
		os.Exit(1)
	}

	jinaProxyURL := fmt.Sprintf("https://r.jina.ai/%s", *url)
	req, err := http.NewRequest("GET", jinaProxyURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jinaAPIKey))

	resp, err := services.DefaultHttpClient().Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to fetch URL: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to read response body: %s\n", err)
		os.Exit(1)
	}

	content := string(body)
	printOutput(*output, content, map[string]string{
		"url":     *url,
		"content": content,
	})
}

// runGeminiFetch performs a web search using Google Gemini AI.
func runGeminiFetch(args []string) {
	fs := flag.NewFlagSet("gemini-fetch", flag.ExitOnError)
	envFile := fs.String("env", ".env", "Path to environment file")
	output := fs.String("output", "text", "Output format: text or json")
	question := fs.String("question", "", "The question to ask (required)")
	questionContext := fs.String("context", "", "Context/purpose of the question (required)")
	fs.Parse(args)

	loadEnv(*envFile)

	if *question == "" {
		fmt.Fprintln(os.Stderr, "error: --question is required")
		fs.Usage()
		os.Exit(1)
	}
	if *questionContext == "" {
		fmt.Fprintln(os.Stderr, "error: --context is required")
		fs.Usage()
		os.Exit(1)
	}

	apiKey := os.Getenv("GOOGLE_AI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: GOOGLE_AI_API_KEY environment variable is not set")
		os.Exit(1)
	}

	cfg := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}
	client, err := genai.NewClient(context.Background(), cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create Gemini client: %s\n", err)
		os.Exit(1)
	}

	systemInstruction := "You are a search engine. You will search the web for the answer to the question. You will then provide the answer to the question. Always try to search the web for the answer first before providing the answer. writing style: short, concise, direct, and to the point."
	systemInstruction += "\n\nContext: " + *questionContext

	resp, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-pro-preview-03-25",
		[]*genai.Content{
			{
				Role: "user",
				Parts: []*genai.Part{
					{Text: *question},
				},
			},
		},
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Role: "system",
				Parts: []*genai.Part{
					{Text: systemInstruction},
				},
			},
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to generate content: %s\n", err)
		os.Exit(1)
	}

	if len(resp.Candidates) == 0 {
		fmt.Fprintln(os.Stderr, "error: no response from Gemini")
		os.Exit(1)
	}

	candidate := resp.Candidates[0]
	var textBuilder strings.Builder
	for _, part := range candidate.Content.Parts {
		textBuilder.WriteString(part.Text)
	}

	if candidate.CitationMetadata != nil {
		for _, citation := range candidate.CitationMetadata.Citations {
			textBuilder.WriteString("\n\nSource: ")
			textBuilder.WriteString(citation.URI)
		}
	}

	if candidate.GroundingMetadata != nil {
		textBuilder.WriteString("\n\nSources: ")
		for _, chunk := range candidate.GroundingMetadata.GroundingChunks {
			if chunk.RetrievedContext != nil {
				textBuilder.WriteString("\n")
				textBuilder.WriteString(chunk.RetrievedContext.Text)
				textBuilder.WriteString(": ")
				textBuilder.WriteString(chunk.RetrievedContext.URI)
			}
			if chunk.Web != nil {
				textBuilder.WriteString("\n")
				textBuilder.WriteString(chunk.Web.Title)
				textBuilder.WriteString(": ")
				textBuilder.WriteString(chunk.Web.URI)
			}
		}
	}

	content := textBuilder.String()
	printOutput(*output, content, map[string]string{
		"question": *question,
		"context":  *questionContext,
		"answer":   content,
	})
}

// runYouTubeTranscript fetches the transcript for a YouTube video.
func runYouTubeTranscript(args []string) {
	fs := flag.NewFlagSet("youtube-transcript", flag.ExitOnError)
	envFile := fs.String("env", ".env", "Path to environment file")
	output := fs.String("output", "text", "Output format: text or json")
	videoID := fs.String("video-id", "", "YouTube video ID or URL (required)")
	fs.Parse(args)

	loadEnv(*envFile)

	if *videoID == "" {
		fmt.Fprintln(os.Stderr, "error: --video-id is required")
		fs.Usage()
		os.Exit(1)
	}

	transcripts, videoTitle, err := tools.FetchTranscript(*videoID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to fetch transcript: %s\n", err)
		os.Exit(1)
	}

	if *output == "json" {
		type entry struct {
			Offset   float64 `json:"offset"`
			Duration float64 `json:"duration"`
			Text     string  `json:"text"`
			Lang     string  `json:"lang"`
		}
		type result struct {
			Title       string  `json:"title"`
			Transcripts []entry `json:"transcripts"`
		}
		entries := make([]entry, 0, len(transcripts))
		for _, t := range transcripts {
			entries = append(entries, entry{
				Offset:   t.Offset,
				Duration: t.Duration,
				Text:     t.Text,
				Lang:     t.Lang,
			})
		}
		data, err := json.MarshalIndent(result{Title: videoTitle, Transcripts: entries}, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to marshal JSON: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
		return
	}

	// text output
	fmt.Printf("Title: %s\n\n", videoTitle)
	for _, t := range transcripts {
		fmt.Printf("%s\n", t.Text)
	}
}

// printOutput writes content in the requested format (text or json).
func printOutput(format string, text string, data map[string]string) {
	if format == "json" {
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to marshal JSON: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(string(out))
		return
	}
	fmt.Print(text)
}
