package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

type VersionInfo struct {
	Version   string `json:"version"`
	BuildCode string `json:"build_code"`
	ReleaseDate string `json:"release_date"`
	Latest    bool   `json:"latest"`
}

func NewClient(baseURL string, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey: apiKey,
	}
}

func (c *Client) CheckLatestVersion() (*VersionInfo, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/version/latest", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var versionInfo VersionInfo
	if err := json.Unmarshal(body, &versionInfo); err != nil {
		return nil, err
	}

	return &versionInfo, nil
}

func (c *Client) SendBuildInfo(project string, version string, buildCode string) error {
	data := map[string]interface{}{
		"project":    project,
		"version":    version,
		"build_code": buildCode,
		"timestamp":  time.Now().Unix(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/builds", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	return nil
}

