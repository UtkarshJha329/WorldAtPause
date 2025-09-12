package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed Assets/*
var EmbeddedAssetsFS embed.FS

func LoadImageFromFileSystem(imageLocation string) *ebiten.Image {
	spriteImage, _, err := ebitenutil.NewImageFromFileSystem(EmbeddedAssetsFS, imageLocation)
	if err != nil {
		log.Fatal("Failed to load sprite from" + imageLocation)
	}

	return spriteImage
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

func LoadGameData(embeddedFileSystem *embed.FS, path string) (map[string]Prefab, []string, error) {
	data, err := fs.ReadFile(embeddedFileSystem, path)
	if err != nil {
		log.Fatal("Failed to read provided game data file.")
		return nil, nil, err
	}

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
