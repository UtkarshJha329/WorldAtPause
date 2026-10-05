package Engine

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

var PrefabMap map[string]Prefab
var UISpritesData []EntitySpawnData
var ScenesData []SceneData

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

func PopulateEntityDataFromPrefab(PrefabMap map[string]Prefab, curSceneRef *Scene, curRoom *Room, prefabName string, entityID int) {

	entityComponentsRef := curSceneRef.EntityComponentsForScene

	for _, assetData := range PrefabMap[prefabName].AssetDatas {
		switch assetData.AssetType {
		case "Position":
			entityComponentsRef.Positions[entityID] = *assetData.Position
		case "Sprite Data":
			entityComponentsRef.Sprites[entityID] = Sprite{
				Image:                 LoadImageFromFileSystem(assetData.SpriteAssetData.SpriteTextureLocation),
				RenderRectStart:       assetData.SpriteAssetData.RenderRectStart,
				RenderRectEnd:         assetData.SpriteAssetData.RenderRectEnd,
				CurrentAnimationIndex: 0,
			}
			for i := 0; i < len(assetData.SpriteAssetData.AnimationsData); i++ {
				curAnimation := Animation{
					currentFrameCounter: 0,
					animationFramesData: assetData.SpriteAssetData.AnimationsData[i].AnimationFramesData,
					AnimationTimer: curSceneRef.AnimationsTimerSystem.SetTimerFromPoolWithDurationLoopAndFunc(time.Duration(float64(assetData.SpriteAssetData.AnimationsData[i].PerFrameTime)*float64(time.Second)), true, func() {

						curSpriteRef := &entityComponentsRef.Sprites[entityID]

						curAnimationFrameIndex := int(curSpriteRef.Animations[curSpriteRef.CurrentAnimationIndex].currentFrameCounter) % len(curSpriteRef.Animations[curSpriteRef.CurrentAnimationIndex].animationFramesData)
						curAnimationFrameData := curSpriteRef.Animations[curSpriteRef.CurrentAnimationIndex].animationFramesData[curAnimationFrameIndex]

						curSpriteRef.RenderRectStart = curAnimationFrameData.StartFramePos
						curSpriteRef.RenderRectEnd = curAnimationFrameData.EndFramePos

						curSpriteRef.Animations[curSpriteRef.CurrentAnimationIndex].currentFrameCounter++
					}),
				}
				entityComponentsRef.Sprites[entityID].Animations = append(entityComponentsRef.Sprites[entityID].Animations, curAnimation)
				entityComponentsRef.Sprites[entityID].Animations[i].AnimationTimer.Item.PauseTimer()
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
				fontPath := PrefabMap[assetData.textAssetData.FontName].AssetDatas[0].fontAssetData.FontLocation
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

		case "Factory":
			var factory Factory
			factoryEntityToCopy := curSceneRef.EntityIDsByName[assetData.factoryAssetData.FactoryTemplateEntityPrefabName]
			entityComponentsRef.EntityDead[factoryEntityToCopy] = true
			factory.InitFactory(entityComponentsRef, curSceneRef, curRoom, assetData.factoryAssetData.FactoryName, assetData.factoryAssetData.FactoryNumEntities, factoryEntityToCopy, FactoryStateFromString[assetData.factoryAssetData.FactoryType])

			curRoom.Factories = append(curRoom.Factories, &factory)
		}
	}

}

func ReloadEntityDataFromPrefab(PrefabMap map[string]Prefab, curSceneRef *Scene, prefabName string, entityID int) {

	entityComponentsRef := curSceneRef.EntityComponentsForScene

	for _, assetData := range PrefabMap[prefabName].AssetDatas {
		switch assetData.AssetType {
		case "Position":
			entityComponentsRef.Positions[entityID] = *assetData.Position
		case "Sprite Data":
			entityComponentsRef.Sprites[entityID].Image = LoadImageFromFileSystem(assetData.SpriteAssetData.SpriteTextureLocation)
			entityComponentsRef.Sprites[entityID].RenderRectStart = assetData.SpriteAssetData.RenderRectStart
			entityComponentsRef.Sprites[entityID].RenderRectEnd = assetData.SpriteAssetData.RenderRectEnd
		case "Collision Shape Data":
			switch assetData.CollisionShapeData.CollisionShapeType {
			case "Box":
				entityComponentsRef.CollisionShapes[entityID].SetCollisionShapeSize(*assetData.CollisionShapeData.Size)
				entityComponentsRef.CollisionShapes[entityID].SetCollisionShapeColliderOriginOffset(*assetData.CollisionShapeData.ColliderOriginOffset)
				entityComponentsRef.CollisionShapes[entityID].CreateCollisionPoints()
			case "Circle":
				entityComponentsRef.CollisionShapes[entityID].SetCollisionShapeSize(Vector2{*assetData.CollisionShapeData.Radius, *assetData.CollisionShapeData.Radius})
				entityComponentsRef.CollisionShapes[entityID].SetCollisionShapeColliderOriginOffset(*assetData.CollisionShapeData.ColliderOriginOffset)
				entityComponentsRef.CollisionShapes[entityID].CreateCollisionPoints()
			}
		case "Text":
			font, ok := curSceneRef.Fonts[assetData.textAssetData.FontName]
			if !ok {
				fontPath := PrefabMap[assetData.textAssetData.FontName].AssetDatas[0].fontAssetData.FontLocation
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
			entityComponentsRef.Texts[entityID].Font = font
			entityComponentsRef.Texts[entityID].TextName = assetData.textAssetData.TextName

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

func ReloadSceneWithSceneData(curScene *Scene, sceneData *SceneData) {

	runningEntityID := 0
	curSceneEntityComponents := curScene.EntityComponentsForScene

	// FILL PLAYER ENTITY DATA
	ReloadEntityDataFromPrefab(PrefabMap, curScene, sceneData.PlayerEntityDataForScene.PrefabName, runningEntityID)
	OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &sceneData.PlayerEntityDataForScene)

	curSceneEntityComponents.PlayerEntityID = runningEntityID
	runningEntityID++

	for _, curRoomParsedData := range sceneData.RoomsData {

		for _, curEntityData := range curRoomParsedData.RoomEntities {

			ReloadEntityDataFromPrefab(PrefabMap, curScene, curEntityData.PrefabName, runningEntityID)
			OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &curEntityData)

			curSceneEntityComponents.EntityDead[runningEntityID] = false

			runningEntityID += 1
		}
	}
}

func FillSceneWithSceneDataFirstTime(curScene *Scene, sceneData *SceneData) {

	runningEntityID := 0
	curSceneEntityComponents := curScene.EntityComponentsForScene

	// FILL PLAYER ENTITY DATA
	PopulateEntityDataFromPrefab(PrefabMap, curScene, nil, sceneData.PlayerEntityDataForScene.PrefabName, runningEntityID)
	OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &sceneData.PlayerEntityDataForScene)
	curSceneEntityComponents.PlayerEntityID = runningEntityID
	curScene.EntityIDsByName[sceneData.PlayerEntityDataForScene.EntityName] = runningEntityID
	runningEntityID++

	for _, curRoomParsedData := range sceneData.RoomsData {

		roomIndex := curRoomParsedData.RoomIndex
		curScene.RoomsData[roomIndex] = &Room{}
		curRoom := curScene.RoomsData[roomIndex]

		for _, curEntityData := range curRoomParsedData.RoomEntities {

			if PrefabMap[curEntityData.PrefabName].AssetType != "Factory" {
				PopulateEntityDataFromPrefab(PrefabMap, curScene, curRoom, curEntityData.PrefabName, runningEntityID)
				OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &curEntityData)

				switch PrefabMap[curEntityData.PrefabName].AssetType {
				case "Enemy":
					curRoom.EnemyEntityIDs = append(curRoom.EnemyEntityIDs, runningEntityID)
				case "Obstacle":
					curRoom.ObstacleEntityIDs = append(curRoom.ObstacleEntityIDs, runningEntityID)
				case "Invisible Trigger":
					curRoom.InvisibleTriggerEntityIDs = append(curRoom.InvisibleTriggerEntityIDs, runningEntityID)
				case "Visible Trigger":
					curRoom.VisibleTriggerEntityIDs = append(curRoom.VisibleTriggerEntityIDs, runningEntityID)
				case "Item":
					curRoom.ItemEntityIDs = append(curRoom.ItemEntityIDs, runningEntityID)
				}

				curRoom.EntitiesInThisRoom = append(curRoom.EntitiesInThisRoom, runningEntityID)
				curScene.EntityIDsByName[curEntityData.EntityName] = runningEntityID
				runningEntityID += 1
			}
		}

		for _, curEntityData := range curRoomParsedData.RoomEntities {

			if PrefabMap[curEntityData.PrefabName].AssetType == "Factory" {

				PopulateEntityDataFromPrefab(PrefabMap, curScene, curRoom, curEntityData.PrefabName, runningEntityID)
				OverridePrefabDataForEntityWithEntityData(runningEntityID, curSceneEntityComponents, &curEntityData)

				curRoom.EntitiesInThisRoom = append(curRoom.EntitiesInThisRoom, runningEntityID)
				curScene.EntityIDsByName[curEntityData.EntityName] = runningEntityID

				runningEntityID += 1

			}
		}
	}
}

func CreateAndPopulateEntitiesAndComponents() *World {

	world := CreateWorldWithNumScenes(len(ScenesData))

	for _, curUISpriteData := range UISpritesData {

		curSpriteAssetData := PrefabMap[curUISpriteData.PrefabName].AssetDatas[0].SpriteAssetData
		curUISprite := &Sprite{
			Image:           LoadImageFromFileSystem(curSpriteAssetData.SpriteTextureLocation),
			RenderRectStart: curSpriteAssetData.RenderRectStart,
			RenderRectEnd:   curSpriteAssetData.RenderRectEnd,
		}
		world.UIRectSprites[curUISpriteData.EntityName] = curUISprite
	}
	world.MakeYesNoUITree()

	for index, sceneData := range ScenesData {

		totalNumEntitiesInScene := 0
		for _, curRoomData := range sceneData.RoomsData {
			totalNumEntitiesInScene += 1 // for player the single entity in each scene.
			totalNumEntitiesInScene += len(curRoomData.RoomEntities)
		}

		world.SceneIndexByName[sceneData.SceneName] = index

		world.Scenes[index] = CreateSceneWithNumEntities(totalNumEntitiesInScene)
		world.Scenes[index].SceneName = sceneData.SceneName

		curScene := world.Scenes[index]
		curScene.AnimationsTimerSystem.InitWithTimers("Scene Animation Timer System", 128)
		FillSceneWithSceneDataFirstTime(curScene, &sceneData)
	}

	return world
}

func CreateAndPopulateWorldScenesAndEntitiesAndComponentsFromGameData(path string) *World {
	var err error
	PrefabMap, UISpritesData, ScenesData, err = LoadGameAssetData(path)
	if err != nil {
		log.Fatal("Failed to load game data from file.")
	}

	entityComponentData := CreateAndPopulateEntitiesAndComponents()
	return entityComponentData
}
