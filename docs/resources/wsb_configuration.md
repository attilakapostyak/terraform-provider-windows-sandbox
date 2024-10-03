# Resource `windows-sandbox_configuration`

This resource manages a Windows Sandbox configuration file (.wsb) with customizable options such as networking, GPU, folder mapping, and more.

## Arguments

| Name                  | Type       | Default Value | Description                                                                                      |
|-----------------------|------------|---------------|--------------------------------------------------------------------------------------------------|
| `name`                | `string`   | None          | (Required) Name of the configuration |
| `path`                | `string`   | None          | (Optional) Path specifying the location of Windows Sandbox configuration files. Overrides provider's path.    |
| `networking`          | `bool`     | `false`       | (Optional) Enables or disables networking in the sandbox.                                                    |
| `virtual_gpu`         | `bool`     | `false`       | (Optional) Enables or disables GPU sharing in the sandbox.                                                   |
| `logon_command`       | `string`   | None          | (Optional) Command (PowerShell script) to execute when the sandbox user logs on.                             |
| `audio_input`         | `bool`     | `false`       | (Optional) Enables or disables audio input in the sandbox.                                                   |
| `video_input`         | `bool`     | `false`       | (Optional) Enables or disables video input in the sandbox.                                                   |
| `protected_client`    | `bool`     | `false`       | (Optional) Enables or disables Protected Client mode, adding an AppContainer Isolation layer.                 |
| `printer_redirection` | `bool`     | `false`       | (Optional) Enables or disables printer sharing from the host to the sandbox.                                 |
| `clipboard_redirection`| `bool`    | `false`       | (Optional) Enables or disables clipboard sharing between the host and the sandbox.                           |
| `memory`              | `string`   | None          | (Optional) Specifies the amount of memory (MB) the sandbox can use.                                          |
| `mapped_folders`     |  Set of [mapped_folders](#block-mapped_folders) | None | (Optional) List of shared folders with details. |

#### Block: `mapped_folders`

The `mapped_folders` attribute allows you to specify an array of folder mappings between the host and the sandbox.

#### Structure

Each entry in the `mapped_folders` set consists of the following fields:

| Name           | Type    | Default Value | Description |
|----------------|---------|---------------|-------------|
| `host_folder`  | String  | N/A           | (Required) Path to the folder on the host machine. |
| `sandbox_folder` | String | N/A          | (Required) Path where the folder will be mounted inside the sandbox. Optional. |
| `read_only`    | Boolean | `false`       | (Optional) Indicates whether the folder is shared in read-only mode. Optional. |



### Resource Example Usage

```hcl
# Example configuration for windows-sandbox_configuration
resource "windows-sandbox_configuration" "example" {
  name     = "example_config"
  path     = "C:\\path\\to\\sandbox_config"
  networking = true

  virtual_gpu = true
  memory      = "2048"

  mapped_folders {
    host_folder    = "C:\\Users\\Public\\Downloads"
    sandbox_folder = "C:\\Users\\WDAGUtilityAccount\\Downloads"
    read_only      = false
  }
}
```