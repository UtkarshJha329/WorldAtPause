package Engine

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

type EntityComponents struct {
	TotalNumEntities int
	PlayerEntityID   int
	CameraData       CameraData
	Tilemap          Tilemap
	EntityDead       []bool
	Positions        []Vector2
	Sprites          []Sprite
	CollisionShapes  []CollisionShape
	Texts            []Text
}

func (entityComponents *EntityComponents) IsEntityPlayer(entityID int) bool {
	return entityComponents.PlayerEntityID == entityID
}

func CreateEntityComponents(totalNumEntitiesToCreate int) *EntityComponents {
	return &EntityComponents{
		TotalNumEntities: totalNumEntitiesToCreate,
		EntityDead:       make([]bool, totalNumEntitiesToCreate),
		Positions:        make([]Vector2, totalNumEntitiesToCreate),
		Sprites:          make([]Sprite, totalNumEntitiesToCreate),
		CollisionShapes:  make([]CollisionShape, totalNumEntitiesToCreate),
		Texts:            make([]Text, totalNumEntitiesToCreate)}
}

func (entityComponents *EntityComponents) IsEntityAlive(entityID int) bool {
	return !entityComponents.EntityDead[entityID]
}

func (entityComponents *EntityComponents) IsEntityDead(entityID int) bool {
	return entityComponents.EntityDead[entityID]
}

func PopulateEntityDataFromPrefab(prefabMap map[string]Prefab, curSceneRef *Scene, prefabName string, entityID int) {

	entityComponentsRef := curSceneRef.EntityComponentsForScene

	for _, assetData := range prefabMap[prefabName].AssetDatas {
		switch assetData.AssetType {
		case "Position":
			entityComponentsRef.Positions[entityID] = *assetData.Position
		case "Sprite Data":
			entityComponentsRef.Sprites[entityID] = Sprite{
				Image:           LoadImageFromFileSystem(assetData.SpriteAssetData.SpriteTextureLocation),
				RenderRectStart: assetData.SpriteAssetData.RenderRectStart,
				RenderRectEnd:   assetData.SpriteAssetData.RenderRectEnd,
			}
		case "Collision Shape Data":
			switch assetData.CollisionShapeData.CollisionShapeType {
			case "Box":
				entityComponentsRef.CollisionShapes[entityID] = &BoxCollider{
					size: *assetData.CollisionShapeData.Size,
					Collider: Collider{
						ColliderOriginOffset: *assetData.CollisionShapeData.ColliderOriginOffset,
					},
				}
				entityComponentsRef.CollisionShapes[entityID].CreateCollisionPoints()
			case "Circle":
				entityComponentsRef.CollisionShapes[entityID] = &CircleCollider{
					radius: *assetData.CollisionShapeData.Radius,
					Collider: Collider{
						ColliderOriginOffset: *assetData.CollisionShapeData.ColliderOriginOffset,
					},
				}
				entityComponentsRef.CollisionShapes[entityID].CreateCollisionPoints()
			}
		case "Text":
			font, ok := curSceneRef.Fonts[assetData.textAssetData.FontName]
			if !ok {
				fontPath := prefabMap[assetData.textAssetData.FontName].AssetDatas[0].fontAssetData.FontLocation
				font = &Font{}
				font.LoadFontWithFontFromPath(fontPath)
				font.Face = &text.GoTextFace{
					Source:   font.Font,
					Size:     16,
					Language: language.English,
				}
				curSceneRef.Fonts[assetData.textAssetData.FontName] = font
			}

			curSceneRef.Texts[assetData.textAssetData.TextName] = assetData.textAssetData.TextString
			entityComponentsRef.Texts[entityID] = Text{
				Font:     font,
				TextName: assetData.textAssetData.TextName,
			}

		case "Tilemap":

			var tilemapDataJSON TilemapDataJSON
			tilemapJSONFileContents := LoadFileFromFileSystem(assetData.tilemapAssetData.TilemapWorldJSONFilePath)
			if err := json.Unmarshal(tilemapJSONFileContents, &tilemapDataJSON); err != nil {
				fmt.Println(err, "Failed to loaad tilemapJSON contents.")
			}

			err := NewTilemap(assetData.tilemapAssetData.TilemapTilesetFilePath, assetData.tilemapAssetData.TilemapWorldJSONFilePath, assetData.tilemapAssetData.TilemapGridSize, &entityComponentsRef.Tilemap, &tilemapDataJSON)
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
		default:
			log.Fatal("Unrecognised asset type to override : ", assetData.AssetType)
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

		world.SceneIndexByName[sceneData.SceneName] = index

		world.Scenes[index] = CreateSceneWithNumEntities(totalNumEntitiesInScene)
		world.Scenes[index].SceneName = sceneData.SceneName

		curScene := world.Scenes[index]

		for _, curUISpriteData := range sceneData.UISpriteData {

			curSpriteAssetData := prefabMap[curUISpriteData.PrefabName].AssetDatas[0].SpriteAssetData
			curUISprite := &Sprite{
				Image:           LoadImageFromFileSystem(curSpriteAssetData.SpriteTextureLocation),
				RenderRectStart: curSpriteAssetData.RenderRectStart,
				RenderRectEnd:   curSpriteAssetData.RenderRectEnd,
			}
			curScene.UIRectSprites[curUISpriteData.EntityName] = curUISprite
		}

		runningEntityID := 0
		curSceneEntityComponents := curScene.EntityComponentsForScene

		PopulateEntityDataFromPrefab(prefabMap, curScene, sceneData.PlayerEntityDataForScene.PrefabName, runningEntityID)
		OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &sceneData.PlayerEntityDataForScene)
		curSceneEntityComponents.PlayerEntityID = runningEntityID
		curScene.EntityIDsByName[sceneData.PlayerEntityDataForScene.EntityName] = runningEntityID
		runningEntityID++

		for _, curRoomParsedData := range sceneData.RoomsData {

			roomIndex := curRoomParsedData.RoomIndex
			curScene.RoomsData[roomIndex] = &Room{}
			curRoom := curScene.RoomsData[roomIndex]

			for _, curEntityData := range curRoomParsedData.RoomEntities {

				if prefabMap[curEntityData.PrefabName].AssetType == "UI Sprite Data" {

				} else {

					PopulateEntityDataFromPrefab(prefabMap, curScene, curEntityData.PrefabName, runningEntityID)
					OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &curEntityData)

					switch prefabMap[curEntityData.PrefabName].AssetType {
					case "Enemy":
						curRoom.EnemyEntityIDs = append(curRoom.EnemyEntityIDs, runningEntityID)
					case "Obstacle":
						curRoom.ObstacleEntityIDs = append(curRoom.ObstacleEntityIDs, runningEntityID)
					case "Trigger":
						curRoom.TriggerEntityIDs = append(curRoom.TriggerEntityIDs, runningEntityID)
					case "Item":
						curRoom.ItemEntityIDs = append(curRoom.ItemEntityIDs, runningEntityID)
					}

					curRoom.EntitiesInThisRoom = append(curRoom.EntitiesInThisRoom, runningEntityID)
					curScene.EntityIDsByName[curEntityData.EntityName] = runningEntityID
					runningEntityID += 1
				}
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
