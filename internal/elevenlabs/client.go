package elevenlabs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.elevenlabs.io"
	defaultVoiceID = "21m00Tcm4TlvDq8ikWAM" // Rachel
	defaultModel   = "eleven_flash_v2_5"
)

// Client calls the ElevenLabs text-to-speech API.
type Client struct {
	apiKey  string
	baseURL string
	voiceID string
	model   string
	http    *http.Client
}

// Option configures a Client.
type Option func(*Client)

func WithBaseURL(url string) Option  { return func(c *Client) { c.baseURL = url } }
func WithVoiceID(id string) Option   { return func(c *Client) { c.voiceID = id } }
func WithModel(model string) Option  { return func(c *Client) { c.model = model } }

// NewClient creates an ElevenLabs API client.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		voiceID: defaultVoiceID,
		model:   defaultModel,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type synthesizeRequest struct {
	Text    string `json:"text"`
	ModelID string `json:"model_id"`
}

// Synthesize sends text to ElevenLabs and returns the MP3 audio bytes.
func (c *Client) Synthesize(text string) ([]byte, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("ELEVENLABS_API_KEY not set")
	}

	body, err := json.Marshal(synthesizeRequest{Text: text, ModelID: c.model})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/text-to-speech/%s", c.baseURL, c.voiceID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(errBody))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return data, nil
}
