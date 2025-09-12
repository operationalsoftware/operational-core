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
	"os"
	"strconv"
	"strings"
)

type GalleryHandler struct {
	fileService    service.FileService
	galleryService service.GalleryService
}

func NewGalleryHandler(
	fileService service.FileService,
	galleryService service.GalleryService,
) *GalleryHandler {
	return &GalleryHandler{
		fileService:    fileService,
		galleryService: galleryService,
	}
}

func (h *GalleryHandler) GalleryPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	galleryID, _ := strconv.Atoi(r.PathValue("galleryID"))

	type urlVals struct {
		HMAC              string
		AllowedOperations string
		Expires           int64
	}
	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}
	allowedOperations := strings.Split(uv.AllowedOperations, ",")

	hmacClaims := apphmac.Claims{
		Entity:            "gallery",
		EntityID:          fmt.Sprintf("%d", galleryID),
		AllowedOperations: allowedOperations,
		Expires:           uv.Expires,
	}
	hmacSecret := os.Getenv("AES_256_ENCRYPTION_KEY")

	validHMAC := apphmac.VerifyHMAC(hmacClaims, uv.HMAC, hmacSecret)
	if !validHMAC {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}

	gallery, err := h.galleryService.GetGallery(r.Context(), galleryID, ctx.User)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching gallery", http.StatusInternalServerError)
		return
	}

	editURL := h.galleryService.GenerateEditTempURL(galleryID, ctx.User.Permissions.Stock.Admin)

	_ = galleryview.GalleryPage(&galleryview.GalleryPageProps{
		Ctx:               ctx,
		Gallery:           gallery,
		EditURL:           editURL,
		AllowedOperations: allowedOperations,
	}).
		Render(w)
}

func (h *GalleryHandler) AddGalleryItem(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	galleryID, _ := strconv.Atoi(r.PathValue("galleryID"))

	type urlVals struct {
		HMAC              string
		AllowedOperations string
		Expires           int64
	}
	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}
	allowedOperations := strings.Split(uv.AllowedOperations, ",")

	hmacClaims := apphmac.Claims{
		Entity:            "gallery",
		EntityID:          fmt.Sprintf("%d", galleryID),
		AllowedOperations: allowedOperations,
		Expires:           uv.Expires,
	}
	hmacSecret := os.Getenv("AES_256_ENCRYPTION_KEY")

	validHMAC := apphmac.VerifyHMAC(hmacClaims, uv.HMAC, hmacSecret)
	if !validHMAC {
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
	position, _ := strconv.Atoi(r.PathValue("position"))

	type urlVals struct {
		HMAC              string
		AllowedOperations string
		Expires           int64
	}
	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}
	allowedOperations := strings.Split(uv.AllowedOperations, ",")

	hmacClaims := apphmac.Claims{
		Entity:            "gallery",
		EntityID:          fmt.Sprintf("%d", galleryID),
		AllowedOperations: allowedOperations,
		Expires:           uv.Expires,
	}
	hmacSecret := os.Getenv("AES_256_ENCRYPTION_KEY")

	validHMAC := apphmac.VerifyHMAC(hmacClaims, uv.HMAC, hmacSecret)
	if !validHMAC {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}

	err = h.galleryService.DeleteGalleryItem(r.Context(), galleryID, galleryItemID, position, ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting gallery item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GalleryHandler) ReorderGalleryItem(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	galleryID, _ := strconv.Atoi(r.PathValue("galleryID"))

	type urlVals struct {
		HMAC              string
		AllowedOperations string
		Expires           int64
	}
	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}
	allowedOperations := strings.Split(uv.AllowedOperations, ",")

	hmacClaims := apphmac.Claims{
		Entity:            "gallery",
		EntityID:          fmt.Sprintf("%d", galleryID),
		AllowedOperations: allowedOperations,
		Expires:           uv.Expires,
	}
	hmacSecret := os.Getenv("AES_256_ENCRYPTION_KEY")

	validHMAC := apphmac.VerifyHMAC(hmacClaims, uv.HMAC, hmacSecret)
	if !validHMAC {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}
	isEditable := false
	for _, op := range allowedOperations {
		if op == "edit" {
			isEditable = true
		}
	}
	if !isEditable {
		http.Error(w, "Resource is not editable", http.StatusForbidden)
		return
	}

	var item model.UpdateGalleryItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	err = h.galleryService.ReorderGalleryItem(
		r.Context(),
		galleryID,
		item,
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

	type urlVals struct {
		HMAC              string
		AllowedOperations string
		Expires           int64
	}
	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}
	allowedOperations := strings.Split(uv.AllowedOperations, ",")

	hmacClaims := apphmac.Claims{
		Entity:            "gallery",
		EntityID:          fmt.Sprintf("%d", galleryID),
		AllowedOperations: allowedOperations,
		Expires:           uv.Expires,
	}
	hmacSecret := os.Getenv("AES_256_ENCRYPTION_KEY")

	validHMAC := apphmac.VerifyHMAC(hmacClaims, uv.HMAC, hmacSecret)
	if !validHMAC {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}

	isEditable := false
	for _, op := range allowedOperations {
		if op == "edit" {
			isEditable = true
		}
	}
	if !isEditable {
		http.Error(w, "Resource is not editable", http.StatusForbidden)
		return
	}

	gallery, err := h.galleryService.GetGallery(r.Context(), galleryID, ctx.User)
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
