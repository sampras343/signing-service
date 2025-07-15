package sink

import (
    "archive/zip"
    "encoding/json"
    "io"
    "os"
    "path/filepath"

    "github.com/sampras343/signing-service/model-signing-service/internal/model"
)

type ZipArtifactSink struct{}

func (z *ZipArtifactSink) Write(bundlePath string, manifest *model.SignedManifest, artifacts []string) error {
    out, err := os.Create(bundlePath)
    if err != nil {
        return err
    }
    defer out.Close()

    zipWriter := zip.NewWriter(out)
    defer zipWriter.Close()

    // Write artifacts
    for _, file := range artifacts {
        f, err := os.Open(file)
        if err != nil {
            return err
        }
        defer f.Close()

        w, err := zipWriter.Create(filepath.Base(file))
        if err != nil {
            return err
        }
        _, err = io.Copy(w, f)
        if err != nil {
            return err
        }
    }

    // Write manifest
    manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
    if err != nil {
        return err
    }
    mw, err := zipWriter.Create("signed_manifest.json")
    if err != nil {
        return err
    }
    _, err = mw.Write(manifestBytes)
    return err
}
