package components

import (
	"app/internal/model"
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type CommentsThreadProps struct {
	Title           g.Node
	Comments        []model.Comment
	CommentThreadID int
	HMACEnvelope    string
}

func CommentsThread(p *CommentsThreadProps) g.Node {

	comments := []g.Node{}
	for _, comment := range p.Comments {

		var nonImageAttachments []g.Node
		var imageURLs []string

		for _, attachment := range comment.Attachments {
			isImage := strings.HasPrefix(attachment.ContentType, "image/")

			if isImage {
				imageURLs = append(imageURLs, attachment.DownloadURL)
			} else {
				link := h.Div(
					h.Class("attachment"),
					Icon(&IconProps{
						Identifier: "open-in-new",
					}),
					h.A(
						h.Class("attachment-link"),
						h.Href(attachment.DownloadURL),
						h.Target("_blank"),
						g.Text(attachment.Filename),
					),
				)
				nonImageAttachments = append(nonImageAttachments, link)
			}

		}
		attachments := []g.Node{
			g.Group(nonImageAttachments),
			Gallery(imageURLs),
		}

		commentNode := h.Div(
			h.Class("comment"),

			h.Div(
				h.Class("comment-details"),

				h.P(
					g.Text(comment.Comment),
				),

				h.Div(
					h.Class("date"),

					h.Strong(
						g.Text(comment.CommentedAt.Format("2006-01-02 15:04:05")),
					),

					h.Div(
						h.Class("username"),
						g.Text(comment.CommentedByUsername),
					),
				),
			),

			h.Div(
				h.Class("attachments"),
				g.Group(attachments),
			),
		)
		comments = append(comments, commentNode)
	}

	return h.Div(
		h.Class("comment-section"),

		g.If(p.Title != nil, p.Title),
		g.If(p.Title == nil,
			h.H3(
				g.Text("Comments"),
			),
		),

		g.If(len(p.Comments) == 0,
			h.Div(
				h.Class("no-entries"),
				g.Text("No comments yet."),
			),
		),

		h.Div(
			h.Class("comments"),

			g.Group(comments),
		),

		h.Form(
			h.Class("comment-form"),
			h.Data("thread-id", fmt.Sprintf("%d", p.CommentThreadID)),
			h.Data("hmac-envelope", p.HMACEnvelope),
			h.Name("comment-form"),
			h.Method("POST"),
			h.EncType("multipart/form-data"),
			g.Attr("onsubmit", "submitComment(event)"),

			h.Div(
				h.Class("comment-box"),

				h.Textarea(
					h.Class("new-comment"),
					h.Name("Comment"),

					h.Placeholder("Enter Comment"),
				),

				h.Div(
					h.Class("file-upload-wrapper"),

					h.Div(
						h.Class("files"),

						h.Label(
							h.Class("file-input-label button small"),

							h.Input(
								h.Class("file-input"),
								h.Name("Files"),
								h.Type("file"),
								h.Multiple(),
								h.Accept("image/*,video/*,application/pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.odt,.ods,.odp"),
								g.Attr("data-max-files", "10"),
							),

							Icon(&IconProps{
								Identifier: "paperclip-plus",
							}),
							g.Text("Attach files"),
						),

						h.Div(
							h.ID("selected-files"),
							h.Class("selected-files"),
						),
					),
				),

				h.Div(
					h.Class("submit-btn"),
					Button(&ButtonProps{
						Classes: c.Classes{
							"add-comment-btn": true,
						},
						ButtonType: "primary",
						Loading:    true,
					},
						g.Attr("type", "submit"),
						Icon(&IconProps{
							Identifier: "comment-text-outline",
						}),
						g.Text(" Comment"),
					),
				),
			),
		),

		InlineStyle("/internal/components/CommentsThread.css"),
		InlineScript("/internal/components/CommentsThread.js"),
	)
}
