# Automatic Registry Installation

To install this provider, copy and paste this code into your Terraform configuration (include a version tag).

```hcl
terraform {
  required_providers {
    wsb = {
      source  = "attilakapostyak/windows-sandbox"
      version = "<version>"
    }
  }
}

provider "wsb" {  
  path = "<output_path>"
}
```


## Initialize Terraform

Initialize Terraform so that it installs the new plugins:

```
$ terraform init
```

You should see the following marking the successful plugin installation:

```shell
[...]
Initializing provider plugins...
- Finding registry.example.com/attilakapostyak/windows-sandbox versions matching ">= 0.1.0"...
- Installing registry.example.com/attilakapostyak/windows-sandbox v0.1.0...
- Installed registry.example.com/attilakapostyak/windows-sandbox v0.1.0 (unauthenticated)

Terraform has been successfully initialized!
[...]
```

Now that the plugin is installed, you can simply create a new terraform directory and do usual terraforming.