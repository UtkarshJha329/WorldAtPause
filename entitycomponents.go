package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type EntityComponents struct {
	totalNumEntities int
	playerEntityID   int
	cameraData       CameraData
	tilemap          Tilemap
	enemyEntityIDs   []int
	itemEntityIDs    []int
	positions        []Vector2
	sprites          []Sprite
	collisionShapes  []CollisionShape
}

func CreateEntityComponents(totalNumEntitiesToCreate int) *EntityComponents {
	return &EntityComponents{
		totalNumEntities: totalNumEntitiesToCreate,
		sprites:          make([]Sprite, totalNumEntitiesToCreate),
		positions:        make([]Vector2, totalNumEntitiesToCreate),
		collisionShapes:  make([]CollisionShape, totalNumEntitiesToCreate)}
}

func PopulateEntityDataFromPrefab(prefabMap map[string]Prefab, entityComponents *EntityComponents, prefabName string, entityID int) {

	for _, assetData := range prefabMap[prefabName].AssetDatas {
		switch assetData.AssetType {
		case "Position":
			entityComponents.positions[entityID] = *assetData.Position
		case "SpriteData":
			entityComponents.sprites[entityID] = Sprite{
				image:           LoadImageFromFileSystem(assetData.SpriteAssetData.SpriteTextureLocation),
				renderRectStart: assetData.SpriteAssetData.RenderRectStart,
				renderRectEnd:   assetData.SpriteAssetData.RenderRectEnd,
			}
		case "CollisionShapeData":
			switch assetData.CollisionShapeData.CollisionShapeType {
			case "Box":
				entityComponents.collisionShapes[entityID] = &BoxCollider{
					size: *assetData.CollisionShapeData.Size,
				}
			case "Circle":
				entityComponents.collisionShapes[entityID] = &CircleCollider{
					radius: *assetData.CollisionShapeData.Radius,
				}
			}
		case "Tilemap":

			var tilemapDataJSON TilemapDataJSON
			tilemapJSONFileContents := LoadFileFromFileSystem(assetData.tilemapAssetData.TilemapWorldJSONFilePath)
			if err := json.Unmarshal(tilemapJSONFileContents, &tilemapDataJSON); err != nil {
				fmt.Println(err, "Failed to loaad tilemapJSON contents.")
			}

			err := NewTilemap(assetData.tilemapAssetData.TilemapTilesetFilePath, assetData.tilemapAssetData.TilemapWorldJSONFilePath, assetData.tilemapAssetData.TilemapGridSize, &entityComponents.tilemap, &tilemapDataJSON)
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
	prefabMap, entityTypesToSpawn, err := LoadGameAssetData(path)
	if err != nil {
		log.Fatal("Failed to load game data from file.")
	}

	entityComponentData := CreateAndPopulateEntitiesAndComponents(prefabMap, entityTypesToSpawn)
	return entityComponentData
}
