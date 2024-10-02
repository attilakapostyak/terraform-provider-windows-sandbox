package provider

import (
	"context"
	"os/user"
	"path"
	"strings"

	"terraform-provider-windows-sandbox/internal/tools"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ContextDataSource{}

func NewContextDataSource() datasource.DataSource {
	return &ContextDataSource{}
}

// ContextDataSource defines the data source implementation.
type ContextDataSource struct {
}

// ContextDataSourceModel describes the data source data model.
type ContextDataSourceModel struct {
	UsersFolder                         types.String `tfsdk:"users_folder"`
	UserName                            types.String `tfsdk:"username"`
	DownloadsFolderName                 types.String `tfsdk:"downloads_folder_name"`
	UserProfileFolder                   types.String `tfsdk:"user_profile_folder"`
	UserDownloadsFolder                 types.String `tfsdk:"user_downloads_folder"`
	SandboxUsersFolder                  types.String `tfsdk:"sandbox_users_folder"`
	SandboxContainerUserName            types.String `tfsdk:"sandbox_container_username"`
	SandboxContainerUserProfileFolder   types.String `tfsdk:"sandbox_container_user_profile_folder"`
	SandboxContainerUserDownloadsFolder types.String `tfsdk:"sandbox_container_user_downloads_folder"`
}

func (d *ContextDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_context"
}

func (d *ContextDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Context data source",

		Attributes: map[string]schema.Attribute{
			"users_folder": schema.StringAttribute{
				MarkdownDescription: "Path to the users directory (e.g., C:\\Users)",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for which the configuration is applicable",
				Optional:            true,
			},
			"downloads_folder_name": schema.StringAttribute{
				MarkdownDescription: "Name of the Downloads folder (Default: Downloads)",
				Optional:            true,
			},
			"user_profile_folder": schema.StringAttribute{
				MarkdownDescription: "Path to the user's profile directory (e.g., C:\\Users\\john)",
				Computed:            true,
			},
			"user_downloads_folder": schema.StringAttribute{
				MarkdownDescription: "Path to the user's Downloads directory (e.g., C:\\Users\\john\\Downloads)",
				Computed:            true,
			},
			"sandbox_users_folder": schema.StringAttribute{
				MarkdownDescription: "Path to the sandbox users directory",
				Computed:            true,
			},
			"sandbox_container_username": schema.StringAttribute{
				MarkdownDescription: "Username for the sandbox container environment",
				Computed:            true,
			},
			"sandbox_container_user_profile_folder": schema.StringAttribute{
				MarkdownDescription: "Path to the Sandbox container user profile directory",
				Computed:            true,
			},
			"sandbox_container_user_downloads_folder": schema.StringAttribute{
				MarkdownDescription: "Path to the Sandbox container account Downloads directory",
				Computed:            true,
			},
		},
	}
}

func (d *ContextDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

}

func (d *ContextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContextDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := user.Current()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving current user", err.Error())
	}
	username_parts := strings.Split(user.Username, "\\")

	if len(data.UsersFolder.ValueString()) == 0 {
		data.UsersFolder = types.StringValue("C:\\Users")
	}

	if len(data.UserName.ValueString()) == 0 {
		data.UserName = types.StringValue(username_parts[len(username_parts)-1])
	}

	if len(data.DownloadsFolderName.ValueString()) == 0 {
		data.DownloadsFolderName = types.StringValue("Downloads")
	}

	data.UserProfileFolder = types.StringValue(
		tools.ConvertToWindowsPath(
			path.Join(data.UsersFolder.ValueString(), data.UserName.ValueString())))
	data.UserDownloadsFolder = types.StringValue(
		tools.ConvertToWindowsPath(
			path.Join(data.UserProfileFolder.ValueString(), data.DownloadsFolderName.ValueString())))
	data.SandboxUsersFolder = types.StringValue("C:\\Users")
	data.SandboxContainerUserName = types.StringValue("WDAGUtilityAccount")
	data.SandboxContainerUserProfileFolder = types.StringValue(
		tools.ConvertToWindowsPath(
			path.Join(data.SandboxUsersFolder.ValueString(), data.SandboxContainerUserName.ValueString())))
	data.SandboxContainerUserDownloadsFolder = types.StringValue(
		tools.ConvertToWindowsPath(
			path.Join(data.SandboxContainerUserProfileFolder.ValueString(), data.DownloadsFolderName.ValueString())))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
