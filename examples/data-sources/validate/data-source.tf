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

data "schemavalidator_validate" "from_fs" {
  document = file("./document.json")
  schema   = file("./schema.json")
}

output "validated" {
  value = data.schemavalidator_validate.test.validated
}