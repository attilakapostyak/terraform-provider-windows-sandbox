package configuration

import (
	"path"
	"runtime"
	"testing"
)

func TestLoadConfiguration_Returns_Valid_Configuration(t *testing.T) {

	_, filename, _, _ := runtime.Caller(0)
	full_path := path.Join(path.Dir(filename), "../../../test_configurations/standard.wsb")
	// Test successful loading of a valid configuration file
	handler := WSBConfigurationHandler{}
	config, err := handler.LoadConfiguration(full_path)

	if err != nil {
		t.Error("Expected LoadConfiguration to return nil error, got:", err)
	}

	if config == nil {
		t.Error("Expected LoadConfiguration to return non-nil configuration, got nil")
	}
}

func TestSaveConfiguration_Save_Valid_Configuration(t *testing.T) {

	_, filename, _, _ := runtime.Caller(0)
	full_path := path.Join(path.Dir(filename), "../../../test_configurations/standard.wsb")

	// Test successful loading of a valid configuration file
	handler := WSBConfigurationHandler{}
	config, err := handler.LoadConfiguration(full_path)

	if err != nil {
		t.Error("Expected LoadConfiguration to return nil error, got:", err)
	}

	if config == nil {
		t.Error("Expected LoadConfiguration to return non-nil configuration, got nil")
	}

	full_path = path.Join(path.Dir(filename), "../../../test_configurations/test_save.wsb")

	err = handler.SaveConfiguration(full_path, config)

	if err != nil {
		t.Error("Expected SaveConfiguration to return nil error, got:", err)
	}
}
