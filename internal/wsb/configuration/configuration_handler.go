package configuration

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

type WSBConfigurationHandler struct {
	DefaultPath string
}

func (c *WSBConfigurationHandler) BuildConfigurationFileName(configurationPath string, configurationName string) string {
	return path.Join(configurationPath, fmt.Sprintf("%s.%s", configurationName, "wsb"))
}

var DefaultWSBConfigurationHandler = &WSBConfigurationHandler{}

func (c *WSBConfigurationHandler) LoadConfiguration(filePath string) (*WSBConfiguration, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	var config WSBConfiguration
	err = xml.Unmarshal(byteValue, &config)

	return &config, err
}

func (c *WSBConfigurationHandler) SaveConfiguration(filePath string, configuration *WSBConfiguration) error {
	xmlFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer xmlFile.Close()

	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "\t")
	err = encoder.Encode(configuration)
	if err != nil {
		return err
	}

	return nil
}

func (c *WSBConfigurationHandler) DeleteConfiguration(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

func ResolvePath(providerPath, path string) (string, error) {
	switch {
	case len(path) > 0:
		return path, nil
	case len(providerPath) > 0:
		return providerPath, nil
	default:
		return "", errors.New("provider path or input path is required")
	}
}
