package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccConfigurationResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wsb_configuration.test", "logon_command", "one"),
				),
			},
			// Update and Read testing
			{
				Config: testAccConfigurationResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wsb_configuration.test", "logon_command", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccConfigurationResourceConfig(logon_command string) string {

	ex, _ := os.Executable()
	directory := filepath.Dir(ex)

	return fmt.Sprintf(`
	resource "wsb_configuration" "test" {
		name                  = "example-configuration"
		path                  = %[1]q
		audio_input           = "true"
		clipboard_redirection   = "true"
		networking   = "true"
		printer_redirection   = "true"
		protected_client      = "true"
		video_input           = "true"
		virtual_gpu           = "true"
		logon_command = %[2]q
}
`, directory, logon_command)
}
