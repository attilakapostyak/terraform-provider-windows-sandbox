# Data Source: `wsb_configuration`

This data source retrieves configuration details for a Windows Sandbox instance based on a provided configuration file.

## Arguments

| Name                 | Type    | Description |
|----------------------|---------|-------------|
| `name`               | String  | (Required) Name of the configuration. |
| `path`               | String  | (Optional) Path specifying the location of Windows Sandbox configuration files. If specified, overrides the provider-level Path. |

## Attributes

| Name                 | Type    | Description |
|----------------------|---------|-------------|
| `name`               | String  | Name of the configuration. |
| `path`               | String  | Resolved file path for the configuration. |
| `network`            | Boolean | Indicates if networking is enabled. |
| `virtual_gpu`        | Boolean | Indicates if GPU sharing is enabled. |
| `logon_command`      | String  | The logon command executed in the sandbox. |
| `audio_input`        | Boolean | Indicates if audio input is enabled. |
| `video_input`        | Boolean | Indicates if video input is enabled. |
| `protected_client`   | Boolean | Indicates if Protected Client mode is enabled. |
| `printer_redirection`| Boolean | Indicates if printer redirection is enabled. |
| `clipboard_redirection` | Boolean | Indicates if clipboard redirection is enabled. |
| `memory`             | String  | Amount of memory allocated to the sandbox. |
| `mapped_folders`     |  Set of [mapped_folders](#block-mapped_folders) | List of shared folders with details. |


#### Block: `mapped_folders`

The `mapped_folders` attribute allows you to specify an array of folder mappings between the host and the sandbox.

Each entry in the `mapped_folders` set consists of the following fields:

| Name           | Type    | Description |
|----------------|---------|-------------|
| `host_folder`  | String  | Path to the folder on the host machine. |
| `sandbox_folder` | String | Path where the folder will be mounted inside the sandbox. Optional. |
| `read_only`    | Boolean | Indicates whether the folder is shared in read-only mode. Optional. |

#### Example

```hcl
mapped_folders = [{
  host_folder    = "C:/path/to/host/folder"
  sandbox_folder = "C:/path/to/sandbox/folder"
  read_only      = true
}]
```

This will share the specified host folder with the sandbox, mounting it at the specified sandbox path and making it read-only.


### Data Source Example Usage

```hcl
data "windows_sandbox_configuration" "example" {
  name = "example-sandbox-config"
  path = "/path/to/config"
}
```

This will read and retrieve the configuration details from the specified file for use in your Terraform plan.
