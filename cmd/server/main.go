package main

import (
	"context"
	"log"
	"net/http"
	"os"

	httpapi "go.peniremit.dev/peniremit-shared-upload-service/internal/adapter/http"
	bootstrap "go.peniremit.dev/peniremit-shared-upload-service/internal/boostrap"
)

func main() {
    ctx := context.Background()

    uc, err := bootstrap.InitUploadService(ctx)
    if err != nil {
        log.Fatal(err)
    }

    handler := httpapi.FileUploadHandler(uc)
    router := httpapi.FileUploadRouter(handler)

   port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not specified
	}
	log.Println("Server running on port: " + port)
	http.ListenAndServe(":"+port, router)
}