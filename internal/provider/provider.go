package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// const baseURL = "https://firecrest.cscs.ch".
const baseURL = "/auth/realms/firecrest-clients/protocol/openid-connect/token"

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &firecrestProvider{}

// var _ provider.ProviderWithFunctions = &ScaffoldingProvider{}

// ScaffoldingProvider defines the provider implementation.
type firecrestProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	client  *FirecrestClient
	version string
}

// type firecrestProviderModel describes the provider data model.
type firecrestProviderModel struct {
	// Endpoint types.String `tfsdk:"endpoint"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *firecrestProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "firecrest"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *firecrestProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Firecrest provider allows you to interact with the Firecrest API.",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "Client ID for firecREST API. Provided by https://oidc-dashboard-prod.cscs.ch/",
				Required:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Client Secret for firecREST API. Provided by https://oidc-dashboard-prod.cscs.ch/",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *firecrestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring FirecrREST client")
	var config firecrestProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Provider Configuration Error", "Unable to parse provider configuration.")
		return
	}

	clientID := config.ClientID.ValueString()
	clientSecret := config.ClientSecret.ValueString()

	// resp.Diagnostics.AddWarning("ci sono", clientSecret)

	if clientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("clientID"),
			"Missing firecREST clientID",
			"The provider cannot create the firecREST API client as there is a missing or empty value for the firecREST API clientID. "+
				"Set the clientID value in the configuration."+
				"Ensure the value is not empty.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("clientSecret"),
			"Missing firecREST clientSecret",
			"The provider cannot create the firecREST API client as there is a missing or empty value for the firecREST API clientSecret. "+
				"Set the clientSecret value in the configuration."+
				"Ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "clientID", clientID)
	ctx = tflog.SetField(ctx, "client_secret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_secret")
	tflog.Debug(ctx, "Creating FirecREST Client")

	client := NewFireCrestClient(baseURL, "")
	token, err := client.GetToken(clientID, clientSecret)
	if err != nil {
		resp.Diagnostics.AddError("Failed to retrieve token", err.Error())
		return
	}

	// client.apiToken = token
	client.SetToken(token)
	p.client = client

	resp.DataSourceData = p
	resp.ResourceData = p

	tflog.Debug(ctx, "API Token"+token)
	tflog.Info(ctx, "Configured FirecREST client successfully! ")

	// tflog.Debug(ctx, fmt.Sprintf("+++ DataSourceData type: %T", p))
	// log.Println("DataSourceData type:", fmt.Sprintf("%T", p))

}

func (p *firecrestProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFirecrestJobResource,
	}
	// return nil
}

func (p *firecrestProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFireDataSource,
	}
}

func (p *firecrestProvider) Functions(ctx context.Context) []func() function.Function {
	// return []func() function.Function{
	// 	NewExampleFunction,
	// }
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &firecrestProvider{
			version: version,
		}
	}
}
