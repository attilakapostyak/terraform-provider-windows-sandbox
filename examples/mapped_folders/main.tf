terraform {
  required_providers {
    wsb = {
      source = "attilakapostyak/windows-sandbox"
    }
  }
}

provider "wsb" {
}

locals {
  mapped_folders = [
    {
      host_folder    = "C:\\Users\\JohnDoe\\Downloads"
      sandbox_folder = "C:\\Users\\WDAGUtilityAccount\\Downloads"
      read_only      = false
    }
  ]
}

resource "wsb_configuration" "example" {
  name                  = "mapped_folders"
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
