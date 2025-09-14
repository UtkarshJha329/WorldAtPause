package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	image           *ebiten.Image
	renderRectStart Vector2Int
	renderRectEnd   Vector2Int
}

type SpriteAssetData struct {
	SpriteTextureLocation string     `json:"SpriteTextureLocation"`
	RenderRectStart       Vector2Int `json:"RenderRectStart"`
	RenderRectEnd         Vector2Int `json:"RenderRectEnd"`
}
