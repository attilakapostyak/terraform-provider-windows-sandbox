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
}

# Create a context for each team member
data "wsb_context" "example" {
  for_each              = toset(local.team_members)
  username              = each.value
  users_folder          = "H:\\Users"
  downloads_folder_name = "Downloads"
}

output "all_contexts" {
  value = data.wsb_context.example
}
