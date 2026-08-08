package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	defaultAuthBaseURL = "https://auth.greenixdb.com"
	defaultAPIBaseURL  = "https://api.greenixdb.com"
)

// AuthBaseURL is the browser-facing auth host (override with GREENIX_AUTH_URL).
func AuthBaseURL() string {
	if v := os.Getenv("GREENIX_AUTH_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return defaultAuthBaseURL
}

// APIBaseURL is the backend host the CLI polls (override with GREENIX_API_URL).
func APIBaseURL() string {
	if v := os.Getenv("GREENIX_API_URL"); v != "" {
		return strings.TrimRight(v, "/")
	}
	return defaultAPIBaseURL
}

var httpClient = &http.Client{Timeout: 20 * time.Second}

// Session is a pending CLI login attempt.
type Session struct {
	ID              string
	Secret          string
	VerificationURL string
}

// SessionStatus is the poll response from the backend.
type SessionStatus struct {
	Status      string `json:"status"` // pending | authorized | denied | expired
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	User        struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Plan  string `json:"plan"`
	} `json:"user"`
	Error string `json:"error"`
}

func randomToken(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// NewSession creates a login session and the browser URL for it.
// The session id is public; the secret never leaves the machine — only its
// SHA-256 challenge is sent through the browser.
func NewSession(cliVersion string) (*Session, error) {
	id, err := randomToken(18)
	if err != nil {
		return nil, err
	}
	secret, err := randomToken(32)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256([]byte(secret))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])

	hostname, _ := os.Hostname()

	q := url.Values{}
	q.Set("session", id)
	q.Set("challenge", challenge)
	q.Set("client", "greenix-cli")
	q.Set("cli_version", cliVersion)
	q.Set("os", runtime.GOOS)
	q.Set("arch", runtime.GOARCH)
	if hostname != "" {
		q.Set("device", hostname)
	}

	return &Session{
		ID:              id,
		Secret:          secret,
		VerificationURL: fmt.Sprintf("%s/cli?%s", AuthBaseURL(), q.Encode()),
	}, nil
}

// Poll asks the backend once for the current state of the session.
func (s *Session) Poll(ctx context.Context) (*SessionStatus, error) {
	endpoint := fmt.Sprintf("%s/cli/auth/session/%s", APIBaseURL(), url.PathEscape(s.ID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Greenix-Session-Secret", s.Secret)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusAccepted {
		return &SessionStatus{Status: "pending"}, nil
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("auth server returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var status SessionStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("unexpected response from auth server: %w", err)
	}
	if status.Status == "" {
		if status.AccessToken != "" {
			status.Status = "authorized"
		} else {
			status.Status = "pending"
		}
	}
	return &status, nil
}

// WaitForApproval polls until the session resolves, the context is cancelled,
// or the deadline passes.
func (s *Session) WaitForApproval(ctx context.Context, interval, timeout time.Duration) (*SessionStatus, error) {
	deadline := time.Now().Add(timeout)
	for {
		status, err := s.Poll(ctx)
		if err == nil && status.Status != "pending" {
			return status, nil
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out waiting for browser sign-in")
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}
	}
}

// User is the identity behind an access token.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Plan  string `json:"plan"`
}

// Me verifies the token server-side and returns the current user.
func Me(ctx context.Context, token string) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APIBaseURL()+"/cli/auth/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrNotLoggedIn
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("auth server returned %s", resp.Status)
	}

	var payload struct {
		User
		Nested *User `json:"user"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unexpected response from auth server: %w", err)
	}
	if payload.Nested != nil && payload.Nested.ID != "" {
		return payload.Nested, nil
	}
	return &payload.User, nil
}

// Revoke invalidates the token server-side. Errors are advisory: the caller
// still clears local credentials.
func Revoke(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, APIBaseURL()+"/cli/auth/logout", strings.NewReader("{}"))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<16))

	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("auth server returned %s", resp.Status)
	}
	return nil
}

// OpenBrowser tries to launch the default browser for the given URL.
func OpenBrowser(target string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start()
	case "darwin":
		return exec.Command("open", target).Start()
	default:
		if _, err := exec.LookPath("xdg-open"); err == nil {
			return exec.Command("xdg-open", target).Start()
		}
		return fmt.Errorf("no browser opener found")
	}
}
