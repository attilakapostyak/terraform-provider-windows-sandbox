terraform {
  required_providers {
    wsb = {
      source = "attilakapostyak/windows-sandbox"
    }
  }
}

provider "wsb" {
  path = "C:\\Users"
}

data "wsb_context" "example" {
  username              = "aaaaaaaaaaa123"
  users_folder          = "Z:\\Uzars"
  downloads_folder_name = "Daunloads"
}

# data "wsb_configuration" "example" {
#   path = "C:\\Users\\attila\\Projects\\terraform-providers\\terraform-provider-windows-sandbox\\test"
#   name = "standard"
# }

locals {
  mapped_folders = [
    {
      host_folder    = data.wsb_context.example.user_downloads_folder
      sandbox_folder = data.wsb_context.example.sandbox_container_user_downloads_folder
      read_only      = false
    }
  ]
}

resource "wsb_logon_script" "example1" {
  name            = "example-configuration"
  winget_packages = ["nvim"]
}

resource "wsb_configuration" "example1" {
  name                  = "example-configuration"
  path                  = "C:\\Users\\attila\\Projects\\terraform-providers\\terraform-provider-windows-sandbox\\test"
  virtual_gpu           = false
  networking            = true
  audio_input           = true
  video_input           = false
  protected_client      = false
  printer_redirection   = true 
  clipboard_redirection = true
  memory                = 8192
  logon_command         = wsb_logon_script.example1.logon_command_script

  dynamic "mapped_folders" {
    for_each = local.mapped_folders
    content {
      host_folder    = mapped_folders.value.host_folder
      sandbox_folder = mapped_folders.value.sandbox_folder
      read_only      = mapped_folders.value.read_only
    }
  }

}

resource "wsb_configuration" "import1" {  
  name = "example-import"
  audio_input = true

}


output "data_wsb_context" {
  value = data.wsb_context.example
}
