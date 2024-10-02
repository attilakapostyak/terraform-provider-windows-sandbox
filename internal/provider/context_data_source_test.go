package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContextDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccContextDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.wsb_context.main", "user_profile_folder", "C:\\Users\\JohnDoe"),
					resource.TestCheckResourceAttr("data.wsb_context.main", "user_downloads_folder", "C:\\Users\\JohnDoe\\Downloads"),
					resource.TestCheckResourceAttr("data.wsb_context.main", "sandbox_users_folder", "C:\\Users"),
					resource.TestCheckResourceAttr("data.wsb_context.main", "sandbox_container_username", "WDAGUtilityAccount"),
					resource.TestCheckResourceAttr("data.wsb_context.main", "sandbox_container_user_profile_folder", "C:\\Users\\WDAGUtilityAccount"),
					resource.TestCheckResourceAttr("data.wsb_context.main", "sandbox_container_user_downloads_folder", "C:\\Users\\WDAGUtilityAccount\\Downloads"),
				),
			},
		},
	})
}

const testAccContextDataSourceConfig = `
data "wsb_context" "main" {
	username              = "JohnDoe"     
	users_folder          = "C:\\Users"   
	downloads_folder_name = "Downloads" 
  }  
`
