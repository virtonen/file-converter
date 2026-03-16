package main

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// makeSampleJPEG returns a minimal valid JPEG (a 4×4 red square).
func makeSampleJPEG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("jpeg.Encode: %v", err)
	}
	return buf.Bytes()
}

func TestConvertJPEGToPNG(t *testing.T) {
	jpegData := makeSampleJPEG(t)

	// Build a multipart/form-data request body
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile("file", "sample.jpg")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := fw.Write(jpegData); err != nil {
		t.Fatalf("write jpeg data: %v", err)
	}
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/convert", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	rec := httptest.NewRecorder()
	convertHandler(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", resp.StatusCode, rec.Body.String())
	}

	// Verify Content-Type
	ct := resp.Header.Get("Content-Type")
	if ct != "image/png" {
		t.Errorf("expected Content-Type image/png, got %s", ct)
	}

	// Verify the body is a valid PNG
	if _, err := png.Decode(resp.Body); err != nil {
		t.Errorf("response body is not a valid PNG: %v", err)
	}

	// Verify Content-Disposition contains the .png filename
	cd := resp.Header.Get("Content-Disposition")
	if cd == "" {
		t.Error("expected Content-Disposition header")
	}
	t.Logf("Content-Disposition: %s", cd)
}

func TestConvertRejectsNonJPEG(t *testing.T) {
	// Create a valid PNG and try to upload it — should be rejected.
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile("file", "sample.png")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := fw.Write(buf.Bytes()); err != nil {
		t.Fatalf("write png data: %v", err)
	}
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/convert", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	rec := httptest.NewRecorder()
	convertHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-JPEG input, got %d", rec.Code)
	}
}

func TestConvertMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/convert", nil)
	rec := httptest.NewRecorder()
	convertHandler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
