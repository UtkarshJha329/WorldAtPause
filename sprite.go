package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	image       *ebiten.Image
	sourceStart Vector2Int
	sourceEnd   Vector2Int
}
