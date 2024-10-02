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
  for_each              = toset(local.team_members)
  username              = each.value
  users_folder          = "C:\\Users"
  downloads_folder_name = "Downloads"
}

# Create a separate wsb_configuration for each team member
resource "wsb_configuration" "team" {
  for_each              = toset(local.team_members)
  name                  = "playground-${lower(each.value)}"
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
    for_each = local.team_mapped_folders[each.value]
    content {
      host_folder    = mapped_folders.value.host_folder
      sandbox_folder = mapped_folders.value.sandbox_folder
      read_only      = mapped_folders.value.read_only
    }
  }
}
