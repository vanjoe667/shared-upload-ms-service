package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/httpResponse"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/port"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/usecase"
)

type Handler struct {
	uc *usecase.UploadUsecase
}

func FileUploadHandler(uc *usecase.UploadUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) PresignSingle(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key         string `json:"key"`
		ContentType string `json:"content_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpResponse.WriteError(w, http.StatusBadRequest, "bad request", err)
		return
	}
	ctx := r.Context()
	resp, err := h.uc.PresignSingle(ctx, req.Key, req.ContentType)
	if err != nil {
		httpResponse.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpResponse.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) PresignDownload(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpResponse.WriteError(w, http.StatusBadRequest, "bad request", err)
		return
	}

	ctx := r.Context()
	resp, err := h.uc.PresignDownload(ctx, req.Key)
	if err != nil {
		httpResponse.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpResponse.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpResponse.WriteError(w, http.StatusBadRequest, "bad request", err)
		return
	}

	ctx := r.Context()
	if err := h.uc.DeleteObject(ctx, req.Key); err != nil {
		httpResponse.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpResponse.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) PresignMulti(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Files []port.UploadInput `json:"files"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpResponse.WriteError(w, http.StatusBadRequest, "bad request", err)
		return
	}
	ctx := r.Context()
	resp, err := h.uc.PresignMany(ctx, req.Files)
	if err != nil {
		httpResponse.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	httpResponse.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) StartMultipart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key         string `json:"key"`
		ContentType string `json:"content_type"`
		MaxSize     int64  `json:"max_size"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpResponse.WriteError(w, http.StatusBadRequest, "bad request", err)
		return
	}
	ctx := r.Context()
	session, err := h.uc.StartMultipart(ctx, req.Key, req.ContentType, req.MaxSize)
	if err != nil {
		httpResponse.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpResponse.WriteJSON(w, http.StatusOK, session)
}

func (h *Handler) PresignPart(w http.ResponseWriter, r *http.Request) {
	uploadID := chi.URLParam(r, "uploadID")
	var req struct {
		Key       string `json:"key"`
		PartNumber int32 `json:"part_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpResponse.WriteError(w, http.StatusBadRequest, "bad request", err)
		return
	}

	ctx := r.Context()
	resp, err := h.uc.GetMultipartPartURL(ctx, req.Key, uploadID, req.PartNumber)
	if err != nil {
		httpResponse.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpResponse.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) CompleteMultipart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key    string          `json:"key"`
		UploadID string        `json:"upload_id"`
		Parts  []port.UploadPart `json:"parts"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpResponse.WriteError(w, http.StatusBadRequest, "bad request", err)
		return
	}
	
	ctx := r.Context()
	if err := h.uc.CompleteMultipart(ctx, req.Key, req.UploadID, req.Parts); err != nil {
		httpResponse.WriteError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpResponse.WriteJSON(w, http.StatusNoContent, nil)
}
