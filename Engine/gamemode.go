package Engine

import "github.com/hajimehoshi/ebiten/v2"

type GameMode interface {
	Update()
	Draw(screenRef *ebiten.Image)
}
