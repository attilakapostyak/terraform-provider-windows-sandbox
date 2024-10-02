terraform {
  required_providers {
    wsb = {
      source = "attilakapostyak/windows-sandbox"
    }
  }
}

provider "wsb" {
}

resource "wsb_logon_script" "example" {
  name            = "example-logon-command"
  scoop_packages  = ["vcredist2022"]
  winget_packages = ["nvim"]

  pre_installation_scripts  = "Write-Host 'Something to perform before software installation. Add your own Powershell script here'"
  post_installation_scripts = "Write-Host 'Something to perform after software installation. Add your own Powershell script here'"
}

resource "wsb_configuration" "example" {
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
  logon_command         = wsb_logon_script.example.logon_command_script
}
