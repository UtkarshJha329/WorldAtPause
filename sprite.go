package main

import (
	"encoding/json"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	image       *ebiten.Image
	sourceStart Vector2Int
	sourceEnd   Vector2Int
}

func (s *Sprite) UnmarshalJSON(data []byte) error {
	var aux struct {
		ImageLocation string     `json:"imageSrc"`
		SourceStart   Vector2Int `json:"sourceRectStart"`
		SourceEnd     Vector2Int `json:"sourceRectEnd"`
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	s.image, _, err = ebitenutil.NewImageFromFile(aux.ImageLocation)
	if err != nil {
		log.Fatal("Failed to load sprite from" + aux.ImageLocation)
	}

	s.sourceStart = aux.SourceStart
	s.sourceEnd = aux.SourceEnd

	return nil
}

func (sprite Sprite) DrawSprite(gRef *Game, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions, drawPositionRef *Vector2) {
	drawImgOptionsRef.GeoM.Translate(drawPositionRef.x, drawPositionRef.y)
	screenRef.DrawImage(sprite.image.SubImage(image.Rect(sprite.sourceStart.x, sprite.sourceStart.y, sprite.sourceEnd.x, sprite.sourceEnd.y)).(*ebiten.Image), drawImgOptionsRef)
	drawImgOptionsRef.GeoM.Reset()
}
