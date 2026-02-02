package printnode

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckStatusNotConfigured(t *testing.T) {
	client := NewClient("")

	status, err := client.CheckStatus(context.Background())
	if err != nil {
		t.Fatalf("expected no error when api key missing, got %v", err)
	}

	if status.Configured {
		t.Fatalf("expected status.Configured to be false when no api key is set")
	}

	if !strings.Contains(status.Message, "not configured") {
		t.Fatalf("expected message to mention not configured, got %q", status.Message)
	}
}

func TestCheckStatusUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient("bad-key")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	status, err := client.CheckStatus(context.Background())
	if err != nil {
		t.Fatalf("expected no error on unauthorized, got %v", err)
	}

	if !status.Configured {
		t.Fatalf("expected status.Configured to be true when api key is provided")
	}

	if status.Reachable {
		t.Fatalf("expected status.Reachable to be false on unauthorized response")
	}

	if !strings.Contains(strings.ToLower(status.Message), "unauthorized") {
		t.Fatalf("expected unauthorized message, got %q", status.Message)
	}
}

func TestCheckStatusSuccess(t *testing.T) {
	apiKey := "test-api-key"
	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(apiKey+":"))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != expectedAuth {
			t.Fatalf("unexpected authorization header: %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"firstname":"Jane","lastname":"Doe","email":"jane@example.com"}`))
	}))
	defer server.Close()

	client := NewClient(apiKey)
	client.baseURL = server.URL
	client.httpClient = server.Client()

	status, err := client.CheckStatus(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !status.Configured || !status.Reachable {
		t.Fatalf("expected configured and reachable, got configured=%v reachable=%v", status.Configured, status.Reachable)
	}

	if status.AccountEmail != "jane@example.com" {
		t.Fatalf("unexpected account email: %s", status.AccountEmail)
	}

	if status.AccountName != "Jane Doe" {
		t.Fatalf("unexpected account name: %s", status.AccountName)
	}

	if !strings.Contains(status.Message, "Connected") {
		t.Fatalf("expected success message, got %q", status.Message)
	}
}

func TestCheckStatusServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient("some-key")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	status, err := client.CheckStatus(context.Background())
	if err == nil {
		t.Fatalf("expected error on server error response")
	}

	if !strings.Contains(status.Message, "returned an error") {
		t.Fatalf("expected error message, got %q", status.Message)
	}

	if status.Reachable {
		t.Fatalf("expected reachable to be false")
	}
}

func TestListPrintersSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/printers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "HP LaserJet",
				"description": "Main floor",
				"default": true,
				"state": "online",
				"computer": {"id": 99, "name": "Server-1", "hostname": "server-1"}
			},
			{
				"id": 2,
				"name": "Label Printer",
				"description": "",
				"default": false,
				"state": "offline",
				"computer": {"id": 100, "name": "", "hostname": "host-100"}
			}
		]`))
	}))
	defer server.Close()

	client := NewClient("key")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	printers, err := client.ListPrinters(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(printers) != 2 {
		t.Fatalf("expected 2 printers, got %d", len(printers))
	}

	first := printers[0]
	if first.ID != 1 || first.Name != "HP LaserJet" || first.ComputerName != "Server-1" || !first.Default || first.State != "online" {
		t.Fatalf("unexpected first printer: %+v", first)
	}

	second := printers[1]
	if second.ComputerName != "host-100" || second.Default {
		t.Fatalf("unexpected second printer: %+v", second)
	}
}

func TestListPrintersUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient("bad-key")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	_, err := client.ListPrinters(context.Background())
	if err == nil || !strings.Contains(err.Error(), "unauthorized") {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

func TestListPrintersNotConfigured(t *testing.T) {
	client := NewClient("")
	_, err := client.ListPrinters(context.Background())
	if err == nil || !strings.Contains(err.Error(), "not configured") {
		t.Fatalf("expected configuration error, got %v", err)
	}
}

func TestSubmitPrintJobUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient("bad")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	_, err := client.SubmitPrintJob(context.Background(), 1, "Test", "pdf_base64", "data", nil)
	if err == nil || !strings.Contains(err.Error(), "unauthorized") {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

func TestSubmitPrintJobSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":123}`))
	}))
	defer server.Close()

	client := NewClient("good")
	client.baseURL = server.URL
	client.httpClient = server.Client()

	jobID, err := client.SubmitPrintJob(context.Background(), 1, "Test", "pdf_base64", "data", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if jobID != 123 {
		t.Fatalf("unexpected job id: %d", jobID)
	}
}
