package engine

import (
	"github.com/BurntSushi/toml"
	"github.com/kipukun/game/engine/errors"
)

const (
	name       = "jayarrpeegee"
	savefile   = "file.save"
	width      = 320
	height     = 240
	samplerate = 44100
)

// Config configures the engine.
type Config struct {
	Name, SaveFile            string
	Width, Height, Samplerate int
}

var defaultConfig = &Config{
	Name:       name,
	SaveFile:   savefile,
	Width:      width,
	Height:     height,
	Samplerate: samplerate,
}

func NewConfigFromFile(file string) (*Config, error) {
	var op errors.Op = "NewConfigFromFile"

	dconf := defaultConfig
	_, err := toml.DecodeFile(file, dconf)
	if err != nil {
		return nil, errors.Error(op, "error decoding TOML", err)
	}
	return dconf, nil
}
