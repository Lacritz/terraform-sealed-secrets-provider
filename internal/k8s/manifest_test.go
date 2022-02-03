package k8s

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSecret(t *testing.T) {
	secretKey, secretValue := "secret", "secret_aaa"
	tests := []struct {
		Name              string
		Input             SecretManifest
		ExpectedDataValue string
	}{
		{
			Name: "happy day",
			Input: SecretManifest{
				Name:      "name_aaa",
				Namespace: "ns_aaa",
				Type:      "type_aaa",
				Data:      map[string][]byte{secretKey: []byte(secretValue)},
			},
			ExpectedDataValue: secretValue,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			secret, err := CreateSecret(&tc.Input)

			assert.NoError(t, err)
			assert.Equal(t, tc.Input.Name, secret.Name)
			assert.Equal(t, tc.Input.Namespace, secret.Namespace)
			assert.Equal(t, tc.Input.Type, string(secret.Type))
			assert.Equal(t, tc.ExpectedDataValue, string(secret.Data[secretKey]))
		})
	}

}
