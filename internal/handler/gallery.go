package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/galleryview"
	"app/pkg/apphmac"
	"app/pkg/appurl"
	"app/pkg/reqcontext"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
)

type GalleryHandler struct {
	fileService    service.FileService
	galleryService service.GalleryService
	hmacService    service.HMACService
}

func NewGalleryHandler(
	fileService service.FileService,
	galleryService service.GalleryService,
	hmacService service.HMACService,
) *GalleryHandler {
	return &GalleryHandler{
		fileService:    fileService,
		galleryService: galleryService,
		hmacService:    hmacService,
	}
}

func (h *GalleryHandler) GalleryPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	galleryID, _ := strconv.Atoi(r.PathValue("galleryID"))

	envStr := r.URL.Query().Get("Envelope")
	if envStr == "" {
		http.Error(w, "Missing envelope", http.StatusUnauthorized)
		return
	}
	envelope := envStr
	secretKey := h.hmacService.Secret()

	ok, err := apphmac.CheckEnvelope(envelope, secretKey, "gallery", fmt.Sprintf("%d", galleryID), "view")
	if err != nil || !ok {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}

	gallery, err := h.galleryService.GetGallery(r.Context(), galleryID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching gallery", http.StatusInternalServerError)
		return
	}

	editURL := ""
	permissions, err := apphmac.GetEnvelopePermissions(envelope, secretKey)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching permissions", http.StatusInternalServerError)
		return
	}
	if slices.Contains(permissions, "edit") {
		editURL = h.galleryService.GenerateEditTempURL(galleryID, true)
	}

	_ = galleryview.GalleryPage(&galleryview.GalleryPageProps{
		Ctx:               ctx,
		Gallery:           gallery,
		EditURL:           editURL,
		AllowedOperations: permissions,
	}).
		Render(w)
}

func (h *GalleryHandler) AddGalleryItem(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	galleryID, _ := strconv.Atoi(r.PathValue("galleryID"))

	envStr := r.URL.Query().Get("Envelope")
	if envStr == "" {
		http.Error(w, "Missing envelope", http.StatusUnauthorized)
		return
	}
	envelope := envStr
	secretKey := h.hmacService.Secret()

	ok, err := apphmac.CheckEnvelope(envelope, secretKey, "gallery", fmt.Sprintf("%d", galleryID), "edit")
	if err != nil || !ok {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}

	type addGalleryItemFormData struct {
		Filename    string
		ContentType string
		SizeBytes   int
	}
	var fd addGalleryItemFormData
	if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	var allowedContentTypes = map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"image/webp":      true,
		"video/mp4":       true,
		"application/pdf": true,
	}
	if !allowedContentTypes[fd.ContentType] {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "File type not allowed",
		})
		return
	}

	file, signedURL, err := h.galleryService.AddGalleryItem(
		r.Context(),
		&model.NewGalleryItem{
			Filename:    fd.Filename,
			ContentType: fd.ContentType,
			SizeBytes:   fd.SizeBytes,
			GalleryID:   galleryID,
		},
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding gallery item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"fileId":    file.FileID,
		"signedUrl": signedURL,
	})
}

func (h *GalleryHandler) DeleteGalleryItem(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	galleryID, _ := strconv.Atoi(r.PathValue("galleryID"))
	galleryItemID, _ := strconv.Atoi(r.PathValue("galleryItemID"))

	// Parse and validate envelope (must have edit)
	envStr := r.URL.Query().Get("Envelope")
	if envStr == "" {
		http.Error(w, "Missing envelope", http.StatusUnauthorized)
		return
	}
	envelope := envStr
	secretKey := h.hmacService.Secret()

	ok, err := apphmac.CheckEnvelope(envelope, secretKey, "gallery", fmt.Sprintf("%d", galleryID), "edit")
	if err != nil || !ok {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}

	err = h.galleryService.DeleteGalleryItem(r.Context(), galleryID, galleryItemID, ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting gallery item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GalleryHandler) SetGalleryItemPosition(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	galleryID, _ := strconv.Atoi(r.PathValue("galleryID"))
	galleryItemID, _ := strconv.Atoi(r.PathValue("galleryItemID"))

	// Parse and validate envelope (must have edit)
	envStr := r.URL.Query().Get("Envelope")
	if envStr == "" {
		http.Error(w, "Missing envelope", http.StatusUnauthorized)
		return
	}
	envelope := envStr
	secretKey := h.hmacService.Secret()

	ok, err := apphmac.CheckEnvelope(envelope, secretKey, "gallery", fmt.Sprintf("%d", galleryID), "edit")
	if err != nil || !ok {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}
	permissions, err := apphmac.GetEnvelopePermissions(envelope, secretKey)
	if err != nil || !slices.Contains(permissions, "edit") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	type galleryItemUpdate struct {
		NewPosition int
	}
	var update galleryItemUpdate
	if err := appurl.Unmarshal(r.Form, &update); err != nil {
		log.Println(err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	err = h.galleryService.SetGalleryItemPosition(
		r.Context(),
		galleryID,
		galleryItemID,
		update.NewPosition,
		ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting gallery item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GalleryHandler) EditPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	galleryID, _ := strconv.Atoi(r.PathValue("galleryID"))

	queryPath := r.URL.Query().Encode()

	// Parse envelope and require edit
	envStr := r.URL.Query().Get("Envelope")
	if envStr == "" {
		http.Error(w, "Missing envelope", http.StatusUnauthorized)
		return
	}
	envelope := envStr
	secretKey := h.hmacService.Secret()

	ok, err := apphmac.CheckEnvelope(envelope, secretKey, "gallery", fmt.Sprintf("%d", galleryID), "edit")
	if err != nil || !ok {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}

	permissions, err := apphmac.GetEnvelopePermissions(envelope, secretKey)
	if err != nil || !slices.Contains(permissions, "edit") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	gallery, err := h.galleryService.GetGallery(r.Context(), galleryID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching gallery", http.StatusInternalServerError)
		return
	}

	_ = galleryview.EditGalleryPage(&galleryview.EditGalleryPageProps{
		Ctx:            ctx,
		Gallery:        gallery,
		GalleryID:      galleryID,
		GalleryPageURL: queryPath,
	}).
		Render(w)
}
