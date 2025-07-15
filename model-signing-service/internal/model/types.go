package model

import "time"

type ModelMetadata struct {
    Author      string    `json:"author"`
    Timestamp   time.Time `json:"timestamp"`
    Framework   string    `json:"framework"`
    DatasetID   string    `json:"dataset_id"`
    Description string    `json:"description"`
}

type ArtifactBundle struct {
    ModelFilePath   string
    DatasetFilePath string
    Metadata        ModelMetadata
}

type SignedManifest struct {
    ModelHash   string        `json:"model_hash"`
    DatasetHash string        `json:"dataset_hash"`
    Metadata    ModelMetadata `json:"metadata"`
    Signature   []byte        `json:"signature"`
    PublicKey   []byte        `json:"public_key"`
}
