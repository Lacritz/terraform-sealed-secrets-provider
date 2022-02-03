package k8s

import (
	"encoding/base64"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretManifest struct {
	Name       string
	Namespace  string
	Type       string
	Data       map[string][]byte
	BinaryData map[string]string
}

func CreateSecret(manifest *SecretManifest) (v1.Secret, error) {
	if manifest.Data == nil {
		manifest.Data = make(map[string][]byte)
	}
	b64DecodeMapValues(manifest.BinaryData, &manifest.Data)

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifest.Name,
			Namespace: manifest.Namespace,
		},
		Type: v1.SecretType(manifest.Type),
		Data: manifest.Data,
	}

	return secret, nil
}

func b64DecodeMapValues(source map[string]string, target *map[string][]byte) {
	for key, value := range source {
		(*target)[key], _ = base64.StdEncoding.DecodeString(value)
	}
}
