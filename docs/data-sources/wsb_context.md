# Data Source: `wsb_context`

The wsb_context data source helps to generate commonly used folder paths needed for mapping directories in a Windows Sandbox instance. This data source simplifies the retrieval of user-specific and sandbox-specific paths, which are essential when configuring folder mappings in sandbox environments.

It is particularly useful when used with the wsb_configuration resource to create mapped folders within a Windows Sandbox configuration, ensuring paths are correctly set for both the host and sandbox environments.

## Arguments

| Name                 | Type    | Description |
|----------------------|---------|-------------|
| `users_folder` | String | (Optional) Path to the users directory (e.g., C:\\Users)", |
| `username` | String | (Optional) Username for which the configuration is applicable", |
| `downloads_folder_name` | String | (Optional) Name of the Downloads folder (Default: Downloads)", |
| `user_profile_folder` | String | (Optional) Path to the user's profile directory (e.g., C:\\Users\\john)", |
| `user_downloads_folder` | String | (Optional) Path to the user's Downloads directory (e.g., C:\\Users\\john\\Downloads)" |
| `sandbox_users_folder` | String | (Optional) Path to the sandbox users directory", |
| `sandbox_container_username` | String | (Optional) Username for the sandbox container environment", |
| `sandbox_container_user_profile_folder` | String | (Optional) Path to the Sandbox container user profile directory", |
| `sandbox_container_user_downloads_folder` | String | (Optional) Path to the Sandbox container account Downloads directory |

## Example usage

```hcl

data "wsb_context" "example" {
  username              = "JohnDoe"
  users_folder          = "H:\\Users"
  downloads_folder_name = "Downloads"
}

locals {
  mapped_folders = [
    {
      host_folder    = data.wsb_context.example.user_downloads_folder
      sandbox_folder = data.wsb_context.example.sandbox_container_user_downloads_folder
      read_only      = false
    }
  ]
}

resource "wsb_configuration" "main" {
  name                  = "example-configuration"
  path                  = "./"
  virtual_gpu           = false
  networking            = true
  audio_input           = true
  video_input           = true
  protected_client      = true
  printer_redirection   = true
  clipboard_redirection = true
  memory                = 8192

  dynamic "mapped_folders" {
    for_each = local.mapped_folders
    content {
      host_folder    = mapped_folders.value.host_folder
      sandbox_folder = mapped_folders.value.sandbox_folder
      read_only      = mapped_folders.value.read_only
    }
  }
}
```