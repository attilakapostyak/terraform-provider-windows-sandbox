terraform {
  required_providers {
    windows-sandbox = {
      source = "attilakapostyak/windows-sandbox"
    }
  }
}

provider "windows-sandbox" {
}

resource "windows-sandbox_logon_script" "example" {
  name            = "example-logon-command"
  scoop_packages  = ["vcredist2022"]
  winget_packages = ["nvim"]

  pre_installation_scripts  = "Write-Host 'Something to perform before software installation. Add your own Powershell script here'"
  post_installation_scripts = "Write-Host 'Something to perform after software installation. Add your own Powershell script here'"
}

resource "windows-sandbox_configuration" "example" {
  name                  = "logon_command"
  path                  = "./"
  virtual_gpu           = false
  networking            = true
  audio_input           = true
  video_input           = false
  protected_client      = false
  printer_redirection   = true
  clipboard_redirection = true
  memory                = 8192
  logon_command         = windows-sandbox_logon_script.example.logon_command_script
}
