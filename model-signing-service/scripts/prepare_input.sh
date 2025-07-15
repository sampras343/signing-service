
#!/bin/bash

mkdir input-folder

echo "Dummy model binary" > input-folder/model.onnx
echo "id,text,label\n1,Hello,happy" > input-folder/dataset.csv

cat <<EOF > input-folder/metadata.json
{
  "author": "Sachin Sampras",
  "timestamp": "2025-07-15T10:00:00Z",
  "framework": "PyTorch",
  "dataset_id": "demo-v1",
  "description": "Dummy model for signing"
}
EOF
