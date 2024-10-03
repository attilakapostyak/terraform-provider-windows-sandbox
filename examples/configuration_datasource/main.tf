terraform {
  required_providers {
    windows-sandbox = {
      source = "attilakapostyak/windows-sandbox"
    }
  }
}

provider "windows-sandbox" {
}

data "windows-sandbox_configuration" "example" {
  path = ".\\"
  name = "configuration_datasource"
}

output "audio_input" {
  value = data.windows-sandbox_configuration.example.audio_input
}



