package Engine

type AssetData struct {
	AssetType          string
	Position           *Vector2
	SpriteAssetData    *SpriteAssetData
	CollisionShapeData *CollisionShapeAssetData
	tilemapAssetData   *TilemapAssetData
}

type Prefab struct {
	Name       string      `json:"Name"`
	EntityType string      `json:"Entity Type"`
	AssetDatas []AssetData `json:"Asset Data"`
}

type EntitySpawnData struct {
	PrefabName string      `json:"Prefab Name"`
	AssetData  []AssetData `json:"Asset Data"`
}

type RoomEntitySpawnData struct {
	RoomIndex    Vector2Int        `json:"Room Index"`
	RoomEntities []EntitySpawnData `json:"Room Entities"`
}

type SceneData struct {
	PlayerEntityDataForScene EntitySpawnData       `json:"Player Entity"`
	RoomsData                []RoomEntitySpawnData `json:"Rooms Data"`
}

type Root struct {
	Prefabs    []map[string][]Prefab `json:"Prefabs"`
	ScenesData []SceneData           `json:"Scenes Data"`
}

func LoadGameAssetData(path string) (map[string]Prefab, []SceneData, error) {
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

	return prefabMap, root.ScenesData, nil
}
