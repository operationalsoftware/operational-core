package printnode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.printnode.com"

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type Status struct {
	Configured   bool
	Reachable    bool
	Message      string
	AccountName  string
	AccountEmail string
}

type Printer struct {
	ID           int
	Name         string
	Description  string
	Default      bool
	State        string
	ComputerID   int
	ComputerName string
}

type PrintJobResponse struct {
	ID int `json:"id"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) CheckStatus(ctx context.Context) (Status, error) {
	status := Status{
		Configured: c.apiKey != "",
	}

	if !status.Configured {
		status.Message = "PrintNode API key is not configured"
		return status, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/whoami", nil)
	if err != nil {
		status.Message = "Failed to build PrintNode status request"
		return status, err
	}

	req.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		status.Message = "Could not reach PrintNode API"
		return status, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		status.Message = "PrintNode rejected the API key (unauthorized)"
		return status, nil
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		status.Message = "PrintNode API returned an error"
		return status, fmt.Errorf("printnode whoami failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var details whoamiResponse
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		status.Message = "Could not decode PrintNode response"
		return status, err
	}

	status.Reachable = true
	status.Message = "Connected to PrintNode"
	status.AccountEmail = details.Email
	status.AccountName = strings.TrimSpace(details.Firstname + " " + details.Lastname)

	return status, nil
}

func (c *Client) ListPrinters(ctx context.Context) ([]Printer, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("printnode api key is not configured")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/printers", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("printnode unauthorized: check api key")
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("printnode printers failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var apiPrinters []printerResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiPrinters); err != nil {
		return nil, err
	}

	printers := make([]Printer, 0, len(apiPrinters))
	for _, p := range apiPrinters {
		printers = append(printers, Printer{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			Default:      p.Default,
			State:        p.State,
			ComputerID:   p.Computer.ID,
			ComputerName: firstNonEmpty(p.Computer.Name, p.Computer.Hostname),
		})
	}

	return printers, nil
}

type whoamiResponse struct {
	Email     string `json:"email"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type printerResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
	State       string `json:"state"`
	Computer    struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Hostname string `json:"hostname"`
	} `json:"computer"`
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func (c *Client) SubmitPrintJob(ctx context.Context,
	printerID int,
	title string,
	contentType string,
	content string) (int, error) {
	if c.apiKey == "" {
		return 0, fmt.Errorf("printnode api key is not configured")
	}

	payload := map[string]any{
		"printerId":   printerID,
		"title":       title,
		"contentType": contentType,
		"content":     content,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/printjobs", strings.NewReader(string(body)))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return 0, fmt.Errorf("printnode unauthorized: check api key")
	}

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("printnode printjob failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var job PrintJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return 0, err
	}

	return job.ID, nil
}
