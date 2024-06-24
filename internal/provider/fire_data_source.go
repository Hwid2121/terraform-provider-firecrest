package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewFireDataSource() datasource.DataSource {
	return &fireDataSource{}
}

type fireDataSource struct {
	client *FirecrestClient
}

type fireDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	Token types.String `tfsdk:"token"`
	ID types.String `tfsdk:"id"`
}

func (d *fireDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fire"
}

func (d *fireDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the authentication token for the firecREST API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier attribute for testing.",
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Data Source.",
				Required: true,
			},
			"token": schema.StringAttribute{
				Description: "The temporary API token for the firecrREST API.",
				Computed: true,
			},
		},
	}
}

func (d *fireDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.AddWarning("Data Source Configuration", fmt.Sprintf("ProviderData type: %T", req.ProviderData))

	if req.ProviderData == nil {
		return
	}
	
	providerConfig, ok := req.ProviderData.(*firecrestProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			("Expected provider to be of type *firecrestProvider"),
		)
		return
	}

	d.client = providerConfig.client
	resp.Diagnostics.AddWarning("Data Source Configuration", "Configured FirecREST data source successfully.")

}


func (d *fireDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {


	var state fireDataSourceModel
	if d.client == nil {
		resp.Diagnostics.AddError("Client Error", "The client is not configured")
		return
	}


	state.ID = types.StringValue("placeholder")


	state.Token = types.StringValue(d.client.apiToken)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return 
	}
}