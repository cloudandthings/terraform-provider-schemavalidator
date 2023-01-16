terraform {
  required_providers {
    schemavalidator = {
      source  = "terraform.local/cloudandthings/schemavalidator"
      version = "1.0.0"
    }
  }
}

data "schemavalidator_validate" "test" {
  document = jsonencode({ "test" = "test" })

  schema = jsonencode(
    {
      "$schema"  = "http://json-schema.org/draft-07/schema#"
      "x-$id"    = "https://example.com"
      "type"     = "object"
      "required" = ["test"]
  })
}

output "sch" {
  value = data.schemavalidator_validate.test.validated
}