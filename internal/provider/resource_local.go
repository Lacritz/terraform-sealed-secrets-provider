package provider

import (
	"context"
	"fmt"
	"github.com/datalbry/sealedsecret/internal/k8s"
	"github.com/datalbry/sealedsecret/internal/kubeseal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v1 "k8s.io/api/core/v1"
	"log"
)

func resourceLocal() *schema.Resource {
	return &schema.Resource{
		Description: "Creates a sealed secret and stores it in yaml_content.",
		ReadContext: resourceLocalRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the secret, must be unique.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace of the secret.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Opaque",
				Description: "The secret type (ex. Opaque). Default type is Opaque.",
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Sensitive:   true,
				Description: "Key/value pairs to populate the secret.",
			},
			"binary_data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Sensitive:   true,
				Description: "Base64 encoded key/value pairs to populate the secret.",
			},
			"yaml_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The produced sealed secret yaml file.",
			},
		},
	}
}

func resourceLocalRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Config)
	name := d.Get("name").(string)

	logDebug("Creating sealed secret " + name)
	k8sSecret, err := createK8sSecret(d)
	if err != nil {
		return diag.FromErr(err)
	}
	pk, err := provider.PublicKeyResolver(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	sealedSecret, err := kubeseal.SealSecret(k8sSecret, pk)
	if err != nil {
		return diag.FromErr(err)
	}

	logDebug("Successfully created sealed secret " + name)

	d.SetId(name)
	err = d.Set("data", d.Get("data").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("binary_data", d.Get("binary_data").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("yaml_content", string(sealedSecret))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createK8sSecret(d *schema.ResourceData) (v1.Secret, error) {
	rawSecret := k8s.SecretManifest{
		Name:      d.Get("name").(string),
		Namespace: d.Get("namespace").(string),
		Type:      d.Get("type").(string),
	}
	if dataRaw, ok := d.GetOk("data"); ok {
		rawSecret.Data = mapValuesToBytes(dataRaw.(map[string]interface{}))
	}
	if binaryDataRaw, ok := d.GetOk("binary_data"); ok {
		rawSecret.BinaryData = mapValuesToString(binaryDataRaw.(map[string]interface{}))
	}

	return k8s.CreateSecret(&rawSecret)
}

func logDebug(s string) {
	log.Printf("[DEBUG] %s", s)
}

func mapValuesToBytes(data map[string]interface{}) map[string][]byte {
	result := make(map[string][]byte)
	for key, value := range data {
		result[key] = []byte(fmt.Sprintf("%v", value))
	}
	return result
}

func mapValuesToString(data map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range data {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}
