package main

import (
	"encoding/json"
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
