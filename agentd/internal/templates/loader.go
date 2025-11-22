package templates

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Template struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Status    int               `json:"status"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Signature string            `json:"signature"`
}

type Loader struct {
	Templates map[string]Template
	secret    string
}

func NewLoader(path, secret string) (*Loader, error) {
	l := &Loader{Templates: make(map[string]Template), secret: secret}

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %w", path, err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		full := filepath.Join(path, e.Name())
		data, err := os.ReadFile(full)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", full, err)
		}

		var t Template
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, fmt.Errorf(
				"JSON parse error in %s:\nError: %w\n\n=== FILE CONTENTS ===\n%s\n=== END ===\n",
				full,
				err,
				string(data),
			)
		}

		if !verifySignature(t.Body, t.Signature, secret) {
			// DEBUGGING: Calculate what it SHOULD be
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write([]byte(t.Body))
			expected := hex.EncodeToString(mac.Sum(nil))

			return nil, fmt.Errorf(
				"signature verification failed for %s in file %s.\nConfigured Secret: '%s'\nBody Length: %d\nExpected: %s\nGot:      %s",
				t.ID, full, secret, len(t.Body), expected, t.Signature,
			)
		}

		l.Templates[t.ID] = t
	}

	return l, nil
}

func verifySignature(body, sigHex, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	sum := mac.Sum(nil)
	expected := hex.EncodeToString(sum)
	return hmac.Equal([]byte(expected), []byte(sigHex))
}

func (l *Loader) Get(id string) (Template, bool) {
	t, ok := l.Templates[id]
	return t, ok
}
