package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccValidateDataSource(t *testing.T) {
	var cases = []struct {
		config        func(string, string) string
		document      string
		schema        string
		errorExpected bool
	}{
		{makeDataSource, "asd asdasd: ^%^*&^%", "{}", true},
		{makeDataSource, "{}", schemaValid, true},
		{makeDataSource, `{"test": "test"}`, schemaValid, false},
		{makeDataSourceYaml, k8sDeployment, schemaWithHttpLink, false},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%s with error expected %t", tt.document, tt.errorExpected), func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: tt.config(tt.document, tt.schema),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.schemavalidator_validate.test", "validated", strconv.FormatBool(!tt.errorExpected)),
						),
					},
				},
				ErrorCheck: func(err error) error {
					if tt.errorExpected {
						if err == nil {
							return fmt.Errorf("error expected")
						} else {
							return nil
						}
					}
					return err
				},
			})
		})
	}
}

func makeDataSource(document string, schema string) string {
	return fmt.Sprintf(`
data "schemavalidator_validate" "test" {
	document = <<EOF
%s
EOF
	schema   = <<EOF
%s
EOF
}
`, document, schema)
}

const schemaValid = `{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"x-$id": "https://example.com",
	"type": "object",
	"required": ["test"]
  }`

func makeDataSourceYaml(document string, schema string) string {
	return fmt.Sprintf(`
data "schemavalidator_validate" "test" {
	document = jsonencode(yamldecode(<<EOF
%s
EOF
))
	schema   = <<EOF
%s
EOF
}
`, document, schema)
}

const k8sDeployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
`

const schemaWithHttpLink = `{
	"description": "Deployment enables declarative updates for Pods and ReplicaSets.",
	"properties": {
	  "apiVersion": {
		"description": "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
		"type": [
		  "string",
		  "null"
		],
		"enum": [
		  "apps/v1"
		]
	  },
	  "kind": {
		"description": "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
		"type": [
		  "string",
		  "null"
		],
		"enum": [
		  "Deployment"
		]
	  },
	  "metadata": {
		"$ref": "https://kubernetesjsonschema.dev/v1.14.0/_definitions.json#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta",
		"description": "Standard object metadata."
	  },
	  "spec": {
		"$ref": "https://kubernetesjsonschema.dev/v1.14.0/_definitions.json#/definitions/io.k8s.api.apps.v1.DeploymentSpec",
		"description": "Specification of the desired behavior of the Deployment."
	  },
	  "status": {
		"$ref": "https://kubernetesjsonschema.dev/v1.14.0/_definitions.json#/definitions/io.k8s.api.apps.v1.DeploymentStatus",
		"description": "Most recently observed status of the Deployment."
	  }
	},
	"type": "object",
	"x-kubernetes-group-version-kind": [
	  {
		"group": "apps",
		"kind": "Deployment",
		"version": "v1"
	  }
	],
	"$schema": "http://json-schema.org/schema#"
  }
`
