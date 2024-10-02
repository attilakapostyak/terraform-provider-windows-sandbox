data "wsb_context" "main" {
  username              = "JohnDoe"     # Optional
  users_folder          = "H:\\Users"   # Optional
  downloads_folder_name = "MyDownloads" # Optional
}

data "wsb_configuration" "basic_sandbox" {
  name = "basic-sandbox" # Required
  path = path.root       # Optional
}
