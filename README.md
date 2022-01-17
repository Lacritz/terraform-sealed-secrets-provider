# Sealed Secrets Provider
The Terraform Sealed Secrets Provider can setup sealed-secrets using a local public key.

* Using a passed controller's public key.
* Encrypts the provided secret

The sealed secret manifest is computed on the `yaml_content` field.

## Getting Started 

### Building (local)
The package can be bundled using
```shell
make build-local
```

### Testing
All packages can be tested using:
```shell
go test ./internal/*
``` 