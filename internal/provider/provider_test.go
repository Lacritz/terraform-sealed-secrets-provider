package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

var testAccProvider *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"gitlabcommit": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func init() {
	testAccProvider = Provider()
	testAccProvider.Configure(context.Background(), &terraform.ResourceConfig{})
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"sealedsecret": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	errMsg := "%s env var must be set for acceptance test"
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Fatalf(errMsg, "TF_ACC")
	}
	if v := os.Getenv("HOST"); v == "" {
		t.Fatalf(errMsg, "HOST")
	}
	if v := os.Getenv("CLIENT_CERTIFICATE"); v == "" {
		t.Fatalf(errMsg, "CLIENT_CERTIFICATE")
	}
	if v := os.Getenv("CLIENT_KEY"); v == "" {
		t.Fatalf(errMsg, "CLIENT_KEY")
	}
	if v := os.Getenv("CLUSTER_CA_CERTIFICATE"); v == "" {
		t.Fatalf(errMsg, "CLUSTER_CA_CERTIFICATE")
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
