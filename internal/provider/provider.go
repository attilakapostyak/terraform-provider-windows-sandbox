package provider

import (
	"context"
	"fmt"
	"os"

	"terraform-provider-windows-sandbox/internal/wsb/configuration"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure WSBProvider satisfies various provider interfaces.
var _ provider.Provider = &WSBProvider{}
var _ provider.ProviderWithFunctions = &WSBProvider{}

// WSBProvider defines the provider implementation.
type WSBProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// WSBProviderModel describes the provider data model.
type WSBProviderModel struct {
	Path types.String `tfsdk:"path"`
}

func (p *WSBProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "wsb"
	resp.Version = p.version
}

func (p *WSBProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				MarkdownDescription: "Path specifying the location of Windows Sandbox configuration files.",
				Optional:            true,
			},
		},
	}
}

func (p *WSBProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data WSBProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var path = data.Path.ValueString()
	if len(path) > 0 {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			resp.Diagnostics.AddError(
				"Invalid path",
				fmt.Sprintf("The specified path '%s' does not exist.", path),
			)
		}
	}

	// Example client configuration for data sources and resources
	handler := configuration.DefaultWSBConfigurationHandler
	handler.DefaultPath = path
	resp.DataSourceData = handler
	resp.ResourceData = handler
}

func (p *WSBProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewConfigurationResource,
		NewLogonScriptResource,
	}
}

func (p *WSBProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewContextDataSource,
		NewConfigurationDataSource,
	}
}

func (p *WSBProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WSBProvider{
			version: version,
		}
	}
}
