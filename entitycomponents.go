package main

import (
	"log"
)

type EntityComponents struct {
	totalNumEntities int
	playerEntityID   int
	cameraEntityID   int
	cameraData       CameraData
	tilemap          Tilemap
	enemyEntityIDs   []int
	itemEntityIDs    []int
	sprites          []Sprite
	positions        []Vector2
}

func CreateEntityComponents(totalNumEntitiesToCreate int) *EntityComponents {
	return &EntityComponents{
		totalNumEntities: totalNumEntitiesToCreate,
		sprites:          make([]Sprite, totalNumEntitiesToCreate),
		positions:        make([]Vector2, totalNumEntitiesToCreate)}
}

func PopulateEntityDataFromPrefab(prefabMap map[string]Prefab, entityComponents *EntityComponents, prefabName string, entityID int) {

	for _, component := range prefabMap[prefabName].Components {
		switch component.Type {
		case "Position":
			entityComponents.positions[entityID] = *component.Position
		case "Sprite":
			entityComponents.sprites[entityID] = *component.Sprite
		case "Tilemap":
			err := NewTilemap(component.tilemapAssetData.TilemapTilesetFilePath, component.tilemapAssetData.TilemapWorldJSONFilePath, component.tilemapAssetData.TilemapGridSize, &entityComponents.tilemap)
			if err != nil {
				log.Fatal(err, "\nFailed to create tilemap.")
			}
		}
	}

}

func CreateAndPopulateEntitiesAndComponents(prefabMap map[string]Prefab, prefabNames []string) *EntityComponents {

	entityComponents := CreateEntityComponents(len(prefabNames))
	for entityID, prefabName := range prefabNames {
		PopulateEntityDataFromPrefab(prefabMap, entityComponents, prefabName, entityID)

		switch prefabMap[prefabName].EntityType {
		case "Player":
			entityComponents.playerEntityID = entityID
		case "Enemy":
			entityComponents.enemyEntityIDs = append(entityComponents.enemyEntityIDs, entityID)
		case "Item":
			entityComponents.itemEntityIDs = append(entityComponents.itemEntityIDs, entityID)
		}
	}

	return entityComponents
}

func CreateAndPopulateEntitiesAndComponentsFromGameData(path string) *EntityComponents {
	prefabMap, entityTypesToSpawn, err := LoadGameData(path)
	if err != nil {
		log.Fatal("Failed to load game data from file.")
	}

	entityComponentData := CreateAndPopulateEntitiesAndComponents(prefabMap, entityTypesToSpawn)
	// err = NewTilemap("Assets/Maps/TilesetFloor.png", "Assets/Maps/WorldMap.json", 16, &entityComponentData.tilemap)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return entityComponentData
}
