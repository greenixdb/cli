package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Credentials is the persisted CLI session for the logged-in user.
type Credentials struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Email       string    `json:"email,omitempty"`
	Name        string    `json:"name,omitempty"`
	Plan        string    `json:"plan,omitempty"`
	LoggedInAt  time.Time `json:"logged_in_at"`
}

// ErrNotLoggedIn is returned when no credentials are stored on this machine.
var ErrNotLoggedIn = errors.New("not logged in")

// ConfigDir is the directory holding greenix CLI state (~/.greenix).
func ConfigDir() (string, error) {
	if dir := os.Getenv("GREENIX_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".greenix"), nil
}

// CredentialsPath is the absolute path of the credentials file.
func CredentialsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

// Save writes credentials to disk with owner-only permissions.
func Save(c *Credentials) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "credentials.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}
	return nil
}

// Load reads stored credentials, returning ErrNotLoggedIn when absent.
func Load() (*Credentials, error) {
	path, err := CredentialsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotLoggedIn
		}
		return nil, err
	}
	var c Credentials
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.AccessToken == "" {
		return nil, ErrNotLoggedIn
	}
	return &c, nil
}

// Clear removes stored credentials. It is a no-op when nothing is stored.
func Clear() error {
	path, err := CredentialsPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Expired reports whether the stored token is past its expiry.
func (c *Credentials) Expired() bool {
	if c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(c.ExpiresAt)
}
