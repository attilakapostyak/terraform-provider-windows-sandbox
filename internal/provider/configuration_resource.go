package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"terraform-provider-windows-sandbox/internal/wsb/configuration"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ConfigurationResource{}
var _ resource.ResourceWithImportState = &ConfigurationResource{}

func NewConfigurationResource() resource.Resource {
	return &ConfigurationResource{}
}

type ConfigurationNotExistsError struct {
	Err error
}

func (c *ConfigurationNotExistsError) Error() string {
	return c.Err.Error()
}

// ConfigurationResource defines the resource implementation.
type ConfigurationResource struct {
	handler *configuration.WSBConfigurationHandler
}

// ConfigurationResourceModel describes the resource data model.
type ConfigurationResourceModel struct {
	Name                 types.String                     `tfsdk:"name"`
	Path                 types.String                     `tfsdk:"path"`
	VirtualGpu           types.Bool                       `tfsdk:"virtual_gpu"`
	Networking           types.Bool                       `tfsdk:"networking"`
	LogonCommand         types.String                     `tfsdk:"logon_command"`
	AudioInput           types.Bool                       `tfsdk:"audio_input"`
	VideoInput           types.Bool                       `tfsdk:"video_input"`
	ProtectedClient      types.Bool                       `tfsdk:"protected_client"`
	PrinterRedirection   types.Bool                       `tfsdk:"printer_redirection"`
	ClipboardRedirection types.Bool                       `tfsdk:"clipboard_redirection"`
	Memory               types.String                     `tfsdk:"memory"`
	MappedFolders        []ConfigurationMappedFolderModel `tfsdk:"mapped_folders"`
}

func (r *ConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration"
}

