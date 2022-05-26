package lib

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func (a *App) SaveSettings(s map[string]interface{}) error {
	p, err := GetSettingsDir()
	if err != nil {
		return err
	}
	p = filepath.Join("settings.yml")
	b, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	err = os.WriteFile(p, b, 0755)
	return err
}

func (a *App) LoadSettings() (map[string]interface{}, error) {
	p, err := GetSettingsDir()
	if err != nil {
		return nil, err
	}
	p = filepath.Join("settings.yml")

	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var in map[string]interface{}
	err = yaml.Unmarshal(b, &in)
	return in, err
}

func GetSettingsDir() (string, error) {
	s, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	s = filepath.Join(s, "treesource")

	return s, err
}
