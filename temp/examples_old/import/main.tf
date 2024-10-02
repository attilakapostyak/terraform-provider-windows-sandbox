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

resource "wsb_configuration" "example" {
  name = "example-configuration"
  path = "C:\\Users\\attila\\Projects\\terraform-providers\\terraform-provider-windows-sandbox\\test"
}