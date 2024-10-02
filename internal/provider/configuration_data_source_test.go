package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigurationDataSource(t *testing.T) {
	ex, _ := os.Executable()
	directory := filepath.Dir(ex)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: NewBasicWSB(directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "network", "true"),
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "virtual_gpu", "true"),
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "logon_command", "explorer.exe C:\\Users\\WDAGUtilityAccount\\Downloads"),
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "audio_input", "true"),
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "video_input", "true"),
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "protected_client", "true"),
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "printer_redirection", "true"),
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "clipboard_redirection", "true"),
					resource.TestCheckResourceAttr("data.wsb_configuration.main", "memory", "4096"),
				),
			},
		},
	})
}

func NewBasicWSB(directory string) string {

	config := `
	<Configuration>
		<Networking>Enable</Networking>
		<vGPU>Enable</vGPU>
		<AudioInput>Enable</AudioInput>  
		<VideoInput>Enable</VideoInput>
		<ProtectedClient>Enable</ProtectedClient>
		<PrinterRedirection>Enable</PrinterRedirection>
		<ClipboardRedirection>Enable</ClipboardRedirection>
		<MemoryInMB>4096</MemoryInMB>  
		<MappedFolders>
		<MappedFolder>
			<HostFolder>C:\Users\Public\Downloads</HostFolder>
			<SandboxFolder>C:\Users\WDAGUtilityAccount\Downloads</SandboxFolder>
			<ReadOnly>true</ReadOnly>
		</MappedFolder>
		</MappedFolders>
		<LogonCommand>
		<Command>explorer.exe C:\Users\WDAGUtilityAccount\Downloads</Command>
		</LogonCommand>
	</Configuration>	
	`

	wsb_filename := filepath.Join(directory, "basic-sandbox.wsb")
	os.WriteFile(wsb_filename, []byte(config), 0644)
	wsb_filename = "basic-sandbox.wsb"

	return fmt.Sprintf("data \"wsb_configuration\" \"main\" {\n  name = \"basic-sandbox\"\n path = \"%s\"\n}", strings.ReplaceAll(directory, "\\", "\\\\"))
}
