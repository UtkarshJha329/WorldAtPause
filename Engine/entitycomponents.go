package Engine

import (
	"encoding/json"
	"fmt"
	"log"
)

type EntityComponents struct {
	TotalNumEntities int
	PlayerEntityID   int
	CameraData       CameraData
	Tilemap          Tilemap
	Positions        []Vector2
	Sprites          []Sprite
	CollisionShapes  []CollisionShape
}

func CreateEntityComponents(totalNumEntitiesToCreate int) *EntityComponents {
	return &EntityComponents{
		TotalNumEntities: totalNumEntitiesToCreate,
		Sprites:          make([]Sprite, totalNumEntitiesToCreate),
		Positions:        make([]Vector2, totalNumEntitiesToCreate),
		CollisionShapes:  make([]CollisionShape, totalNumEntitiesToCreate)}
}

func PopulateEntityDataFromPrefab(prefabMap map[string]Prefab, entityComponents *EntityComponents, prefabName string, entityID int) {

	for _, assetData := range prefabMap[prefabName].AssetDatas {
		switch assetData.AssetType {
		case "Position":
			entityComponents.Positions[entityID] = *assetData.Position
		case "SpriteData":
			entityComponents.Sprites[entityID] = Sprite{
				Image:           LoadImageFromFileSystem(assetData.SpriteAssetData.SpriteTextureLocation),
				RenderRectStart: assetData.SpriteAssetData.RenderRectStart,
				RenderRectEnd:   assetData.SpriteAssetData.RenderRectEnd,
			}
		case "CollisionShapeData":
			switch assetData.CollisionShapeData.CollisionShapeType {
			case "Box":
				entityComponents.CollisionShapes[entityID] = &BoxCollider{
					size: *assetData.CollisionShapeData.Size,
					Collider: Collider{
						ColliderOriginOffset: *assetData.CollisionShapeData.ColliderOriginOffset,
					},
				}
				entityComponents.CollisionShapes[entityID].CreateCollisionPoints()
			case "Circle":
				entityComponents.CollisionShapes[entityID] = &CircleCollider{
					radius: *assetData.CollisionShapeData.Radius,
					Collider: Collider{
						ColliderOriginOffset: *assetData.CollisionShapeData.ColliderOriginOffset,
					},
				}
				entityComponents.CollisionShapes[entityID].CreateCollisionPoints()
			}
		case "Tilemap":

			var tilemapDataJSON TilemapDataJSON
			tilemapJSONFileContents := LoadFileFromFileSystem(assetData.tilemapAssetData.TilemapWorldJSONFilePath)
			if err := json.Unmarshal(tilemapJSONFileContents, &tilemapDataJSON); err != nil {
				fmt.Println(err, "Failed to loaad tilemapJSON contents.")
			}

			err := NewTilemap(assetData.tilemapAssetData.TilemapTilesetFilePath, assetData.tilemapAssetData.TilemapWorldJSONFilePath, assetData.tilemapAssetData.TilemapGridSize, &entityComponents.Tilemap, &tilemapDataJSON)
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
			entityComponents.Positions[entityID] = *assetData.Position

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

		world.Scenes[index] = CreateSceneWithNumEntities(totalNumEntitiesInScene)
		world.Scenes[index].SceneType = sceneData.SceneType

		runningEntityID := 0
		curSceneEntityComponents := world.Scenes[index].EntityComponentsForScene

		PopulateEntityDataFromPrefab(prefabMap, curSceneEntityComponents, sceneData.PlayerEntityDataForScene.PrefabName, runningEntityID)
		OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &sceneData.PlayerEntityDataForScene)
		curSceneEntityComponents.PlayerEntityID = runningEntityID
		runningEntityID++

		for _, curRoomParsedData := range sceneData.RoomsData {

			roomIndex := curRoomParsedData.RoomIndex
			world.Scenes[index].RoomsData[roomIndex] = &Room{}
			curRoom := world.Scenes[index].RoomsData[roomIndex]

			for _, curEntityData := range curRoomParsedData.RoomEntities {

				PopulateEntityDataFromPrefab(prefabMap, curSceneEntityComponents, curEntityData.PrefabName, runningEntityID)
				OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &curEntityData)

				switch prefabMap[curEntityData.PrefabName].EntityType {
				case "Enemy":
					curRoom.EnemyEntityIDs = append(curRoom.EnemyEntityIDs, runningEntityID)
				case "Obstacle":
					curRoom.ObstacleEntityIDs = append(curRoom.ObstacleEntityIDs, runningEntityID)
				case "Item":
					curRoom.ItemEntityIDs = append(curRoom.ItemEntityIDs, runningEntityID)
				}

				curRoom.EntitiesInThisRoom = append(curRoom.EntitiesInThisRoom, runningEntityID)
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
