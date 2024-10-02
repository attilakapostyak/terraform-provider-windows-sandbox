terraform {
  required_providers {
    wsb = {
      source = "attilakapostyak/windows-sandbox"
    }
  }
}

provider "wsb" {
}

data "wsb_configuration" "example" {
  path = ".\\"
  name = "configuration_datasource"
}

output "audio_input" {
  value = data.wsb_configuration.example.audio_input
}



