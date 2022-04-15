package engine

import (
	"bytes"
	"io"
	"io/fs"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type audioHandler struct {
	actx  *audio.Context
	files fs.FS
	sr    int
}

// AudioPlayer is a concurrent-safe wrapper around an audio.Player.
type AudioPlayer struct {
	mu sync.Mutex
	p  *audio.Player
}

func (ap *AudioPlayer) Play() {
	ap.mu.Lock()
	ap.p.Pause()
	ap.p.Rewind()
	ap.p.Play()
	ap.mu.Unlock()
}

// Player is a helper method to give an audio.Player from a wav asset.
func (ah *audioHandler) Player(path string) (*AudioPlayer, error) {
	fd, err := asset(ah.files, path)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	mp3, err := mp3.DecodeWithSampleRate(ah.sr, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	p, err := ah.actx.NewPlayer(mp3)
	if err != nil {
		return nil, err
	}
	return &AudioPlayer{p: p}, nil
}
