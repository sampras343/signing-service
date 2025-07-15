package source

import "github.com/sampras343/signing-service/model-signing-service/internal/model"

type ArtifactSource interface {
    Load(path string) (*model.ArtifactBundle, error)
}
