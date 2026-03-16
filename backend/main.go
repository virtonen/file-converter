package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strings"
)

const maxUploadSize = 10 << 20 // 10 MB

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "failed to parse upload", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		http.Error(w, "cannot decode image: "+err.Error(), http.StatusBadRequest)
		return
	}
	if format != "jpeg" {
		http.Error(w, "only JPEG input is supported", http.StatusBadRequest)
		return
	}

	// Derive output filename from original name (case-insensitive extension handling)
	base := header.Filename
	dotIdx := strings.LastIndex(base, ".")
	var outName string
	if dotIdx > 0 {
		ext := strings.ToLower(base[dotIdx+1:])
		switch ext {
		case "jpg", "jpeg":
			outName = base[:dotIdx] + ".png"
		default:
			outName = base + ".png"
		}
	} else {
		outName = base + ".png"
	}

	// Encode to an in-memory buffer so we can detect encoding errors before
	// writing the response headers.
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		http.Error(w, "failed to encode PNG: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, outName))
	w.Write(buf.Bytes())
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", convertHandler)

	// Serve the frontend index.html
	mux.Handle("/", http.FileServer(http.Dir("../frontend")))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, corsMiddleware(mux)))
}

// corsMiddleware adds permissive CORS headers so the frontend can talk to the
// backend regardless of how it is served during development.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
