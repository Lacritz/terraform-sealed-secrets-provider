package kubeseal

import (
	"context"
	"github.com/datalbry/sealedsecret/internal/k8s"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/yaml"
	"log"
	"testing"
)

const pem = `-----BEGIN CERTIFICATE-----
MIIErjCCApagAwIBAgIRAIrkLt+H5TI6sZojiRnT0KswDQYJKoZIhvcNAQELBQAw
ADAeFw0yMTA3MDUxMzExMjhaFw0zMTA3MDMxMzExMjhaMAAwggIiMA0GCSqGSIb3
DQEBAQUAA4ICDwAwggIKAoICAQDQymZt7IoS0gQn8lA0UNCFpbFFPF5VK+zygi0f
+HHG4OrNMbCTpHVx3wSWIAkMyG+JvIg6yPb+oYA6SS+n8T3OVH1I+PiIqf4ZQOZW
yB0CH9b5l+lQ3pFgaysByrc2ONBsnIuqSNbm5z02P+d9oeFcC+htxQsWFPIG8TdI
2d/FuRsMa/mx67hM39raa0D3RiWKMpp9gR5H+eUskGIyFuZvmPHWVyGUt38lJQ8I
4jhwCGwZv4P7DtjaoDS3O9FmE0O/aao4vn0gSGJa5EBjInYv5bakck1TPoXPbnxL
FwF50rCg0drt6evWJWFDSt9FHYqVa4qucxcHo7FI1Kz6WC7IlMhC9i6PRQcsuCEW
Izsbjuxct9cUg5GsNgbJQR+FXcTjmY8SJA0a+fwVAMWNO8SQjJIJXJ1hgCIMQtgZ
3GzEgULq9EA4Hp59cW5px1XPG7UhxbNlJpZ0oZ7hNw+FbuyyfAgVH/QiXg+Zk7Pa
2cpGuiwDhDjOWDiuOONfuYa6a2KDrdDXXpxvI+lMltkQzM/rTtDgEXyd5Cvn6NIh
HuDqq4ffvrJrnwtzUEnmqsE56Gv/tFfcCdh+tqSTwkpbXj649yA2sQ8ByozEiWoV
Wb8mjs/ckrmtjrl6IoABYTrQZ3j8TXY5F/kOkygx2nidd2zsJyD5RGEEgAkzY7uc
dFa8TQIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAAEwDwYDVR0TAQH/BAUwAwEB/zAN
BgkqhkiG9w0BAQsFAAOCAgEAQ0nc7NsAim2DIvd4KOsHDt7PbOsoAOe4bscvJsLT
GTJaYWkp85oHXDoOrDpsN7baLtyxeVW4+QHFZrauXcQGg45a8lHo4LV+RUzylms8
xJkLRJPWHW5YsovXBS/roBVTH3irC6VRSgEf3j0VQNK0jVXr+V5VPSnZJw4mXeW2
pN640DHGxHfdniSRXeiACWxAD+k77eg6VoMb+hk3U6em43TlraP9vNz5nC16Vd2Y
o9cVQ+Z0G0hGyR3vBWEOSFWsh11NJKimMtUjdE2qGokLEDDur7Rktqq6oRvhz4na
+TRymQ+up9zAtY1Sj396lPzf/s8KHNgC3Y0lC0YMHvLWZiza4SVb35cn5oHnqn2L
/4zIjOijx4i/wVS3j1nGqYuBbKbF6zDZuu4vWElEl4F7LY0h6816Z2nCwH/wMyZt
ebWk/voTwBDRp/u+/GwRRL6nH9Bsrx7zN3xRv/Lm6jmhLrt7rfdAaZXgUY4hEIJ7
TyU/weNmm52KUO9KQCLOe9z0cogrjKf5L87cyCddI2pKtI/IdA1qDLpIjtaxPSA9
YvfoFuqfA/Ps8zzHv+yyYJb98v9clUyoSyXYFoRmHqvLVTEEZZG8Ggo83v9EzERu
lffms06J9FgeTsCaydCx+jGFW1eOeBZc7Cbh9kO/DuoIqIa2RvBMVoOBf4eWkbEV
p+0=
-----END CERTIFICATE-----
`

func TestResolvePK(t *testing.T) {
	pk, err := ResolvePK([]byte(pem))(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, 65537, pk.E)
}

func TestSealSecret(t *testing.T) {
	sm := k8s.SecretManifest{
		Name:      "name_aa",
		Namespace: "ns_aa",
		Type:      "type_aa",
		Data: map[string]interface{}{
			"keyAA": "secret",
		},
	}

	pk, err := ResolvePK([]byte(pem))(context.Background())
	assert.Nil(t, err)

	secret, err := k8s.CreateSecret(&sm)
	assert.Nil(t, err)
	sealedSecretRaw, err := SealSecret(secret, pk)
	assert.Nil(t, err)

	actualSS := struct {
		Kind     string `yaml:"kind"`
		Metadata struct {
			Name      string `yaml:"name"`
			Namespace string `yaml:"namespace"`
		} `yaml:"metadata"`
		Spec struct {
			EncryptedData map[string]string `yaml:"encryptedData"`
			Template      struct {
				Data     interface{} `yaml:"data"`
				Metadata struct {
					Name      string `yaml:"name"`
					Namespace string `yaml:"namespace"`
				} `yaml:"metadata"`
				Type string `yaml:"type"`
			} `yaml:"template"`
		} `yaml:"spec"`
	}{}

	err = yaml.Unmarshal(sealedSecretRaw, &actualSS)
	assert.Nil(t, err)

	assert.Equal(t, sm.Name, actualSS.Metadata.Name)
	assert.Equal(t, sm.Name, actualSS.Spec.Template.Metadata.Name)

	assert.Equal(t, sm.Namespace, actualSS.Metadata.Namespace)
	assert.Equal(t, sm.Namespace, actualSS.Spec.Template.Metadata.Namespace)

	assert.Equal(t, "SealedSecret", actualSS.Kind)
	assert.Equal(t, sm.Type, actualSS.Spec.Template.Type)
	if len(actualSS.Spec.EncryptedData["keyAA"]) < 600 {
		log.Println(actualSS)
		t.Errorf("expected long encrypted string, got %s", actualSS.Spec.EncryptedData["keyAA"])
	}
}