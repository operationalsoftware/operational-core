package aidocsview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type DocsPageProps struct {
	Ctx reqcontext.ReqContext
}

func DocsPage(p DocsPageProps) g.Node {
	content := h.Div(
		h.Class("ai-docs-page"),

		h.Div(
			h.Class("ai-docs-header"),
			h.H1(g.Text("AI Documentation Assistant")),
			h.P(g.Text("Ask questions about modules, page links, and sample test data.")),
		),

		h.Div(
			h.Class("ai-docs-layout"),
			components.Card(
				h.Div(
					h.Class("ai-docs-side-panel"),
					h.H3(g.Text("Quick Prompts")),
					h.Div(
						h.Class("ai-docs-prompts"),
						h.Button(
							h.Type("button"),
							h.Class("button"),
							h.Data("prompt", "Does stock items module exist?"),
							g.Text("Stock items availability"),
						),
						h.Button(
							h.Type("button"),
							h.Class("button"),
							h.Data("prompt", "How do resources work and where are they?"),
							g.Text("Resource workflow"),
						),
						h.Button(
							h.Type("button"),
							h.Class("button"),
							h.Data("prompt", "Give me sample test payload for services."),
							g.Text("Service sample payload"),
						),
					),
					h.P(
						h.Class("ai-docs-note"),
						g.Text("Use prompts to ask about module availability, links, behavior, and sample data."),
					),
				),
			),
			components.Card(
				h.Div(
					h.Class("ai-docs-chat"),
					h.Div(
						h.ID("ai-docs-messages"),
						h.Class("ai-docs-messages"),
						h.Div(
							h.Class("ai-docs-message assistant"),
							h.P(g.Text("Ask me about any module and I will map pages and sample data where available.")),
						),
					),
					h.FormEl(
						h.ID("ai-docs-query-form"),
						h.Class("ai-docs-query-form"),
						h.Data("endpoint", "/ai/docs/query"),
						h.Label(
							h.For("ai-docs-question"),
							g.Text("Question"),
						),
						h.Textarea(
							h.ID("ai-docs-question"),
							h.Name("question"),
							h.Rows("3"),
							h.Placeholder("Example: How do stock items work and where can I find them?"),
						),
						h.Div(
							h.Class("ai-docs-form-actions"),
							h.Button(
								h.Type("submit"),
								h.Class("button primary"),
								g.Text("Ask"),
							),
						),
					),
				),
			),
		),
	)

	return layout.Page(layout.PageProps{
		Title:   "AI Docs",
		Content: content,
		Ctx:     p.Ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "ocr",
				Title:          "AI",
				URLPart:        "ai",
			},
			{
				Title:   "Docs",
				URLPart: "docs",
			},
		},
		AppendHead: []g.Node{
			h.StyleEl(g.Raw(aiDocsPageCSS)),
		},
		AppendBody: []g.Node{
			h.Script(g.Raw(aiDocsPageJS)),
		},
	})
}

const aiDocsPageCSS = `
.ai-docs-page {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.ai-docs-header {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.ai-docs-header h1,
.ai-docs-header p {
  margin: 0;
}

.ai-docs-header p {
  color: var(--text-color-light);
}

.ai-docs-layout {
  display: grid;
  gap: var(--spacing-lg);
  grid-template-columns: minmax(250px, 320px) 1fr;
}

.ai-docs-side-panel {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.ai-docs-side-panel h3 {
  margin: 0;
}

.ai-docs-prompts {
  display: grid;
  gap: var(--spacing-sm);
}

.ai-docs-prompts .button {
  text-align: left;
  justify-content: flex-start;
}

.ai-docs-note {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--text-color-light);
}

.ai-docs-chat {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.ai-docs-messages {
  min-height: 360px;
  max-height: 500px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  background: var(--background-color-grey);
}

.ai-docs-message {
  max-width: 85%;
  border-radius: var(--border-radius-sm);
  padding: var(--spacing-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.ai-docs-message p {
  margin: 0;
  white-space: pre-wrap;
}

.ai-docs-message.user {
  align-self: flex-end;
  background: var(--primary-color);
  color: var(--text-color-contrast);
}

.ai-docs-message.assistant {
  align-self: flex-start;
  background: var(--background-color);
  border: 1px solid var(--border-color);
}

.ai-docs-message ul {
  margin: 0;
  padding-left: var(--spacing-lg);
}

.ai-docs-message pre {
  margin: 0;
  padding: var(--spacing-md);
  border-radius: var(--border-radius-sm);
  background: var(--background-color-grey);
  border: 1px solid var(--border-color);
  overflow-x: auto;
  font-size: 12px;
}

.ai-docs-query-form {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.ai-docs-query-form label {
  font-size: var(--font-size-sm);
  color: var(--text-color-light);
}

.ai-docs-query-form textarea {
  resize: vertical;
  min-height: 96px;
}

.ai-docs-form-actions {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 980px) {
  .ai-docs-layout {
    grid-template-columns: 1fr;
  }

  .ai-docs-messages {
    min-height: 300px;
  }
}
`

