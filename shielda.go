// Package shielda provides the official Go client for Shielda Captcha verification.
package shielda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	DefaultVerifyURL = "https://shielda.ir/api/v1/siteverify"
	DefaultWidgetURL = "https://shielda.ir/widget/api.js"
	DefaultTimeout   = 10 * time.Second
)

// Client handles interaction with the Shielda Captcha API.
type Client struct {
	siteKey    string
	secretKey  string
	verifyURL  string
	widgetURL  string
	httpClient *http.Client
}

// Option configures a Shielda Client.
type Option func(*Client)

// WithVerifyURL overrides the default verify URL.
func WithVerifyURL(url string) Option {
	return func(c *Client) {
		c.verifyURL = url
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets custom request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// New creates a new Shielda Client instance.
func New(siteKey, secretKey string, opts ...Option) *Client {
	c := &Client{
		siteKey:   siteKey,
		secretKey: secretKey,
		verifyURL: DefaultVerifyURL,
		widgetURL: DefaultWidgetURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// VerifyResponse represents the JSON response returned by Shielda siteverify API.
type VerifyResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts,omitempty"`
	Hostname    string   `json:"hostname,omitempty"`
	ErrorCodes  []string `json:"error_codes,omitempty"`
}

type verifyPayload struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	RemoteIP string `json:"remoteip,omitempty"`
}

// Verify checks if the provided captcha response token is valid.
func (c *Client) Verify(ctx context.Context, responseToken string, remoteIP ...string) (bool, error) {
	resp, err := c.VerifyDetails(ctx, responseToken, remoteIP...)
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

// VerifyDetails checks the token and returns the complete API response.
func (c *Client) VerifyDetails(ctx context.Context, responseToken string, remoteIP ...string) (*VerifyResponse, error) {
	if responseToken == "" {
		return &VerifyResponse{
			Success:    false,
			ErrorCodes: []string{"missing-input-response"},
		}, errors.New("shielda: missing response token")
	}

	if c.secretKey == "" {
		return &VerifyResponse{
			Success:    false,
			ErrorCodes: []string{"missing-secret-key"},
		}, errors.New("shielda: missing secret key")
	}

	payload := verifyPayload{
		Secret:   c.secretKey,
		Response: responseToken,
	}
	if len(remoteIP) > 0 && remoteIP[0] != "" {
		payload.RemoteIP = remoteIP[0]
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("shielda: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.verifyURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("shielda: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Shielda-Go/1.0.0")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("shielda: http request failed: %w", err)
	}
	defer res.Body.Close()

	var verifyResp VerifyResponse
	if err := json.NewDecoder(res.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("shielda: failed to decode response: %w", err)
	}

	return &verifyResp, nil
}

// GetScriptTag returns the HTML script tag for the widget.
func (c *Client) GetScriptTag() string {
	return fmt.Sprintf(`<script src="%s" async defer></script>`, c.widgetURL)
}

// GetWidgetTag returns the HTML widget container tag.
func (c *Client) GetWidgetTag(theme, lang string) string {
	if theme == "" {
		theme = "auto"
	}
	if lang == "" {
		lang = "fa"
	}
	return fmt.Sprintf(`<div class="darban-captcha" data-sitekey="%s" data-theme="%s" data-lang="%s"></div>`, c.siteKey, theme, lang)
}
