package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	image           *ebiten.Image
	renderRectStart Vector2Int
	renderRectEnd   Vector2Int
}

type SpriteData struct {
	SpriteTextureLocation string     `json:"SpriteTextureLocation"`
	RenderRectStart       Vector2Int `json:"renderRectStart"`
	RenderRectEnd         Vector2Int `json:"renderRectEnd"`
}
