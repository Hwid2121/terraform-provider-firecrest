package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
	// "github.com/hashicorp/terraform-svchost/disco"
)

var (
	_ resource.Resource = &firecrestFileResource{}
	_ resource.ResourceWithConfigure = &firecrestFileResource{}
)

func NewFirecrestFileResource() resource.Resource {
	return &firecrestFileResource{}
}

type firecrestFileResource struct {
	client *FirecrestClient
}


type firecrestFileResourceModel struct {
	ID	types.String `tfsdk:"id"`
	SourcePath types.String `tfsdk:"source_path"`
	DestinationPath types.String `tfsdk:"destination_path"`
}

func (r *firecrestFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *firecrestFileResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"source_path": schema.StringAttribute{
				Required: true,
			},
			"destination_path": schema.StringAttribute{
				Required: true,
			},
		},
	}
}


func (r *firecrestFileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp * resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerConfig, ok := req.ProviderData.(*firecrestProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected provider to be of type *firecrestProvider, got %T", req.ProviderData),
		)
		return
	}
	r.client = providerConfig.client

}

func (r *firecrestFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan firecrestFileResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	
}

func (r *firecrestFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state firecrestFileResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TOOD implement read ogic to verify file status
}

func (r *firecrestFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *firecrestFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}