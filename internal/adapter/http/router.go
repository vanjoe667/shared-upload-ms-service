package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func FileUploadRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()
	basePath := "/api/uploads"

	r.Post(basePath+"/presign", handler.PresignSingle)
	r.Post(basePath+"/presign-multi", handler.PresignMulti)
	r.Post(basePath+"/delete", handler.DeleteObject)
	r.Post(basePath+"/presign-download", handler.PresignDownload)

	// Multipart flows
	r.Post(basePath+"/multipart/start", handler.StartMultipart)
	r.Post(basePath+"/multipart/{uploadID}/presign", handler.PresignPart)
	r.Post(basePath+"/multipart/{uploadID}/complete", handler.CompleteMultipart)

	return r
}