func (r *ConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Configuration resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the configuration",
				Required:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path specifying the location of Windows Sandbox configuration files. Overrides the Path configured in the provider if specified.",
				Optional:            true,
			},
			"networking": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables networking in the sandbox. You can disable network access to decrease the attack surface exposed by the sandbox.",
				Optional:            true,
			},
			"virtual_gpu": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables GPU sharing.",
				Optional:            true,
			},
			"logon_command": schema.StringAttribute{
				MarkdownDescription: "Command to execute when a user logs on to the sandbox. The command should be PowerShell scripts.",
				Optional:            true,
			},
			"audio_input": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables audio input to the sandbox.",
				Optional:            true,
			},
			"video_input": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables video input to the sandbox.",
				Optional:            true,
			},
			"protected_client": schema.BoolAttribute{
				MarkdownDescription: "When Protected Client mode is enabled, Sandbox adds a new layer of security boundary by running inside an AppContainer Isolation execution environment.",
				Optional:            true,
			},
			"printer_redirection": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables printer sharing from the host into the sandbox.",
				Optional:            true,
			},
			"clipboard_redirection": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables sharing of the host clipboard with the sandbox.",
				Optional:            true,
			},
			"memory": schema.StringAttribute{
				MarkdownDescription: "Specifies the amount of memory that the sandbox can use in megabytes (MB).",
				Optional:            true,
			},
		},

		Blocks: map[string]schema.Block{
			"mapped_folders": schema.ListNestedBlock{
				MarkdownDescription: "An array of folders, each representing a location on the host machine that is shared with the sandbox at the specified path.",
				NestedObject: schema.NestedBlockObject{
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

func (r *ConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.handler = handler
}

func (r *ConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConfigurationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := ValidateConfigurationResourceModel(&data); err != nil {
		resp.Diagnostics.AddError("Validation Error", fmt.Sprintf("validation failed: %s", err.Error()))
		return
	}

	err := SaveConfigurationResourceModel(r.handler, &data)
	if err != nil {
		resp.Diagnostics.AddError("Save error", fmt.Sprintf("Unable to create configuration: %s", err.Error()))
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConfigurationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := ValidateConfigurationResourceModel(&data); err != nil {
		resp.Diagnostics.AddError("Validation Error", fmt.Sprintf("validation failed: %s", err.Error()))
		return
	}

	err := LoadConfigurationResourceModel(r.handler, &data)
	if err != nil {
		resp.State.RemoveResource(ctx)

		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConfigurationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := ValidateConfigurationResourceModel(&data); err != nil {
		resp.Diagnostics.AddError("Validation Error", fmt.Sprintf("validation failed: %s", err.Error()))
		return
	}

	err := SaveConfigurationResourceModel(r.handler, &data)
	if err != nil {
		resp.Diagnostics.AddError("Save error", fmt.Sprintf("Unable to update configuration: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConfigurationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	err := DeleteConfiguration(r.handler, &data)
	if err != nil {
		resp.Diagnostics.AddWarning("Delete error", fmt.Sprintf("Unable to delete configuration: %s", err))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data ConfigurationResourceModel

	data.Path = types.StringValue(filepath.Dir(req.ID))
	var filename = filepath.Base(req.ID)
	filename = strings.Replace(filename, ".wsb", "", -1)

	data.Name = types.StringValue(filename)

	resp.Diagnostics.AddWarning("DEBUG", fmt.Sprintf("DATA %+v\n", data))

	err := LoadConfigurationResourceModel(r.handler, &data)
	if err != nil {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("Load error", err.Error())

		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func ValidateConfigurationResourceModel(data *ConfigurationResourceModel) error {

	sandboxFolderMap := make(map[string]bool)
	for _, mappedFolder := range data.MappedFolders {
		if sandboxFolderMap[mappedFolder.SandboxFolder.ValueString()] {
			return fmt.Errorf("duplicate SandboxFolder found: %s", mappedFolder.SandboxFolder)
		}
		sandboxFolderMap[mappedFolder.SandboxFolder.ValueString()] = true
	}

	return nil
}

func LoadConfigurationResourceModel(handler *configuration.WSBConfigurationHandler, data *ConfigurationResourceModel) error {
	path, err := configuration.ResolvePath(handler.DefaultPath, data.Path.ValueString())
	if err != nil {
		return fmt.Errorf("Error resolving path: \n%w", err)
	}

	fileName := handler.BuildConfigurationFileName(path, data.Name.ValueString())

	config, load_err := handler.LoadConfiguration(fileName)

	if config != nil {

		value, err := configuration.ParseFeatureFlag(config.VGPU)
		if err != nil {
			return fmt.Errorf("Error reading configuration \n%w", err)
		}
		data.VirtualGpu = value.ToBool(VirtualGpuDefaultValue)

		value, err = configuration.ParseFeatureFlag(config.Networking)
		if err != nil {
			return fmt.Errorf("Error reading configuration \n%w", err)
		}
		data.Networking = value.ToBool(NetworkingDefaultValue)

		value, err = configuration.ParseFeatureFlag(config.AudioInput)
		if err != nil {
			return fmt.Errorf("Error reading configuration \n%w", err)
		}
		data.AudioInput = value.ToBool(AudioInputDefaultValue)

		value, err = configuration.ParseFeatureFlag(config.VideoInput)
		if err != nil {
			return fmt.Errorf("Error reading configuration \n%w", err)
		}
		data.VideoInput = value.ToBool(VideoInputDefaultValue)

		value, err = configuration.ParseFeatureFlag(config.ProtectedClient)
		if err != nil {
			return fmt.Errorf("Error reading configuration \n%w", err)
		}
		data.ProtectedClient = value.ToBool(ProtectedClientDefaultValue)

		value, err = configuration.ParseFeatureFlag(config.PrinterRedirection)
		if err != nil {
			return fmt.Errorf("Error reading configuration \n%w", err)
		}
		data.PrinterRedirection = value.ToBool(PrinterRedirectionDefaultValue)

		value, err = configuration.ParseFeatureFlag(config.ClipboardRedirection)
		if err != nil {
			return fmt.Errorf("Error reading configuration \n%w", err)
		}
		data.ClipboardRedirection = value.ToBool(ClipboardRedirectionDefaultValue)

		switch len(config.Memory) {
		case 0:
			data.Memory = types.StringNull()
		default:
			data.Memory = types.StringValue(config.Memory)
		}

		switch len(config.LogonCommand.Command) {
		case 0:
			data.LogonCommand = types.StringNull()
		default:
			data.LogonCommand = types.StringValue(config.LogonCommand.Command)
		}

		data.MappedFolders = make([]ConfigurationMappedFolderModel, 0)
		for _, mappedFolder := range config.MappedFolders.MappedFolder {
			readonly, err := strconv.ParseBool(mappedFolder.ReadOnly)
			if err != nil {
				return fmt.Errorf("Error parsing read_only flag for host folder: %s\n%w", mappedFolder.HostFolder, err)
			}
			data.MappedFolders = append(data.MappedFolders, ConfigurationMappedFolderModel{
				HostFolder:    types.StringValue(mappedFolder.HostFolder),
				SandboxFolder: types.StringValue(mappedFolder.SandboxFolder),
				ReadOnly:      types.BoolValue(readonly),
			})
		}
	}

	if load_err != nil {
		new_data := ConfigurationResourceModel{}
		*data = new_data

		return &ConfigurationNotExistsError{
			Err: load_err,
		}
	}

	return nil
}

func SaveConfigurationResourceModel(handler *configuration.WSBConfigurationHandler, data *ConfigurationResourceModel) error {
	path, err := configuration.ResolvePath(handler.DefaultPath, data.Path.ValueString())
	if err != nil {
		return fmt.Errorf("error resolving path: %s", err.Error())
	}

	fileName := handler.BuildConfigurationFileName(path, data.Name.ValueString())

	config := *configuration.NewWSBConfiguration()

	if !data.VirtualGpu.IsNull() {
		config.VGPU = configuration.ToFeatureFlag(data.VirtualGpu.ValueBool()).String()
	}
	if !data.Networking.IsNull() {
		config.Networking = configuration.ToFeatureFlag(data.Networking.ValueBool()).String()
	}
	if !data.AudioInput.IsNull() {
		config.AudioInput = configuration.ToFeatureFlag(data.AudioInput.ValueBool()).String()
	}
	if !data.VideoInput.IsNull() {
		config.VideoInput = configuration.ToFeatureFlag(data.VideoInput.ValueBool()).String()
	}
	if !data.ProtectedClient.IsNull() {
		config.ProtectedClient = configuration.ToFeatureFlag(data.ProtectedClient.ValueBool()).String()
	}
	if !data.PrinterRedirection.IsNull() {
		config.PrinterRedirection = configuration.ToFeatureFlag(data.PrinterRedirection.ValueBool()).String()
	}
	if !data.ClipboardRedirection.IsNull() {
		config.ClipboardRedirection = configuration.ToFeatureFlag(data.ClipboardRedirection.ValueBool()).String()
	}
	if !data.Memory.IsNull() {
		config.Memory = data.Memory.ValueString()
	}
	if !data.LogonCommand.IsNull() {
		config.LogonCommand.Command = data.LogonCommand.ValueString()
	}

	if len(data.MappedFolders) > 0 {
		var mappedFolders = make([]configuration.MappedFolder, len(data.MappedFolders))
		for i, mappedFolder := range data.MappedFolders {
			mappedFolders[i] = configuration.MappedFolder{
				HostFolder:    mappedFolder.HostFolder.ValueString(),
				SandboxFolder: mappedFolder.SandboxFolder.ValueString(),
				ReadOnly:      strconv.FormatBool(mappedFolder.ReadOnly.ValueBool()),
			}
		}
		config.MappedFolders.MappedFolder = mappedFolders
	}

	err = handler.SaveConfiguration(fileName, &config)
	if err != nil {
		return fmt.Errorf("error reading configuration.\n%s", err.Error())
	}

	return nil
}

func DeleteConfiguration(handler *configuration.WSBConfigurationHandler, data *ConfigurationResourceModel) error {
	path, err := configuration.ResolvePath(handler.DefaultPath, data.Path.ValueString())
	if err != nil {
		return fmt.Errorf("error resolving path: %s", err.Error())
	}

	fileName := handler.BuildConfigurationFileName(path, data.Name.ValueString())

	err = handler.DeleteConfiguration(fileName)
	if err != nil {
		return fmt.Errorf("error deleting configuration.\n%s", err.Error())
	}

	return nil
}
