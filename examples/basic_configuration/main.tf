terraform {
  required_providers {
    windows-sandbox = {
      source = "attilakapostyak/windows-sandbox"
    }
  }
}

provider "windows-sandbox" {
}

resource "windows-sandbox_configuration" "example" {
  path        = "./"
  name        = "basic_configuration"
  audio_input = true
}



