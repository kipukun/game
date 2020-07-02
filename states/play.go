package states

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/kipukun/game/engine"
	"github.com/markbates/pkger"
)

type PlayState struct {
	cloud       *ebiten.Image
	audioPlayer *audio.Player
}

func (p *PlayState) Init(e *engine.Engine) {
	f, err := pkger.Open("/assets/sprites/cloud.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	p.cloud, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	fd, err := pkger.Open("/assets/audio/pause.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	d, err := wav.Decode(e.AudioCtx(), fd)
	if err != nil {
		log.Fatal(err)
	}
	p.audioPlayer, err = audio.NewPlayer(e.AudioCtx(), d)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *PlayState) Update(e *engine.Engine) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.audioPlayer.Pause()
		p.audioPlayer.Rewind()
		p.audioPlayer.Play()
		// e.PushState(new(MenuState))
	}
	return nil
}

func (p *PlayState) Draw(e *engine.Engine, screen *ebiten.Image) {
	w, h := e.Size()
	co := &ebiten.DrawImageOptions{}
	co.GeoM.Translate(w/2, h/4)
	screen.DrawImage(p.cloud, co)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}

type MenuState struct {
	sidebg             *ebiten.Image
	audio              map[string]*audio.Player
	pos, vel, friction float64
	options            []string

	sel struct {
		border                *ebiten.Image
		pos, vel, fric, accel float64
		idx                   int
	}
}

func (m *MenuState) Init(e *engine.Engine) {
	m.vel = 20
	m.friction = 0.75
	m.options = []string{"Equip", "Items", "Skills", "Map", "System"}
	m.sidebg, _, _ = ebitenutil.NewImageFromFile(`sprites\menu_bg.png`, ebiten.FilterDefault)
	m.sel.border, _, _ = ebitenutil.NewImageFromFile(`sprites\select_outline.png`, ebiten.FilterDefault)
	m.audio = map[string]*audio.Player{
		"pause_exit": nil,
		"menu_move":  nil,
		"select":     nil}
	for k := range m.audio {
		f, err := os.Open(fmt.Sprintf(`audio\%s.wav`, k))
		if err != nil {
			log.Fatal(err)
		}
		d, err := wav.Decode(e.AudioCtx(), f)
		if err != nil {
			log.Fatal(err)
		}
		m.audio[k], err = audio.NewPlayer(e.AudioCtx(), d)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (m *MenuState) Update(e *engine.Engine) error {
	m.pos += m.vel
	m.vel *= m.friction
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		m.audio["select"].Pause()
		m.audio["select"].Rewind()
		m.audio["select"].Play()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		m.audio["pause_exit"].Pause()
		m.audio["pause_exit"].Rewind()
		m.audio["pause_exit"].Play()
		e.PopState()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.audio["menu_move"].Pause()
		m.audio["menu_move"].Rewind()
		m.audio["menu_move"].Play()
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			if m.sel.idx == 0 {
				m.sel.idx = len(m.options) - 1
				return nil
			}
			m.sel.idx -= 1
		} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			if m.sel.idx == len(m.options)-1 {
				m.sel.idx = 0
				return nil
			}
			m.sel.idx += 1
		}
	}

	return nil
}

func (m *MenuState) Draw(e *engine.Engine, screen *ebiten.Image) {
	w, _ := e.Size()
	o := &ebiten.DrawImageOptions{}
	offset := 32
	sw, sh := m.sidebg.Size()
	img, _ := ebiten.NewImage(sw, sh, ebiten.FilterDefault)
	img.DrawImage(m.sidebg, o)
	for i, option := range m.options {
		text.Draw(img, option, e.Font(), sw/2-engine.TileSize, sh/4+(offset*i), color.White)
	}
	bord := &ebiten.DrawImageOptions{}
	bord.GeoM.Translate(float64(sw/2-24), float64(sh/4)-11)
	img.DrawImage(m.sel.border, bord)
	o.GeoM.Translate(w-m.pos, -engine.TileSize)
	screen.DrawImage(img, o)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}
