# Signing Service

## 📑 Table of Contents

- [Signing Service](#signing-service)
  - [Pre-Requisites](#pre-requisites)
  - [How To Run?](#how-to-run)
  - [Features](#features)


## Pre-Requisites
- Clone the repository

- Navigate to the source directory
```bash
cd signing-service/model-signing-service/
```

- Run script to generate necessary keys
```bash
bash scripts/gen_keys.sh
```

- Run script to generate dummy input files
```bash
bash scripts/prepare_input.sh
```

## How To Run?
- Below steps have to be run from `signing-service/model-signing-service/`
- Replace the existing dummy artifacts with the actuals in the input-folder or a custom folder and ensure to provide the correct path in the next step. 

- Run the signer
```bash
go run cmd/signer/main.go --input input-folder --output output/signed_bundle.zip
```

- Optional: Overridet the path of public & private keys as well
```bash
go run cmd/signer/main.go --input <input-folder> --output <output.zip> --priv <path/to/private.pem> --pub <path/to/public.pem>
```

- This will create a signed artifact under the folder output.

- Optional: Unzip this signed artifact to explore what is available in the zip

- Verify the integrity of the atrifact
```bash
go run cmd/verifier/main.go output/signed_bundle.zip
```

## Features

### Input 

1. The signing service will take the input artifacts from a folder
    - The input artifacts should consist of details such as the model file, the dataset used to train, the author details
    - The interface should be flexible enough to accept inputs from any source (example: RESTful endpoint or Hub) at a later stage


### Core Logic
1. The signing service will accept certificate from a local directory but it should be easily replacable with any other interface at a later stage of the project.
2. The signing service should create hashes for each of the artifacts an a given instance to maintain the authenticity and avoid tampering at any given stage.
3. The core signing logic should also be abstracted away in case a different signing mechanism is used at a later stage.

### Outputs
1. The signing service will provide a manifest of the output
2. The signing service should bundle these artifacts along with the signature in .zip or .tar format
3. The output will be written onto a local directory however it should open and easy enough to replace with any kind of output mechanism that could be evolved to in future. (Example: Pushing the resources to a remote repository)

### Verification
1. The complete bundle will be provided for verification. All checks must take place to determine if the digital signature has been tampered with or not.

### Miscellaneous
1. Should be built and run as a simple executable.
2. CLI interface
3. REST interface

