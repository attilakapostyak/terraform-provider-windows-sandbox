resource "wsb_logon_script" "example_configuration" {
  name            = "example-configuration"
  scoop_packages  = ["vcredist2022"]
  winget_packages = ["nvim"]

  pre_installation_scripts  = "Write-Host 'Something to perform before software installation. Add your own Powershell script here'"
  post_installation_scripts = "Write-Host 'Something to perform after software installation. Add your own Powershell script here'"
}

resource "wsb_configuration" "example" {
  name                  = "example-configuration"
  path                  = path.root
  virtual_gpu           = false
  networking            = true
  audio_input           = true
  video_input           = false
  protected_client      = false
  printer_redirection   = true
  clipboard_redirection = true
  memory                = 8192
  logon_command         = wsb_logon_script.example_configuration.logon_command_script

  dynamic "mapped_folders" {
    for_each = local.mapped_folders
    content {
      host_folder    = mapped_folders.value.host_folder
      sandbox_folder = mapped_folders.value.sandbox_folder
      read_only      = mapped_folders.value.read_only
    }
  }
}
