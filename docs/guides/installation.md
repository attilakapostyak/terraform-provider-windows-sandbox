# Automatic Registry Installation

To install this provider, copy and paste this code into your Terraform configuration (include a version tag).

```hcl
terraform {
  required_providers {
    windows-sandbox = {
      source  = "attilakapostyak/windows-sandbox"
      version = "<version>"
    }
  }
}

provider "windows-sandbox" {  
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
Initializing the backend...

Initializing provider plugins...
- Finding latest version of attilakapostyak/windows-sandbox...
- Installing attilakapostyak/windows-sandbox v0.1.1...
- Installed attilakapostyak/windows-sandbox v0.1.1 (...)

Terraform has been successfully initialized!
[...]
```

Now that the plugin is installed, you can simply create a new terraform directory and do usual terraforming.