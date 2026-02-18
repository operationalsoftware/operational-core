package handler

import (
	"app/internal/views/imagetotextview"
	"app/pkg/reqcontext"
	"net/http"
)

type ImageToTextHandler struct {
}

func NewImageToTextHandler() *ImageToTextHandler {
	return &ImageToTextHandler{}
}

func (h *ImageToTextHandler) ImageToTextResolvePage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	_ = imagetotextview.ImageToTextResolvePage(&imagetotextview.ImageToTextResolvePageProps{
		Ctx: ctx,
	}).Render(w)
}
