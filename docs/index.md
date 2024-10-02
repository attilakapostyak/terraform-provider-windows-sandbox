
# Windows Sandbox Provider

This Terraform provider enables users to create and manage WSB files for defining Windows Sandbox configurations. The provider includes a data source for fetching existing configurations, a resource for creating new configurations, and a context data source to generate paths for specific folders in the sandbox environment.

## Table of Contents

- [Resource: wsb_configuration](#resource-wsb_configuration)
- [Resource: wsb_logon_script](#resource-wsb_logon_script)
- [Data Source: wsb_configuration](#data-source-wsb_configuration)
- [Data Source: wsb_context](#data-source-wsb_context)

---

## Provider Configuration

A basic provider configuration

```hcl
provider "wsb" {

}
```

A provider configuration where output path for .wsb file is specified
```hcl
provider "wsb" {
  path = "C:\\Windows\\Temp"
}
```



## Resource: `wsb_configuration`

The `wsb_configuration` resource is used to create a new Windows Sandbox configuration.

### Example Usage

The following example creates a Windows Sandbox configuration file located at `C:\Windows\Temp\example-configuration.wsb`. The Windows Sandbox instance will have the Downloads folder mapped to JohnDoe's Downloads folder.

```hcl
locals {
  mapped_folders = [
    {
      host_folder = "C:\\Users\\JohnDoe\\Downloads"
      sandbox_folder = "C:\\Users\\WDAGUtilityAccount\\Downloads"
      read_only = false
    }
  ]
}

resource "wsb_configuration" "example" {
  name        = "example-configuration"
  path        = "C:\\Windows\\Temp"
  virtual_gpu = false
  networking  = true
  audio_input = true
  video_input = false
  protected_client = false
  printer_redirection = true
  clipboard_redirection = true
  memory = 8192

  dynamic "mapped_folders" {
    for_each = local.mapped_folders
    content {
      host_folder      = mapped_folders.value.host_folder
      sandbox_folder  = mapped_folders.value.sandbox_folder
      read_only      = mapped_folders.value.read_only
    }
  }
}

```

## Resource: `wsb_logon_script`

This resource helps you build logon command scripts that run when Windows Sandbox starts up.
The resource has the following capabilities:
1. Specify a script to be executed before package installation
1. Install packages using winget
1. Install packages using scoop
1. Specify a script to be executed after package installation


### Example Usage
```hcl
resource "wsb_logon_script" "example_configuration" {
  name            = "example-configuration"
  scoop_packages  = ["vcredist2022"]
  winget_packages = ["nvim"]

  pre_installation_scripts  = "Write-Host 'Something to perform before software installation. Add your own Powershell script here'"
  post_installation_scripts = "Write-Host 'Something to perform after software installation. Add your own Powershell script here'"
}
```

## Data Source: `wsb_configuration`

The `wsb_configuration` data source allows you to retrieve existing Windows Sandbox configurations.

### Example Usage

The following data source will open a Windows Sandbox configuration file named `example_wsb.wsb` file located in the default path specified in the provider configuration.

```hcl
data "wsb_configuration" "example" {
  name = "example_wsb"
}
```


## Data Source: `wsb_context`

The `wsb_context` data source provides contextual paths for various folders in the Windows Sandbox environment.

### Example Usage

```hcl
data "wsb_context" "example" {

}
```

## Practical Example

The following Terraform script will create a Windows Sandbox configuration file for each of your team members, with the correct mapped folders for each member of your team.

```hcl
locals {
  team_members = ["Sofia", "John", "Maria", "Trevor", "Carl", "Heather"]

  # Dynamically create a mapped folder structure for each team member
  team_mapped_folders = { 
    for member in local.team_members : 
    member => [
      {
        host_folder    = data.wsb_context.example[member].user_downloads_folder
        sandbox_folder = data.wsb_context.example[member].sandbox_container_user_downloads_folder
        read_only      = false
      }
    ]
  }
}

data "wsb_context" "example" {
  for_each           = toset(local.team_members)
  username              = each.value
  users_folder          = "C:\\Users" 
  downloads_folder_name = "Downloads" 
}

# Create a separate wsb_configuration for each team member
resource "wsb_configuration" "team" {
  for_each           = toset(local.team_members)
  name               = "playground-${lower(each.value)}"
  path               = "./"
  virtual_gpu        = false
  networking         = true
  audio_input        = true
  video_input        = true
  protected_client   = true
  printer_redirection = true
  clipboard_redirection = true
  memory = 8192

  dynamic "mapped_folders" {
    for_each = local.team_mapped_folders[each.value]
    content {
      host_folder    = mapped_folders.value.host_folder
      sandbox_folder = mapped_folders.value.sandbox_folder
      read_only      = mapped_folders.value.read_only
    }
  }
}
```