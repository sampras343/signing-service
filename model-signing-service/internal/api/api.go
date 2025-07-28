package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"github.com/sampras343/signing-service/model-signing-service/internal/app"
	"github.com/sampras343/signing-service/model-signing-service/internal/service"
)

func SignHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	inputDir, err := saveUploadedDir(r, "inputDir")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer os.RemoveAll(inputDir)

	cfg := app.Config{
		InputDir:    inputDir,
		OutputZip:   "output_api/signed_bundle.zip",
		PrivKeyPath: "keys/private.pem",
		PubKeyPath:  "keys/public.pem",
	}
	svc, err := app.BuildSigningService(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := svc.SignAndBundle(cfg.InputDir, cfg.OutputZip); err != nil {
		http.Error(w, "Signing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, cfg.OutputZip)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("bundle")
	if err != nil {
		http.Error(w, "Bundle file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "bundle-*.zip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := service.VerifyBundle(tmpFile.Name()); err != nil {
		http.Error(w, "Verification failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]string{"status": "✅ Signature verified. Bundle authentic."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func saveUploadedDir(r *http.Request, field string) (string, error) {
	files := r.MultipartForm.File[field]
	if len(files) == 0 {
		return "", fmt.Errorf("no files in form field %s", field)
	}
	tmpDir, _ := os.MkdirTemp("", "inputDir")
	for _, f := range files {
		src, _ := f.Open()
		defer src.Close()

		dstPath := filepath.Join(tmpDir, f.Filename)
		dst, _ := os.Create(dstPath)
		defer dst.Close()
		io.Copy(dst, src)
	}
	return tmpDir, nil
}
