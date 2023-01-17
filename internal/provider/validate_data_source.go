package provider

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ValidateDataSource{}

func NewValidateDataSource() datasource.DataSource {
	return &ValidateDataSource{}
}

// ValidateDataSource defines the data source implementation.
type ValidateDataSource struct{}

// ValidateDataSourceModel describes the data source data model.
type ValidateDataSourceModel struct {
	Document types.String `tfsdk:"document"`
	Schema   types.String `tfsdk:"schema"`

	Validated types.Bool   `tfsdk:"validated"`
	ID        types.String `tfsdk:"id"`
}

func (d *ValidateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_validate"
}

func (d *ValidateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Validate a json document against a [json-schema](https://json-schema.org/) schema file.",

		Attributes: map[string]schema.Attribute{
			"document": schema.StringAttribute{
				MarkdownDescription: "Content of a json file",
				Required:            true,
			},
			"schema": schema.StringAttribute{
				MarkdownDescription: "Content of a [json-schema](https://json-schema.org/) file",
				Required:            true,
			},
			// Computed Attributes
			"validated": schema.BoolAttribute{
				MarkdownDescription: "True if the document's schema could be validated.",
				Computed:            true,
				Optional:            false,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Id generated for document and schema hashes",
				Computed:            true,
			},
		},
	}
}

func (d *ValidateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data ValidateDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	data.ID = types.StringValue(hash(data.Document.ValueString() + data.Schema.ValueString()))
	data.Validated = types.BoolValue(false)

	if resp.Diagnostics.HasError() {
		return
	}

	sch, err := jsonschema.CompileString("schema.json", data.Schema.ValueString())
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Invalid Schema: %s", err))
		resp.Diagnostics.AddError("Invalid Schema", fmt.Sprintf("Got error: %s", err))
		return
	}

	var document interface{}
	if err := json.Unmarshal([]byte(data.Document.ValueString()), &document); err != nil {
		tflog.Error(ctx, fmt.Sprintf("Invalid JSON Document: %s", err))
		resp.Diagnostics.AddError("Invalid JSON Document", fmt.Sprintf("Got error: %s", err))
		return
	}

	if err := sch.Validate(document); err != nil {
		tflog.Error(ctx, fmt.Sprintf("Document not valid, got error: %s", err))
		resp.Diagnostics.AddError("Document not valid", fmt.Sprintf("Got error: %s", err))
		return
	}

	data.Validated = types.BoolValue(true)

	// // Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
