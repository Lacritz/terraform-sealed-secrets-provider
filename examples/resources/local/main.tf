terraform {
  required_providers {
    sealedsecret = {
      version = ">=1.0.0"
      source  = "datalbry/sealed_secret"
    }
  }
}

provider "sealedsecret" {
  controller_name      = "sealed-secret-controller-sealed-secrets"
  controller_namespace = "kube-system"
  pem                  = "weathuawetl...awethiawe"
}

data "sealed_secret" "example" {
  name      = "example-secret"
  namespace = "default"
  data      = {
    "key" : "value"
  }
}

resource "local_file" "example" {
  filename = "sealed-secret.yaml"
  content  = data.sealed_secret.example.yaml_content
}
