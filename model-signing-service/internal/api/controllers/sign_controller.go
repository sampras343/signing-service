package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"github.com/sampras343/signing-service/model-signing-service/internal/app"
	"github.com/sampras343/signing-service/model-signing-service/internal/api/responses"
)

func Sign(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid form: "+err.Error())
		return
	}

	files := r.MultipartForm.File["inputDir"]
	if len(files) == 0 {
		responses.Error(w, http.StatusBadRequest, "No files uploaded")
		return
	}

	tmpDir, _ := os.MkdirTemp("", "inputDir")
	defer os.RemoveAll(tmpDir)

	for _, f := range files {
		src, _ := f.Open()
		defer src.Close()
		dst, _ := os.Create(filepath.Join(tmpDir, f.Filename))
		defer dst.Close()
		io.Copy(dst, src)
	}

	cfg := app.Config{
		InputDir:    tmpDir,
		OutputZip:   "output_api/signed_bundle.zip",
		PrivKeyPath: "keys/private.pem",
		PubKeyPath:  "keys/public.pem",
	}

	svc, err := app.BuildSigningService(cfg)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := svc.SignAndBundle(cfg.InputDir, cfg.OutputZip); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.ServeFile(w, r, cfg.OutputZip)
}
