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
					Collider: Collider{
						colliderOriginOffset: *assetData.CollisionShapeData.ColliderOriginOffset,
					},
				}
				entityComponents.collisionShapes[entityID].CreateCollisionPoints()
			case "Circle":
				entityComponents.collisionShapes[entityID] = &CircleCollider{
					radius: *assetData.CollisionShapeData.Radius,
					Collider: Collider{
						colliderOriginOffset: *assetData.CollisionShapeData.ColliderOriginOffset,
					},
				}
				entityComponents.collisionShapes[entityID].CreateCollisionPoints()
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

func OverridePrefabDataForEntityWithEntityData(entityID int, entityComponents *EntityComponents, const_entitySpawnData *EntitySpawnData) {

	for _, assetData := range const_entitySpawnData.AssetData {
		switch assetData.AssetType {

		case "Position":
			entityComponents.positions[entityID] = *assetData.Position

		}
	}
}

func CreateAndPopulateEntitiesAndComponents(prefabMap map[string]Prefab, scenesData []SceneData) *World {

	world := CreateWorldWithNumScenes(len(scenesData))
	for index, sceneData := range scenesData {

		totalNumEntitiesInScene := 0
		for _, curRoomData := range sceneData.RoomsData {
			totalNumEntitiesInScene += 1 // for player the single entity in each scene.
			totalNumEntitiesInScene += len(curRoomData.RoomEntities)
		}

		world.scenes[index] = CreateSceneWithNumEntities(totalNumEntitiesInScene)

		runningEntityID := 0
		curSceneEntityComponents := world.scenes[index].entityComponentsForScene

		PopulateEntityDataFromPrefab(prefabMap, curSceneEntityComponents, sceneData.PlayerEntityDataForScene.PrefabName, runningEntityID)
		OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &sceneData.PlayerEntityDataForScene)
		curSceneEntityComponents.playerEntityID = runningEntityID
		runningEntityID++

		for _, curRoomParsedData := range sceneData.RoomsData {

			roomIndex := curRoomParsedData.RoomIndex
			world.scenes[index].roomsData[roomIndex] = &Room{}
			curRoom := world.scenes[index].roomsData[roomIndex]

			for _, curEntityData := range curRoomParsedData.RoomEntities {

				PopulateEntityDataFromPrefab(prefabMap, curSceneEntityComponents, curEntityData.PrefabName, runningEntityID)
				OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &curEntityData)

				switch prefabMap[curEntityData.PrefabName].EntityType {
				case "Enemy":
					curRoom.enemyEntityIDs = append(curRoom.enemyEntityIDs, runningEntityID)
				case "Obstacle":
					curRoom.obstacleEntityIDs = append(curRoom.obstacleEntityIDs, runningEntityID)
				case "Item":
					curRoom.itemEntityIDs = append(curRoom.itemEntityIDs, runningEntityID)
				}

				curRoom.entitiesInThisRoom = append(curRoom.entitiesInThisRoom, runningEntityID)
				runningEntityID += 1
			}
		}

	}

	return world
}

func CreateAndPopulateWorldScenesAndEntitiesAndComponentsFromGameData(path string) *World {
	prefabMap, scenesData, err := LoadGameAssetData(path)
	if err != nil {
		log.Fatal("Failed to load game data from file.")
	}

	entityComponentData := CreateAndPopulateEntitiesAndComponents(prefabMap, scenesData)
	return entityComponentData
}
