package engine

import (
	"github.com/BurntSushi/toml"
)

const (
	name       = "jayarrpeegee"
	width      = 320
	height     = 240
	samplerate = 44100
)

type Config struct {
	Name                      string
	Width, Height, Samplerate int
}

var defaultConfig = &Config{
	Name:       name,
	Width:      width,
	Height:     height,
	Samplerate: samplerate,
}

func NewConfigFromFile(file string) (*Config, error) {
	dconf := defaultConfig
	_, err := toml.DecodeFile(file, dconf)
	if err != nil {
		return nil, err
	}
	return dconf, nil
}
