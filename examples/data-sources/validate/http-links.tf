terraform {
  required_providers {
    http = {
      source  = "hashicorp/http"
      version = "3.2.1"
    }
  }
}

data "http" "k8s_deployment_schema" {
  url = "https://kubernetesjsonschema.dev/v1.14.0/deployment-apps-v1.json"
  request_headers = {
    Accept = "application/json"
  }
}

data "http" "k8s_deployment_example" {
  url = "https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/controllers/nginx-deployment.yaml"
  request_headers = {
    Accept = "application/json"
  }
}

data "schemavalidator_validate" "k8s_deployment" {
  schema   = data.http.k8s_deployment_schema.response_body
  document = jsonencode(yamldecode(data.http.k8s_deployment_example.response_body))
}

output "validated" {
  value = data.schemavalidator_validate.k8s_service.validated
}