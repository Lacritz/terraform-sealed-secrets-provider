package provider

import (
	"context"
	"github.com/datalbry/sealedsecret/internal/kubeseal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"controller_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of k8s service for the sealed-secret-controller.",
				Default:     "sealed-secret-controller-sealed-secrets",
			},
			"controller_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The namespace the controller is running in.",
				Default:     "kube-system",
			},
			"pem": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The content of the pem file of the sealed-secret-controller.",
			},
		},
		ConfigureContextFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"sealed_secret": resourceLocal(),
		},
	}
}

type Config struct {
	ControllerName      string
	ControllerNamespace string
	PublicKeyResolver   kubeseal.PKResolverFunc
}

func configureProvider(_ context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
	cName := rd.Get("controller_name").(string)
	cNs := rd.Get("controller_namespace").(string)
	pk := []byte(rd.Get("pem").(string))

	return &Config{
		ControllerName:      cName,
		ControllerNamespace: cNs,
		PublicKeyResolver:   kubeseal.ResolvePK(pk),
	}, nil
}

