# Resource `windows-sandbox_logon_script`

The `logon_script` resource allows you to manage Windows Sandbox wsb configuration files by defining logon scripts and package installations. This resource is particularly useful when used in conjunction with the `windows-sandbox_configuration` resource, as its attributes can be leveraged to create mapped folders in the Windows Sandbox environment. By utilizing the `logon_script` resource, you can automate the setup of your sandbox environment, including installing packages and running custom scripts, while also providing the necessary content for mapped folders.


## Argument Reference

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| name | String |  | (Required) The name of the configuration. |
| pre_installation_scripts | String | "" | (Optional) Scripts to be executed before the installation of packages starts. |
| post_installation_scripts | String | "" | (Optional) Scripts to be executed after the installation of packages is completed. |
| winget_packages | List<String> | [] | (Optional) List of packages to be installed in the sandbox using Winget package manager. |
| scoop_packages | List<String> | [] | (Optional) List of packages to be installed in the sandbox Scoop package manager. |

## Attributes Reference

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| logon_command_script | String | N/A | The script built automatically by Terraform and executed when the sandbox starts. |

## Example Usage

```hcl
resource "windows-sandbox_logon_script" "example" {
  name = "example"

  pre_installation_scripts = <<-EOT
    Write-Host 'Run Pre-installation scripts ...'
    # Your pre-installation script here
  EOT

  post_installation_scripts = <<-EOT
    Write-Host 'Run Post-installation scripts ...'
    # Your post-installation script here
  EOT

  winget_packages = [
    "package1",
    "package2",
  ]

  scoop_packages = [
    "scoop-package1",
    "scoop-package2",
  ]
}
```