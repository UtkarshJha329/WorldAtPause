package main

import (
	"encoding/json"
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

func (assetData *AssetData) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, val := range raw {
		assetData.AssetType = key
		switch key {
		case "Position":
			var pos Vector2
			if err := json.Unmarshal(val, &pos); err != nil {
				return err
			}
			assetData.Position = &pos
		case "SpriteData":
			var spr SpriteData
			if err := json.Unmarshal(val, &spr); err != nil {
				return err
			}
			assetData.SpriteData = &spr
		case "Tilemap":
			var tilemapAssetData TilemapAssetData
			if err := json.Unmarshal(val, &tilemapAssetData); err != nil {
				return err
			}
			assetData.tilemapAssetData = &tilemapAssetData
		}
	}
	return nil
}

func UnmarshalRoot(data *[]byte, root *Root) {
	if err := json.Unmarshal(*data, root); err != nil {
		log.Fatal("Failed to unmarshal provided game data file.")
	}
}
