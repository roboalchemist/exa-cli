package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var Version = "dev"

// Client is the Exa API client.
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	debug      func(string, ...interface{})
}

// NewClient creates a new API client.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// SetDebug enables debug logging.
func (c *Client) SetDebug(fn func(string, ...interface{})) {
	c.debug = fn
}

func (c *Client) debugLog(format string, args ...interface{}) {
	if c.debug != nil {
		c.debug(format, args...)
	}
}

func (c *Client) doJSON(ctx context.Context, method, endpoint string, body, result interface{}) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimLeft(endpoint, "/"))

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		c.debugLog("%s %s body=%s", method, url, string(jsonBody))
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		c.debugLog("%s %s", method, url)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("User-Agent", "exa-cli/"+Version)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	c.debugLog("Response status: %d", resp.StatusCode)
	c.debugLog("Response body: %s", truncate(string(respBody), 2000))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, truncate(string(respBody), 500))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
	}

	return nil
}

// Search performs a web search.
func (c *Client) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	var resp SearchResponse
	if err := c.doJSON(ctx, http.MethodPost, "/search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetContents retrieves page contents by URL or ID.
func (c *Client) GetContents(ctx context.Context, req *ContentsRequest) (*ContentsResponse, error) {
	var resp ContentsResponse
	if err := c.doJSON(ctx, http.MethodPost, "/contents", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindSimilar finds pages similar to a URL.
func (c *Client) FindSimilar(ctx context.Context, req *FindSimilarRequest) (*FindSimilarResponse, error) {
	var resp FindSimilarResponse
	if err := c.doJSON(ctx, http.MethodPost, "/findSimilar", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Answer gets an LLM-generated answer with citations.
func (c *Client) Answer(ctx context.Context, req *AnswerRequest) (*AnswerResponse, error) {
	var resp AnswerResponse
	if err := c.doJSON(ctx, http.MethodPost, "/answer", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AnswerStream streams an LLM-generated answer with citations.
// It calls textFn for each text chunk and doneFn with the final response.
func (c *Client) AnswerStream(ctx context.Context, req *AnswerRequest, textFn func(string), doneFn func(*AnswerResponse)) error {
	req.StreamOutput = true
	url := fmt.Sprintf("%s/answer", c.baseURL)

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	c.debugLog("POST %s body=%s", url, string(jsonBody))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("User-Agent", "exa-cli/"+Version)
	httpReq.Header.Set("Accept", "text/event-stream")

	// Use a separate client without timeout for streaming
	streamClient := &http.Client{}
	resp, err := streamClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, truncate(string(body), 500))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk AnswerStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			c.debugLog("Failed to parse SSE chunk: %s", err)
			continue
		}

		if chunk.Text != "" && textFn != nil {
			textFn(chunk.Text)
		}
		// Final chunk with citations
		if chunk.Citations != nil && doneFn != nil {
			doneFn(&AnswerResponse{
				Answer:      chunk.Answer,
				Citations:   chunk.Citations,
				CostDollars: chunk.CostDollars,
			})
		}
	}

	return scanner.Err()
}

// GetContext retrieves code context.
func (c *Client) GetContext(ctx context.Context, req *ContextRequest) (*ContextResponse, error) {
	var resp ContextResponse
	if err := c.doJSON(ctx, http.MethodPost, "/context", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAPIKeys lists API keys for the team.
func (c *Client) ListAPIKeys(ctx context.Context) (*APIKeysResponse, error) {
	var resp APIKeysResponse
	if err := c.doJSON(ctx, http.MethodGet, "/team-management/api-keys", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetUsage retrieves usage data for an API key.
func (c *Client) GetUsage(ctx context.Context, keyID, startDate, endDate string) (*UsageResponse, error) {
	endpoint := fmt.Sprintf("/team-management/api-keys/%s/usage?startDate=%s&endDate=%s", keyID, startDate, endDate)
	var resp UsageResponse
	if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
