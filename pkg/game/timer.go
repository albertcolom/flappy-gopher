package game

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"strconv"
	"sync"
)

type timer struct {
	mu   sync.RWMutex
	time int
}

func newTimer() *timer {
	return &timer{time: 0}
}

func (t *timer) paint(r *sdl.Renderer) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	f, err := ttf.OpenFont(font, fontSize)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()

	s, err := f.RenderUTF8Solid(strconv.Itoa(t.seconds()), sdl.Color{R: 255, G: 100, B: 0, A: 255})
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	texture, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer texture.Destroy()

	rect := &sdl.Rect{X: 10, Y: 0, W: 75, H: 75}
	if err := r.Copy(texture, nil, rect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	return nil
}

func (t *timer) seconds() int {
	return t.time / 60
}

func (t *timer) update() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.time++
}

func (t *timer) restart() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.time = 0
}
