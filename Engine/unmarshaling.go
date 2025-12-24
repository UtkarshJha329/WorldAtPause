package Engine

import (
	"encoding/json"
	"log"
)

func (v *Vector2) UnmarshalJSON(data []byte) error {
	var arr [2]float64
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	v.X = arr[0]
	v.Y = arr[1]
	return nil
}

func (v *Vector2Int) UnmarshalJSON(data []byte) error {
	var arr [2]int
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	v.X = arr[0]
	v.Y = arr[1]
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
		case "Sprite Data":
			var spriteAssetData SpriteAssetData
			if err := json.Unmarshal(val, &spriteAssetData); err != nil {
				return err
			}
			assetData.SpriteAssetData = &spriteAssetData
		case "Collision Shape Data":
			var collisionShapeAssetData CollisionShapeAssetData
			if err := json.Unmarshal(val, &collisionShapeAssetData); err != nil {
				return err
			}
			assetData.CollisionShapeData = &collisionShapeAssetData
		case "Font":
			var fontAssetData FontAssetData
			if err := json.Unmarshal(val, &fontAssetData); err != nil {
				return err
			}
			assetData.fontAssetData = &fontAssetData
		case "Text":
			var textAssetData TextAssetData
			if err := json.Unmarshal(val, &textAssetData); err != nil {
				return err
			}
			assetData.textAssetData = &textAssetData
		case "Tilemap":
			var tilemapAssetData TilemapAssetData
			if err := json.Unmarshal(val, &tilemapAssetData); err != nil {
				return err
			}
			assetData.tilemapAssetData = &tilemapAssetData
		case "Factory":
			var factoryAssetData FactoryAssetData
			if err := json.Unmarshal(val, &factoryAssetData); err != nil {
				return err
			}
			assetData.factoryAssetData = &factoryAssetData
		}
	}
	return nil
}

func UnmarshalRoot(data *[]byte, root *Root) {
	if err := json.Unmarshal(*data, root); err != nil {
		log.Fatal("Failed to unmarshal provided game data file.")
	}
}
