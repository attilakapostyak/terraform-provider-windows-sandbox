terraform {
  required_providers {
    wsb = {
      source = "attilakapostyak/windows-sandbox"
    }
  }
}

provider "wsb" {
}

resource "wsb_configuration" "example" {
  path        = "./"
  name        = "basic_configuration"
  audio_input = true
}



