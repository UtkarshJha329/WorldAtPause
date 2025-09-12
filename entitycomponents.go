package main

import (
	"embed"
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

func PopulateEntityDataFromPrefab(embeddedFileSystem *embed.FS, prefabMap map[string]Prefab, entityComponents *EntityComponents, prefabName string, entityID int) {

	for _, component := range prefabMap[prefabName].Components {
		switch component.Type {
		case "Position":
			entityComponents.positions[entityID] = *component.Position
		case "Sprite":
			entityComponents.sprites[entityID] = *component.Sprite
		case "Tilemap":
			err := NewTilemap(embeddedFileSystem, component.tilemapAssetData.TilemapTilesetFilePath, component.tilemapAssetData.TilemapWorldJSONFilePath, component.tilemapAssetData.TilemapGridSize, &entityComponents.tilemap)
			if err != nil {
				log.Fatal(err, "\nFailed to create tilemap.")
			}
		}
	}

}

func CreateAndPopulateEntitiesAndComponents(embeddedFileSystem *embed.FS, prefabMap map[string]Prefab, prefabNames []string) *EntityComponents {

	entityComponents := CreateEntityComponents(len(prefabNames))
	for entityID, prefabName := range prefabNames {
		PopulateEntityDataFromPrefab(embeddedFileSystem, prefabMap, entityComponents, prefabName, entityID)

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

func CreateAndPopulateEntitiesAndComponentsFromGameData(embeddedFileSystem *embed.FS, path string) *EntityComponents {
	prefabMap, entityTypesToSpawn, err := LoadGameData(embeddedFileSystem, path)
	if err != nil {
		log.Fatal("Failed to load game data from file.")
	}

	entityComponentData := CreateAndPopulateEntitiesAndComponents(embeddedFileSystem, prefabMap, entityTypesToSpawn)
	// err = NewTilemap("Assets/Maps/TilesetFloor.png", "Assets/Maps/WorldMap.json", 16, &entityComponentData.tilemap)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return entityComponentData
}
