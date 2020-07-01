package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
)

type intro struct {
	grass, cloud, dude *ebiten.Image
	num                int
	audioPlayer        *audio.Player
}

func (i *intro) Init(e *Engine) {
	i.grass, _, _ = ebitenutil.NewImageFromFile(`sprites\grass.png`, ebiten.FilterDefault)
	i.cloud, _, _ = ebitenutil.NewImageFromFile(`sprites\cloud.png`, ebiten.FilterDefault)
	i.dude, _, _ = ebitenutil.NewImageFromFile(`sprites\alman.png`, ebiten.FilterDefault)
	i.num = int(math.Floor(WIDTH / TILESIZE))
	f, err := os.Open(`audio\pause.wav`)
	if err != nil {
		log.Fatal(err)
	}
	d, err := wav.Decode(e.audioCtx, f)
	if err != nil {
		log.Fatal(err)
	}
	i.audioPlayer, err = audio.NewPlayer(e.audioCtx, d)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *intro) Update(e *Engine) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		i.audioPlayer.Pause()
		i.audioPlayer.Rewind()
		i.audioPlayer.Play()
		e.PushState(new(menu))
	}
	return nil
}

func (i *intro) Draw(e *Engine, screen *ebiten.Image) {
	co := &ebiten.DrawImageOptions{}
	co.GeoM.Translate(WIDTH/2, HEIGHT/4)
	screen.DrawImage(i.cloud, co)
	do := &ebiten.DrawImageOptions{}
	do.GeoM.Translate(WIDTH/2, (HEIGHT/1.5)-TILESIZE*2)
	screen.DrawImage(i.dude, do)
	for j := 0; j < i.num; j++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(j*TILESIZE), HEIGHT/1.5)
		screen.DrawImage(i.grass, op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}

type menu struct {
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

func (m *menu) Init(e *Engine) {
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
		d, err := wav.Decode(e.audioCtx, f)
		if err != nil {
			log.Fatal(err)
		}
		m.audio[k], err = audio.NewPlayer(e.audioCtx, d)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (m *menu) Update(e *Engine) error {
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

func (m *menu) Draw(e *Engine, screen *ebiten.Image) {
	o := &ebiten.DrawImageOptions{}
	offset := 32
	sw, sh := m.sidebg.Size()
	img, _ := ebiten.NewImage(sw, sh, ebiten.FilterDefault)
	img.DrawImage(m.sidebg, o)
	for i, option := range m.options {
		text.Draw(img, option, e.tf, sw/2-TILESIZE, sh/4+(offset*i), color.White)
	}
	bord := &ebiten.DrawImageOptions{}
	bord.GeoM.Translate(float64(sw/2-24), float64(sh/4)-11)
	img.DrawImage(m.sel.border, bord)
	o.GeoM.Translate(WIDTH-m.pos, -TILESIZE)
	screen.DrawImage(img, o)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
}
