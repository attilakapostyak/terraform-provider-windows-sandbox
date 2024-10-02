package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &LogonScriptResource{}
var _ resource.ResourceWithImportState = &LogonScriptResource{}

func NewLogonScriptResource() resource.Resource {
	return &LogonScriptResource{}
}

// LogonScriptResource defines the resource implementation.
type LogonScriptResource struct {
}

// LogonScriptResourceModel describes the resource data model.
type LogonScriptResourceModel struct {
	Name                    types.String   `tfsdk:"name"`
	LogonCommandScript      types.String   `tfsdk:"logon_command_script"`
	PreInstallationScripts  types.String   `tfsdk:"pre_installation_scripts"`
	PostInstallationScripts types.String   `tfsdk:"post_installation_scripts"`
	WingetPackages          []types.String `tfsdk:"winget_packages"`
	ScoopPackages           []types.String `tfsdk:"scoop_packages"`
}

func (r *LogonScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_logon_script"
}

func (r *LogonScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "LogonScript resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the configuration",
				Required:            true,
			},
			"logon_command_script": schema.StringAttribute{
				MarkdownDescription: "The script built automatically by Terraform and execuuted when the sandbox starts.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pre_installation_scripts": schema.StringAttribute{
				MarkdownDescription: "Scripts to be executed before the installation of packages starts. The scripts should be PowerShell scripts.",
				Optional:            true,
			},
			"post_installation_scripts": schema.StringAttribute{
				MarkdownDescription: "Scripts to be executed after the installation of packages is completed. The scripts should be PowerShell scripts.",
				Optional:            true,
			},
			"winget_packages": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of packages to be installed in the sandbox using winget package manager.",
				Optional:            true,
			},
			"scoop_packages": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of packages to be installed in the sandbox scoop package manager.",
				Optional:            true,
			},
		},
	}
}

func (r *LogonScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (r *LogonScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LogonScriptResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	script := types.StringValue(BuildLogonCommandScript(&data))
	data.LogonCommandScript = script

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LogonScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LogonScriptResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	script := types.StringValue(BuildLogonCommandScript(&data))
	if data.LogonCommandScript != script {
		data.PreInstallationScripts = types.StringNull()
		data.PostInstallationScripts = types.StringNull()
		data.WingetPackages = make([]types.String, 0)
		data.WingetPackages = make([]types.String, 0)
	}
	data.LogonCommandScript = script

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LogonScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LogonScriptResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	script := types.StringValue(BuildLogonCommandScript(&data))
	if data.LogonCommandScript != script {
		data.PreInstallationScripts = types.StringNull()
		data.PostInstallationScripts = types.StringNull()
		data.WingetPackages = make([]types.String, 0)
		data.WingetPackages = make([]types.String, 0)
	}
	data.LogonCommandScript = script

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LogonScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LogonScriptResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *LogonScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func BuildLogonCommandScript(data *LogonScriptResourceModel) string {
	baseScript := "powershell -ExecutionPolicy Unrestricted -Command \"# Generated with Terraform wsb_logon_script." + data.Name.ValueString() + "  \nstart powershell {-noexit -command \" [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String('%s')) | Out-File %s ; . %s \" } \" "

	var sb strings.Builder

	// Install Winget
	if len(data.WingetPackages) > 0 {
		sb.WriteString("Write-Host 'Install Winget ...'\n")
		sb.WriteString("$progressPreference = 'silentlyContinue'\n")
		sb.WriteString("Write-Information \"Downloading WinGet and its dependencies...\"\n")
		sb.WriteString("Invoke-WebRequest -Uri https://aka.ms/getwinget -OutFile Microsoft.DesktopAppInstaller_8wekyb3d8bbwe.msixbundle\n")
		sb.WriteString("Invoke-WebRequest -Uri https://aka.ms/Microsoft.VCLibs.x64.14.00.Desktop.appx -OutFile Microsoft.VCLibs.x64.14.00.Desktop.appx\n")
		sb.WriteString("Invoke-WebRequest -Uri https://github.com/microsoft/microsoft-ui-xaml/releases/download/v2.8.6/Microsoft.UI.Xaml.2.8.x64.appx -OutFile Microsoft.UI.Xaml.2.8.x64.appx\n")
		sb.WriteString("Add-AppxPackage Microsoft.VCLibs.x64.14.00.Desktop.appx\n")
		sb.WriteString("Add-AppxPackage Microsoft.UI.Xaml.2.8.x64.appx\n")
		sb.WriteString("Add-AppxPackage Microsoft.DesktopAppInstaller_8wekyb3d8bbwe.msixbundle\n")
		sb.WriteString("\n\n")
	}

	// Install Scoop
	if len(data.ScoopPackages) > 0 {
		sb.WriteString("Write-Host 'Install Scoop ...'\n")
		sb.WriteString("# Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser\n")
		sb.WriteString("Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression\n")
		sb.WriteString("scoop install git\n")
		sb.WriteString("scoop bucket add extras\n")
		sb.WriteString("\n\n")
	}

	if len(data.PreInstallationScripts.ValueString()) > 0 {
		sb.WriteString("Write-Host 'Run Pre-installation scripts ...'\n")
		sb.WriteString(data.PreInstallationScripts.ValueString())
		sb.WriteString("\n\n")
	}

	// Install packages from Winget
	if len(data.WingetPackages) > 0 {
		var packageCommand string
		sb.WriteString("Write-Host 'Winget packages installation ...'\n")
		packageCommand = "winget install %s --accept-source-agreements --accept-package-agreements \n"
		for _, script := range data.WingetPackages {
			sb.WriteString(fmt.Sprintf(packageCommand, script.ValueString()))
		}
		sb.WriteString("\n\n")
	}

	// Install packages from Scoop
	if len(data.ScoopPackages) > 0 {
		var packageCommand string
		sb.WriteString("Write-Host 'Scoop packages installation ...'\n")
		packageCommand = "scoop install %s \n"
		for _, script := range data.ScoopPackages {
			sb.WriteString(fmt.Sprintf(packageCommand, script.ValueString()))
		}
		sb.WriteString("\n\n")
	}

	if len(data.PostInstallationScripts.ValueString()) > 0 {
		sb.WriteString("Write-Host 'Run Post-installation scripts ...'\n")
		sb.WriteString(data.PostInstallationScripts.ValueString())
		sb.WriteString("\n\n")
	}

	var logonCommandScript = "C:\\Windows\\Temp\\longoncommand_script.ps1"
	return fmt.Sprintf(baseScript, base64.StdEncoding.EncodeToString([]byte(sb.String())), logonCommandScript, logonCommandScript)
}
