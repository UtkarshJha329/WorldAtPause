package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed Assets/*
var EmbeddedAssetsFS embed.FS

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

	s.image, _, err = ebitenutil.NewImageFromFileSystem(EmbeddedAssetsFS, aux.ImageLocation)
	if err != nil {
		log.Fatal("Failed to load sprite from" + aux.ImageLocation)
	}

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

func LoadGameData(embeddedFileSystem *embed.FS, path string) (map[string]Prefab, []string, error) {
	data, err := fs.ReadFile(embeddedFileSystem, path)
	if err != nil {
		log.Fatal("Failed to read provided game data file.")
		return nil, nil, err
	}

	var root Root
	if err := json.Unmarshal(data, &root); err != nil {
		log.Fatal("Failed to unmarshal provided game data file.")
		return nil, nil, err
	}

	prefabMap := make(map[string]Prefab)

	for _, entry := range root.Prefabs {
		for _, prefabList := range entry {
			for _, p := range prefabList {
				prefabMap[p.Name] = p
			}
		}
	}

	return prefabMap, root.Entities, nil
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
