package game

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	bgTexture = "resources/images/background.jpeg"
	font      = "resources/fonts/coolvetica.otf"
	fontSize  = 20
)

type scene struct {
	bg      *sdl.Texture
	bird    *bird
	pipes   *pipes
	timer   *timer
	started bool
}

func NewScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, bgTexture)
	if err != nil {
		return nil, fmt.Errorf("could not load backgound: %v", err)
	}

	b, err := newBird(r)
	if err != nil {
		return nil, err
	}

	ps, err := newPipes(r)
	if err != nil {
		return nil, err
	}

	t := newTimer()

	return &scene{bg: bg, bird: b, pipes: ps, timer: t}, nil
}

func (s *scene) Run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if quit := s.handleEvent(e); quit {
					return
				}
			case <-tick:
				if !s.hasStarted() {
					s.drawTitle(r, "PRESS TO START")
					continue
				}
				if s.bird.isDead() {
					s.drawTitle(r, fmt.Sprintf("SCORE %d", s.timer.seconds()))
					continue
				}
				s.update()
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent, *sdl.KeyboardEvent:
		if s.bird.isDead() || !s.hasStarted() {
			s.restart()
		}
		s.bird.jump()
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.AudioDeviceEvent, *sdl.TextInputEvent:
	default:
		log.Printf("Unknown event %T", event)
	}
	return false
}

func (s *scene) drawTitle(r *sdl.Renderer, text string) error {
	r.Clear()

	f, err := ttf.OpenFont(font, fontSize)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()

	surface, err := f.RenderUTF8Solid(text, sdl.Color{R: 255, G: 100, B: 0, A: 255})
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer surface.Free()

	t, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	rect := &sdl.Rect{X: 200, Y: 250, W: 400, H: 100}
	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present()

	return nil
}

func (s *scene) hasStarted() bool {
	return s.started
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.timer.update()
	s.pipes.touch(s.bird)
}

func (s *scene) restart() {
	s.started = true
	s.bird.restart()
	s.pipes.restart()
	s.timer.restart()
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	if err := s.bird.paint(r); err != nil {
		return err
	}
	if err := s.pipes.paint(r); err != nil {
		return err
	}
	if err := s.timer.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

func (s *scene) Destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}
