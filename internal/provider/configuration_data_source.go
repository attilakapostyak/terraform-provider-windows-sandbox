package provider

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-windows-sandbox/internal/wsb/configuration"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const VirtualGpuDefaultValue = true
const NetworkingDefaultValue = true
const AudioInputDefaultValue = true
const VideoInputDefaultValue = true
const ProtectedClientDefaultValue = true
const PrinterRedirectionDefaultValue = true
const ClipboardRedirectionDefaultValue = true

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ConfigurationDataSource{}

func NewConfigurationDataSource() datasource.DataSource {
	return &ConfigurationDataSource{}
}

// ConfigurationDataSource defines the data source implementation.
type ConfigurationDataSource struct {
	handler *configuration.WSBConfigurationHandler
}

// ConfigurationDataSourceModel describes the data source data model.
type ConfigurationDataSourceModel struct {
	Name                 types.String                     `tfsdk:"name"`
	Path                 types.String                     `tfsdk:"path"`
	VirtualGpu           types.Bool                       `tfsdk:"virtual_gpu"`
	Networking           types.Bool                       `tfsdk:"network"`
	LogonCommand         types.String                     `tfsdk:"logon_command"`
	AudioInput           types.Bool                       `tfsdk:"audio_input"`
	VideoInput           types.Bool                       `tfsdk:"video_input"`
	ProtectedClient      types.Bool                       `tfsdk:"protected_client"`
	PrinterRedirection   types.Bool                       `tfsdk:"printer_redirection"`
	ClipboardRedirection types.Bool                       `tfsdk:"clipboard_redirection"`
	Memory               types.String                     `tfsdk:"memory"`
	MappedFolders        []ConfigurationMappedFolderModel `tfsdk:"mapped_folders"`
}

type ConfigurationMappedFolderModel struct {
	HostFolder    types.String `tfsdk:"host_folder"`
	SandboxFolder types.String `tfsdk:"sandbox_folder"`
	ReadOnly      types.Bool   `tfsdk:"read_only"`
}

func (d *ConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration"
}

func (d *ConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Configuration data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the configuration",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path specifying the location of Windows Sandbox configuration files. Overrides the Path configured in the provider if specified.",
				Optional:            true,
			},
			"network": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables networking in the sandbox. You can disable network access to decrease the attack surface exposed by the sandbox.",
				Computed:            true,
			},
			"virtual_gpu": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables GPU sharing.",
				Computed:            true,
			},
			"logon_command": schema.StringAttribute{
				MarkdownDescription: "Specifies a single command that will be invoked automatically after the sandbox logs on. Apps in the sandbox are run under the container user account. The container user account should be an administrator account.",
				Computed:            true,
			},
			"audio_input": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables audio input to the sandbox.",
				Computed:            true,
			},
			"video_input": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables video input to the sandbox.",
				Computed:            true,
			},
			"protected_client": schema.BoolAttribute{
				MarkdownDescription: "When Protected Client mode is enabled, Sandbox adds a new layer of security boundary by running inside an AppContainer Isolation execution environment.",
				Computed:            true,
			},
			"printer_redirection": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables printer sharing from the host into the sandbox.",
				Computed:            true,
			},
			"clipboard_redirection": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables sharing of the host clipboard with the sandbox.",
				Computed:            true,
			},
			"memory": schema.StringAttribute{
				MarkdownDescription: "Specifies the amount of memory that the sandbox can use in megabytes (MB).",
				Computed:            true,
			},
			"mapped_folders": schema.SetNestedAttribute{
				MarkdownDescription: "An array of folders, each representing a location on the host machine that is shared with the sandbox at the specified path.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host_folder": schema.StringAttribute{
							Required: true,
						},
						"sandbox_folder": schema.StringAttribute{
							Optional: true,
						},
						"read_only": schema.BoolAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (d *ConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	handler, ok := req.ProviderData.(*configuration.WSBConfigurationHandler)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *WSBConfigurationHandler, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.handler = handler
}

func (d *ConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConfigurationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	path, err := configuration.ResolvePath(d.handler.DefaultPath, data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error resolving path", err.Error())
		return
	}

	fileName := d.handler.BuildConfigurationFileName(path, data.Name.ValueString())
	tflog.Trace(ctx, fmt.Sprintf("Try to read configuration file: %s", fileName))

	config, err := d.handler.LoadConfiguration(fileName)
	if err != nil {
		resp.Diagnostics.AddError("Error reading configuration", err.Error())
		return
	}

	value, err := configuration.ParseFeatureFlag(config.VGPU)
	if err != nil {
		resp.Diagnostics.AddError("Error reading configuration", err.Error())
	}
	data.VirtualGpu = value.ToBool(VirtualGpuDefaultValue)

	value, err = configuration.ParseFeatureFlag(config.Networking)
	if err != nil {
		resp.Diagnostics.AddError("Error reading configuration", err.Error())
	}
	data.Networking = value.ToBool(NetworkingDefaultValue)

	data.LogonCommand = types.StringValue(config.LogonCommand.Command)

	value, err = configuration.ParseFeatureFlag(config.AudioInput)
	if err != nil {
		resp.Diagnostics.AddError("Error reading configuration", err.Error())
	}
	data.AudioInput = value.ToBool(AudioInputDefaultValue)

	value, err = configuration.ParseFeatureFlag(config.VideoInput)
	if err != nil {
		resp.Diagnostics.AddError("Error reading configuration", err.Error())
	}
	data.VideoInput = value.ToBool(VideoInputDefaultValue)

	value, err = configuration.ParseFeatureFlag(config.ProtectedClient)
	if err != nil {
		resp.Diagnostics.AddError("Error reading configuration", err.Error())
	}
	data.ProtectedClient = value.ToBool(ProtectedClientDefaultValue)

	value, err = configuration.ParseFeatureFlag(config.PrinterRedirection)
	if err != nil {
		resp.Diagnostics.AddError("Error reading configuration", err.Error())
	}
	data.PrinterRedirection = value.ToBool(PrinterRedirectionDefaultValue)

	value, err = configuration.ParseFeatureFlag(config.ClipboardRedirection)
	if err != nil {
		resp.Diagnostics.AddError("Error reading configuration", err.Error())
	}
	data.ClipboardRedirection = value.ToBool(ClipboardRedirectionDefaultValue)

	data.Memory = types.StringValue(config.Memory)

	for _, mappedFolder := range config.MappedFolders.MappedFolder {
		readonly, err := strconv.ParseBool(mappedFolder.ReadOnly)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error parsing read_only flag for host folder: %s", mappedFolder.HostFolder), err.Error())
		}
		data.MappedFolders = append(data.MappedFolders, ConfigurationMappedFolderModel{
			HostFolder:    types.StringValue(mappedFolder.HostFolder),
			SandboxFolder: types.StringValue(mappedFolder.SandboxFolder),
			ReadOnly:      types.BoolValue(readonly),
		})
	}

	tflog.Trace(ctx, fmt.Sprintf("Done reading configuration file: %s", fileName))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
