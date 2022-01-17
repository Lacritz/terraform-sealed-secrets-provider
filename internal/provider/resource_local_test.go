package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/yaml"

	"testing"
)

func TestAccResourceLocal(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSealedSecretLocal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckYamlContentProduced(t, "sealedsecret_local.test"),
				),
			},
		},
	})
}

const testAccResourceSealedSecretLocal = `
provider "sealedsecret"{
	kubernetes {}
}
resource "sealedsecret_local" "test" {
	name = "secret"
	namespace = "default"
	data = {
		"secret-key": "secret-value"		
	}
}
`

func testAccCheckYamlContentProduced(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		assert.Equal(t, rs.Primary.ID, "secret")
		sealedSecret := struct {
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
		err := yaml.Unmarshal([]byte(rs.Primary.Attributes["yaml_content"]), &sealedSecret)
		assert.NoError(t, err)

		assert.Equal(t, "secret", sealedSecret.Metadata.Name)
		assert.Equal(t, "secret", sealedSecret.Spec.Template.Metadata.Name)
		assert.Equal(t, "default", sealedSecret.Metadata.Namespace)
		assert.Equal(t, "default", sealedSecret.Spec.Template.Metadata.Namespace)

		// Can only check length since encrypted content changes even though the public key is the same
		assert.Len(t, sealedSecret.Spec.EncryptedData["secret-key"], 724)
		return nil
	}
}
