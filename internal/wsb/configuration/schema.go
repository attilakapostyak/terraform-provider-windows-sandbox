package configuration

import "encoding/xml"

type WSBConfiguration struct {
	TopComment string `xml:",comment"`

	XMLName       xml.Name      `xml:"Configuration"`
	VGPU          string        `xml:"vGPU,omitempty"`
	Networking    string        `xml:"Networking,omitempty"`
	Memory        string        `xml:"MemoryInMB,omitempty"`
	MappedFolders MappedFolders `xml:"MappedFolders,omitempty"`

	AudioInput           string `xml:"AudioInput,omitempty"`
	VideoInput           string `xml:"VideoInput,omitempty"`
	ProtectedClient      string `xml:"ProtectedClient,omitempty"`
	PrinterRedirection   string `xml:"PrinterRedirection,omitempty"`
	ClipboardRedirection string `xml:"ClipboardRedirection,omitempty"`

	LogonCommand LogonCommand `xml:"LogonCommand,omitempty"`
}

type MappedFolders struct {
	XMLName      xml.Name       `xml:"MappedFolders"`
	MappedFolder []MappedFolder `xml:"MappedFolder,omitempty"`
}

type MappedFolder struct {
	HostFolder    string `xml:"HostFolder,omitempty"`
	SandboxFolder string `xml:"SandboxFolder,omitempty"`
	ReadOnly      string `xml:"ReadOnly,omitempty"`
}

type LogonCommand struct {
	XMLName xml.Name `xml:"LogonCommand,omitempty"`
	Command string   `xml:"Command,omitempty"`
}

func NewWSBConfiguration() *WSBConfiguration {
	return &WSBConfiguration{
		TopComment: "// Managed by Terraform using Windows Sandbox provider: attilakapostyak/windows-sandbox ",
	}
}
