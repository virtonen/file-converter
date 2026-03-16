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

- **[Go](https://go.dev/dl/) 1.21+**
  - **Windows (PowerShell):** Download and run the `.msi` installer. After installation, **restart your terminal** so the `go` command is recognized.
  - **macOS:** Download and run the `.pkg` installer or use Homebrew (`brew install go`).
- A modern **web browser** (Chrome, Firefox, Edge, Safari)

## Running

### 1 – Start the backend

#### Windows (PowerShell)
```powershell
cd backend
go run main.go
# Listening on :8080
```

#### macOS / Linux (Bash/Zsh)
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

#### Windows (PowerShell)
```powershell
cd backend
go test -v ./...
```

#### macOS / Linux (Bash/Zsh)
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