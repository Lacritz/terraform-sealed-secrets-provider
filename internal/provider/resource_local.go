package provider

import (
	"context"
	"crypto/rsa"
	"crypto/sha1"
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
		Description:   "Creates a sealed secret and stores it in yaml_content.",
		ReadContext:   resourceLocalRead,
		UpdateContext: resourceLocalRead,
		CreateContext: resourceLocalCreate,
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			d.SetId("")
			return nil
		},
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
				Description: "Key/value pairs to populate the secret. The value will be base64 encoded",
			},
			"yaml_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The produced sealed secret yaml file.",
			},
			"public_key_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "The public key hashed to detect if the public key changes.",
			},
		},
	}
}

// resourceLocalRead creates only a hash of the public key.
// If the hash changes then the resource is forced recreated.
func resourceLocalRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Config)
	pk, err := provider.PublicKeyResolver(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("name").(string))
	d.Set("data", d.Get("data").(map[string]interface{}))

	newPkHash := hashPublicKey(pk)
	if oldPkHash, ok := d.GetOk("public_key_hash"); ok && oldPkHash.(string) != newPkHash {
		d.SetId("")
	}
	d.Set("public_key_hash", newPkHash)

	return nil
}

func resourceLocalCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	d.Set("data", d.Get("data").(map[string]interface{}))
	d.Set("yaml_content", string(sealedSecret))
	d.Set("public_key_hash", hashPublicKey(pk))

	return nil
}

func createK8sSecret(d *schema.ResourceData) (v1.Secret, error) {
	rawSecret := k8s.SecretManifest{
		Name:      d.Get("name").(string),
		Namespace: d.Get("namespace").(string),
		Type:      d.Get("type").(string),
	}
	if dataRaw, ok := d.GetOk("data"); ok {
		rawSecret.Data = dataRaw.(map[string]interface{})
	}

	return k8s.CreateSecret(&rawSecret)
}

func hashPublicKey(pk *rsa.PublicKey) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%v%v", pk.N, pk.E))))
}

func logDebug(s string) {
	log.Printf("[DEBUG] %s", s)
}
