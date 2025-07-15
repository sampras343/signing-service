package sink

import "github.com/sampras343/signing-service/model-signing-service/internal/model"

type ArtifactSink interface {
    Write(bundlePath string, manifest *model.SignedManifest, artifacts []string) error
}
