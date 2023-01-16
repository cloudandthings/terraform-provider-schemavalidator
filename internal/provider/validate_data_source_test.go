package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccValidateDataSource(t *testing.T) {
	var cases = []struct {
		document      string
		schema        string
		errorExpected bool
	}{
		{"asd asdasd: ^%^*&^%", "{}", true},
		{"{}", schemaValid, true},
		{`{"test": "test"}`, schemaValid, false},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%s with error expected %t", tt.document, tt.errorExpected), func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: makeDataSource(tt.document, tt.schema),
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
