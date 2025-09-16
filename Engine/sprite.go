package Engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image           *ebiten.Image
	RenderRectStart Vector2Int
	RenderRectEnd   Vector2Int
}

type SpriteAssetData struct {
	SpriteTextureLocation string     `json:"SpriteTextureLocation"`
	RenderRectStart       Vector2Int `json:"RenderRectStart"`
	RenderRectEnd         Vector2Int `json:"RenderRectEnd"`
}
