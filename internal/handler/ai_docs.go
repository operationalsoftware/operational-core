package handler

import (
	"app/internal/views/aidocsview"
	"app/pkg/reqcontext"
	"encoding/json"
	"net/http"
	"strings"
)

type AIDocsHandler struct{}

func NewAIDocsHandler() *AIDocsHandler {
	return &AIDocsHandler{}
}

type aiDocsQueryRequest struct {
	Question string `json:"question"`
}

type aiDocsPageLink struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

type aiDocsQueryResponse struct {
	Answer            string                 `json:"answer"`
	Module            string                 `json:"module,omitempty"`
	ModuleStatus      string                 `json:"module_status"`
	PageLinks         []aiDocsPageLink       `json:"page_links,omitempty"`
	SupportsTestData  bool                   `json:"supports_test_data"`
	SampleData        map[string]interface{} `json:"sample_data,omitempty"`
	PublicSteps       []string               `json:"public_steps,omitempty"`
	RelatedModuleHint []string               `json:"related_modules,omitempty"`
}

func (h *AIDocsHandler) DocsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	_ = aidocsview.DocsPage(aidocsview.DocsPageProps{
		Ctx: ctx,
	}).Render(w)
}

func (h *AIDocsHandler) Query(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req aiDocsQueryRequest

	r.Body = http.MaxBytesReader(w, r.Body, 16384)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"answer":"Invalid request payload.","module_status":"not_found","supports_test_data":false}`, http.StatusBadRequest)
		return
	}

	question := strings.TrimSpace(req.Question)
	if question == "" {
		http.Error(w, `{"answer":"Please enter a question.","module_status":"not_found","supports_test_data":false}`, http.StatusBadRequest)
		return
	}

	response := aiDocsQueryResponse{
		Answer:           "AI docs backend is not wired yet. UI is active; retrieval/generation will be added step by step.",
		ModuleStatus:     "not_found",
		SupportsTestData: false,
		PublicSteps: []string{
			"Question received.",
			"RAG pipeline not connected yet.",
		},
	}

	_ = json.NewEncoder(w).Encode(response)
}
