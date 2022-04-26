package engine

import (
	"bytes"
	"io"
	"io/fs"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/kipukun/game/engine/errors"
)

type audioHandler struct {
	actx  *audio.Context
	ps    []*AudioPlayer
	files fs.FS
	sr    int
	v     float64
}

// AudioPlayer is a concurrent-safe wrapper around an audio.Player.
type AudioPlayer struct {
	mu sync.Mutex
	p  *audio.Player
}

// Play plays the AudioPlayer from the beginning.
func (ap *AudioPlayer) Play() {
	ap.mu.Lock()
	ap.p.Pause()
	ap.p.Rewind()
	ap.p.Play()
	ap.mu.Unlock()
}

// SetVolume sets the volume for all players registered in the handler.
func (ap *audioHandler) SetVolume(v float64) {
	if v > 1.0 || v < 0.0 {
		return
	}
	ap.v = v
	for _, p := range ap.ps {
		p.mu.Lock()
		p.p.SetVolume(v)
		p.mu.Unlock()
	}
}

// Player is a helper method to give an audio.Player from an MP3 asset.
func (ah *audioHandler) Player(path string) (*AudioPlayer, error) {
	var op errors.Op = "Player"

	fd, err := asset(ah.files, path)
	if err != nil {
		return nil, errors.Error(op, "error getting asset", err)
	}
	bs, err := io.ReadAll(fd)
	if err != nil {
		return nil, errors.Error(op, "error reading all", err)
	}
	mp3, err := mp3.DecodeWithSampleRate(ah.sr, bytes.NewReader(bs))
	if err != nil {
		return nil, errors.Error(op, "error decoding mp3", err)
	}

	p, err := ah.actx.NewPlayer(mp3)
	if err != nil {
		return nil, errors.Error(op, "error creating Player from context", err)
	}

	ap := &AudioPlayer{p: p}
	ap.p.SetVolume(ah.v)
	ah.ps = append(ah.ps, ap)
	return ap, nil
}
