# file-converter

A minimal JPG → PNG file converter with a Go backend and a single-file React frontend.

```
file-converter/
├── backend/        # Go HTTP server
│   ├── go.mod
│   ├── main.go
│   └── main_test.go
└── frontend/
    └── index.html  # Single-file React app (CDN React, no build step)
```

## Requirements

- [Go](https://golang.org/dl/) 1.21+
- A modern web browser (Chrome, Firefox, Edge, Safari)

## Running

### 1 – Start the backend

```bash
cd backend
go run .
# Listening on :8080
```

The server also serves the frontend at `http://localhost:8080/`.

### 2 – Open the UI

Open **http://localhost:8080** in your browser.

1. Click the drop-zone (or drag-and-drop) to pick a `.jpg` / `.jpeg` file.
2. Click **Convert to PNG**.
3. Download the resulting `.png` file.

## Running tests

```bash
cd backend
go test -v ./...
```

Three tests are included:

| Test | What it checks |
|---|---|
| `TestConvertJPEGToPNG` | Uploads a synthetic JPEG, verifies a valid PNG comes back |
| `TestConvertRejectsNonJPEG` | Uploading a PNG returns HTTP 400 |
| `TestConvertMethodNotAllowed` | GET requests return HTTP 405 |