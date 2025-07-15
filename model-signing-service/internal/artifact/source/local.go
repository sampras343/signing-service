package source

import (
    "encoding/json"
    "os"
    "path/filepath"
    "github.com/sampras343/signing-service/model-signing-service/internal/model"
)

type LocalArtifactSource struct{}

func (l *LocalArtifactSource) Load(path string) (*model.ArtifactBundle, error) {
    metadataFile := filepath.Join(path, "metadata.json")
    modelFile := filepath.Join(path, "model.onnx")
    datasetFile := filepath.Join(path, "dataset.csv")

    var meta model.ModelMetadata
    data, err := os.ReadFile(metadataFile)
    if err != nil {
        return nil, err
    }
    if err := json.Unmarshal(data, &meta); err != nil {
        return nil, err
    }

    return &model.ArtifactBundle{
        ModelFilePath:   modelFile,
        DatasetFilePath: datasetFile,
        Metadata:        meta,
    }, nil
}