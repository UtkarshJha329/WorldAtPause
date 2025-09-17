package Engine

import "github.com/hajimehoshi/ebiten/v2"

type GameMode interface {
	Init()
	Update()
	Draw(screenRef *ebiten.Image)
}
