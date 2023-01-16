package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure SchemaValidatorProvider satisfies various provider interfaces.
var _ provider.Provider = &SchemaValidatorProvider{}

// SchemaValidatorProvider defines the provider implementation.
type SchemaValidatorProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SchemaValidatorProviderModel describes the provider data model.
type SchemaValidatorProviderModel struct{}

func (p *SchemaValidatorProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "schemavalidator"
	resp.Version = p.version
}

func (p *SchemaValidatorProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description:         "Validate a json document against a [json-schema](https://json-schema.org/) file",
		MarkdownDescription: "Validate a json document against a [json-schema](https://json-schema.org/) file",

		Attributes: map[string]schema.Attribute{},
	}
}

func (p *SchemaValidatorProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *SchemaValidatorProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *SchemaValidatorProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewValidateDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SchemaValidatorProvider{
			version: version,
		}
	}
}
