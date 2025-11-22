package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Forwarder struct {
	Url       string
	JWTSecret string
	Client    *http.Client
}

func NewForwarder(url, secret string) *Forwarder {
	return &Forwarder{Url: url, JWTSecret: secret, Client: &http.Client{Timeout: 5 * time.Second}}
}

func (f *Forwarder) SendEvent(ctx context.Context, ev map[string]interface{}) error {
	if ev == nil {
		return nil
	}
	// sign short-lived token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "agentd",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(2 * time.Minute).Unix(),
	})
	signed, err := token.SignedString([]byte(f.JWTSecret))
	if err != nil {
		return err
	}
	b, _ := json.Marshal(ev)
	req, err := http.NewRequestWithContext(ctx, "POST", f.Url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+signed)
	req.Header.Set("Content-Type", "application/json")
	resp, err := f.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("aggregator responded %d", resp.StatusCode)
	}
	return nil
}
