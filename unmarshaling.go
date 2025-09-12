package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
)

func (v *Vector2) UnmarshalJSON(data []byte) error {
	var arr [2]float64
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	v.x = arr[0]
	v.y = arr[1]
	return nil
}

func (v *Vector2Int) UnmarshalJSON(data []byte) error {
	var arr [2]int
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	v.x = arr[0]
	v.y = arr[1]
	return nil
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

	s.image = LoadImageFromFileSystem(aux.ImageLocation)

	s.sourceStart = aux.SourceStart
	s.sourceEnd = aux.SourceEnd

	return nil
}

func (c *Components) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, val := range raw {
		c.Type = key
		switch key {
		case "Position":
			var pos Vector2
			if err := json.Unmarshal(val, &pos); err != nil {
				return err
			}
			c.Position = &pos
		case "Sprite":
			var spr Sprite
			if err := json.Unmarshal(val, &spr); err != nil {
				return err
			}
			c.Sprite = &spr
		case "Tilemap":
			var tilemapAssetData TilemapAssetData
			if err := json.Unmarshal(val, &tilemapAssetData); err != nil {
				return err
			}
			c.tilemapAssetData = &tilemapAssetData
		}
	}
	return nil
}

func UnmarshalRoot(data *[]byte, root *Root) {
	if err := json.Unmarshal(*data, root); err != nil {
		log.Fatal("Failed to unmarshal provided game data file.")
	}
}

func NewTileMapJson(embeddedFileSystem *embed.FS, filepath string, tilemapJSON *TilemapDataJSON) error {
	contents, err := fs.ReadFile(embeddedFileSystem, filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, tilemapJSON)
	if err != nil {
		return err
	}

	return nil

}
