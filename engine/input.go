package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// input describes the different types of inputs.
type input interface {
	ebiten.Key | ebiten.GamepadButton
}

// inputHandler is able to keep track of currently registered keys.
type inputHandler[T input] struct {
	keys             map[T][]func(e *Engine)
	currentStateKeys []T
	keepFlag         bool
	held             bool
}

// handle registers f to be called when k is either pressed once or held down.
func (ih *inputHandler[T]) handle(k T, f func(e *Engine)) {
	if ih.keys[k] == nil {
		ih.keys[k] = make([]func(e *Engine), 0)
	}
	if ih.currentStateKeys == nil {
		ih.currentStateKeys = make([]T, 0)
	}
	ih.currentStateKeys = append(ih.currentStateKeys, k)
	ih.keys[k] = append(ih.keys[k], f)
}

// deregister removes the current State's keys, unless KeepHandlers was called.
func (ih *inputHandler[T]) deregister() {
	if ih.keepFlag {
		ih.currentStateKeys = nil
		ih.keepFlag = false
		return
	}
	if len(ih.currentStateKeys) < 1 {
		return
	}
	for _, k := range ih.currentStateKeys {
		if len(ih.keys[k]) < 1 {
			return
		}
		ih.keys[k] = ih.keys[k][:len(ih.keys[k])-1]
	}
}

// run runs all registered handlers depending on current input.
func (ih *inputHandler[T]) run(e *Engine) {
	if ih.held {
		keys := inpututil.AppendPressedKeys(make([]ebiten.Key, 0))
		if len(keys) < 1 {
			return
		}
		for k, v := range ih.keys {
			for _, key := range keys {
				if key == any(k).(ebiten.Key) {
					for _, f := range v {
						f(e)
					}
				}
			}
		}
		return
	}
	for k, v := range ih.keys {
		// we know that T can only be a type in input
		switch any(k).(type) {
		case ebiten.Key:
			if !inpututil.IsKeyJustPressed(ebiten.Key(k)) {
				continue
			}
		case ebiten.GamepadButton:
			if !inpututil.IsGamepadButtonJustPressed(0, ebiten.GamepadButton(k)) {
				continue
			}
		// if it's not, idk how it got past the compiler.
		// hopefully we can get compile-time guarantees on a type switch like this.
		default:
			panic("not a type in input constraint")
		}
		for _, f := range v {
			f(e)
		}
	}
}

func newInputHandler[T input](held bool) *inputHandler[T] {
	ih := new(inputHandler[T])
	ih.keys = make(map[T][]func(e *Engine))
	ih.held = held
	return ih
}