const aiDocsPageJS = `
(function() {
  const form = document.getElementById("ai-docs-query-form");
  if (!form) {
    return;
  }

  const endpoint = form.dataset.endpoint || "/ai/docs/query";
  const questionInput = document.getElementById("ai-docs-question");
  const messages = document.getElementById("ai-docs-messages");
  const submitButton = form.querySelector("button[type='submit']");
  const promptButtons = document.querySelectorAll("[data-prompt]");

  function scrollToBottom() {
    messages.scrollTop = messages.scrollHeight;
  }

  function createMessage(role, text) {
    const wrapper = document.createElement("div");
    wrapper.className = "ai-docs-message " + role;

    const textNode = document.createElement("p");
    textNode.textContent = text;
    wrapper.appendChild(textNode);
    return wrapper;
  }

  function appendLinks(wrapper, links) {
    if (!Array.isArray(links) || links.length === 0) {
      return;
    }

    const list = document.createElement("ul");
    for (const link of links) {
      if (!link || !link.path) {
        continue;
      }

      const item = document.createElement("li");
      const anchor = document.createElement("a");
      anchor.href = link.path;
      anchor.textContent = link.label || link.path;
      item.appendChild(anchor);
      list.appendChild(item);
    }

    if (list.children.length > 0) {
      wrapper.appendChild(list);
    }
  }

  function appendSampleData(wrapper, sampleData) {
    if (!sampleData || typeof sampleData !== "object") {
      return;
    }

    const title = document.createElement("p");
    title.textContent = "Sample test data:";
    wrapper.appendChild(title);

    const pre = document.createElement("pre");
    pre.textContent = JSON.stringify(sampleData, null, 2);
    wrapper.appendChild(pre);
  }

  function appendPublicSteps(wrapper, steps) {
    if (!Array.isArray(steps) || steps.length === 0) {
      return;
    }

    const title = document.createElement("p");
    title.textContent = "How this answer was prepared:";
    wrapper.appendChild(title);

    const list = document.createElement("ul");
    for (const step of steps) {
      const item = document.createElement("li");
      item.textContent = step;
      list.appendChild(item);
    }
    wrapper.appendChild(list);
  }

  promptButtons.forEach((button) => {
    button.addEventListener("click", function() {
      questionInput.value = button.dataset.prompt || "";
      questionInput.focus();
    });
  });

  form.addEventListener("submit", async function(e) {
    e.preventDefault();

    const question = questionInput.value.trim();
    if (!question) {
      questionInput.focus();
      return;
    }

    messages.appendChild(createMessage("user", question));
    const loadingMessage = createMessage("assistant", "Thinking...");
    messages.appendChild(loadingMessage);
    scrollToBottom();

    questionInput.value = "";
    submitButton.disabled = true;

    try {
      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ question: question })
      });

      let payload = {};
      try {
        payload = await response.json();
      } catch (err) {
        payload = {};
      }

      const answerText = payload.answer || "I could not prepare an answer right now.";
      const assistantMessage = createMessage("assistant", answerText);
      appendLinks(assistantMessage, payload.page_links);
      appendSampleData(assistantMessage, payload.sample_data);
      appendPublicSteps(assistantMessage, payload.public_steps);

      loadingMessage.replaceWith(assistantMessage);
    } catch (err) {
      loadingMessage.replaceWith(createMessage("assistant", "Request failed. Please try again."));
    } finally {
      submitButton.disabled = false;
      scrollToBottom();
    }
  });
})();
`
