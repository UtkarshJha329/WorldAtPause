package main

type AssetData struct {
	AssetType        string
	Position         *Vector2
	SpriteData       *SpriteData
	tilemapAssetData *TilemapAssetData
}

type Prefab struct {
	Name       string      `json:"Name"`
	EntityType string      `json:"Entity Type"`
	AssetDatas []AssetData `json:"AssetData"`
}

type Root struct {
	Prefabs  []map[string][]Prefab `json:"Prefabs"`
	Entities []string              `json:"Entities"`
}

func LoadGameAssetData(path string) (map[string]Prefab, []string, error) {
	data := LoadFileFromFileSystem(path)
	var root Root
	UnmarshalRoot(&data, &root)
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
