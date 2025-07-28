package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/sampras343/signing-service/model-signing-service/internal/api/responses"
	"github.com/sampras343/signing-service/model-signing-service/internal/service"
)

func Verify(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		responses.Error(w, http.StatusBadRequest, "Invalid form: "+err.Error())
		return
	}

	file, _, err := r.FormFile("bundle")
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Bundle file required")
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "bundle-*.zip")
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, file); err != nil {
		responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := service.VerifyBundle(tmpFile.Name()); err != nil {
		responses.Error(w, http.StatusBadRequest, "Verification failed: "+err.Error())
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{
		"status": "Signature verified. Bundle authentic.",
		"code": "200",
	})
}
