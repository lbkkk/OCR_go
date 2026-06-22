// Package llm provides a minimal client for OpenAI-compatible chat completion
// APIs, including optional vision (image) input.
package llm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultTimeout = 60 * time.Second

// Client talks to an OpenAI-compatible chat completions endpoint.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	model      string
}

// Option customizes a Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.httpClient = h }
}

// WithTimeout sets the request timeout on the default HTTP client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// NewClient creates a Client for the given base URL (e.g. "https://host/v1"),
// API key and model name.
func NewClient(baseURL, apiKey, model string, opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		model:      model,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Complete sends the messages to the model and returns the assistant's reply.
func (c *Client) Complete(ctx context.Context, messages []Message) (string, error) {
	wireMessages := make([]wireMessage, 0, len(messages))
	for _, m := range messages {
		wireMessages = append(wireMessages, toWireMessage(m))
	}

	payload, err := json.Marshal(chatRequest{Model: c.model, Messages: wireMessages})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call llm: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed chatResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("llm: response contained no choices")
	}
	return parsed.Choices[0].Message.Content, nil
}

// --- wire types: exact JSON shapes for the OpenAI-compatible API ---

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []wireMessage `json:"messages"`
}

type wireMessage struct {
	Role string `json:"role"`
	// Content is either a string (text-only) or a []contentPart (vision).
	Content any `json:"content"`
}

type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// toWireMessage converts a public Message into its JSON wire representation.
func toWireMessage(m Message) wireMessage {
	if len(m.Images) == 0 {
		return wireMessage{Role: string(m.Role), Content: m.Text}
	}

	parts := make([]contentPart, 0, len(m.Images)+1)
	if m.Text != "" {
		parts = append(parts, contentPart{Type: "text", Text: m.Text})
	}
	for _, img := range m.Images {
		parts = append(parts, contentPart{
			Type:     "image_url",
			ImageURL: &imageURL{URL: dataURL(img)},
		})
	}
	return wireMessage{Role: string(m.Role), Content: parts}
}

// dataURL encodes an image as a base64 data URL usable by vision models.
func dataURL(img Image) string {
	return fmt.Sprintf("data:image/%s;base64,%s", img.Format, base64.StdEncoding.EncodeToString(img.Data))
}
