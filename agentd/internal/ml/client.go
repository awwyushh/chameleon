package ml

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MLClient struct {
	url    string
	client *http.Client
}

type MLResponse struct {
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
}

func New(url string, timeoutSec int) *MLClient {
	return &MLClient{
		url: url,
		client: &http.Client{Timeout: time.Duration(timeoutSec) * time.Second},
	}
}

func (m *MLClient) Predict(ctx context.Context, text string) (*MLResponse, error) {
	payload := map[string]string{"text": text}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", m.url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ml request failed: %w", err)
	}
	defer resp.Body.Close()
	var out MLResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("ml decode failed: %w", err)
	}
	return &out, nil
}
