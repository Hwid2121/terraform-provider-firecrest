package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// const baseURL = "https://firecrest.cscs.ch".
// const baseURL = "/auth/realms/firecrest-clients/protocol/openid-connect/token"

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &firecrestProvider{}

// var _ provider.ProviderWithFunctions = &ScaffoldingProvider{}

type firecrestProvider struct {
	client  *FirecrestClient
	version string
}

func (p *firecrestProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "firecrest"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *firecrestProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Firecrest provider allows you to interact with the Firecrest API.",
	}
}

func (p *firecrestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Firecrest client")
	client := NewFireCrestClient("", "")
	p.client = client
	resp.DataSourceData = p
	resp.ResourceData = p
	tflog.Info(ctx, "Configured Firecrest client successfully")
}

func (p *firecrestProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFirecrestJobResource,
	}
}

func (p *firecrestProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFireDataSource,
	}
}

func (p *firecrestProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &firecrestProvider{
			version: version,
		}
	}
}
