package main

import (
	"encoding/json"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TilemapAssetData struct {
	TilemapTilesetFilePath   string `json:"tilemapTilesetsSource"`
	TilemapWorldJSONFilePath string `json:"tilemapSourceJSON"`
	TilemapGridSize          int    `json:"gridSize"`
}

type TileDataJSON struct {
	TilePos             Vector2    `json:"px"`
	TileTextureSrcStart Vector2Int `json:"src"`
	TileID              int        `json:"t"`
}

type LayerDataJSON struct {
	LayerName string         `json:"__identifier"`
	GridSize  int            `json:"__gridSize"`
	TilesData []TileDataJSON `json:"gridTiles"`
}

type LevelDataJSON struct {
	LevelName        string          `json:"identifier"`
	LevelPosInWorldX float64         `json:"worldX"`
	LevelPosInWorldY float64         `json:"worldY"`
	Layers           []LayerDataJSON `json:"layerInstances"`
}

type TilemapDataJSON struct {
	WorldGridWidth  float64         `json:"worldGridWidth"`
	WorldGridHeight float64         `json:"worldGridHeight"`
	Levels          []LevelDataJSON `json:"levels"`
}

func NewTileMapJson(filepath string, tilemapJSON *TilemapDataJSON) error {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(contents, tilemapJSON)
	if err != nil {
		return err
	}

	return nil

}

type Tile struct {
	TileSetTilesID int
	position       Vector2
}

type TileSet struct {
	tileSetTexture   *ebiten.Image
	tileSourceImages map[int]*ebiten.Image
}

type Layer struct {
	layerName string
	tiles     []Tile
}

type Level struct {
	levelName            string
	levelPositionInWorld Vector2
	layers               []Layer
}

type Tilemap struct {
	tileSet       TileSet
	gridSize      int
	worldGridSize Vector2
	levels        map[Vector2Int]*Level
}

func NewTilemap(tileSetTexturePath string, tilemapPath string, gridSize int, tilemapToFill *Tilemap) error {

	// Will not work with android directly because of lack of filesystem, use go:embed!
	var err error
	tilemapToFill.tileSet.tileSetTexture, _, err = ebitenutil.NewImageFromFile(tileSetTexturePath)
	if err != nil {
		return err
	}

	var tilemapDataJSON TilemapDataJSON
	NewTileMapJson(tilemapPath, &tilemapDataJSON)

	tilemapToFill.gridSize = gridSize
	tilemapToFill.levels = make(map[Vector2Int]*Level)

	tilemapToFill.tileSet.tileSourceImages = make(map[int]*ebiten.Image)

	tilemapToFill.worldGridSize = Vector2{tilemapDataJSON.WorldGridWidth, tilemapDataJSON.WorldGridHeight}

	for _, level := range tilemapDataJSON.Levels {

		levelIndex := Vector2Int{int(level.LevelPosInWorldX) / int(tilemapDataJSON.WorldGridWidth), int(level.LevelPosInWorldY) / int(tilemapDataJSON.WorldGridHeight)}

		currentLevel := Level{
			levelName:            level.LevelName,
			levelPositionInWorld: Vector2{level.LevelPosInWorldX, level.LevelPosInWorldY},
			layers:               make([]Layer, len(level.Layers)),
		}
		tilemapToFill.levels[levelIndex] = &currentLevel

		for currentLayerIndex, layer := range level.Layers {

			currentLevel.layers[currentLayerIndex].layerName = layer.LayerName
			currentLevel.layers[currentLayerIndex].tiles = make([]Tile, len(layer.TilesData))

			for currentTileIndex, tile := range layer.TilesData {

				_, ok := tilemapToFill.tileSet.tileSourceImages[tile.TileID]
				if !ok {

					tilemapToFill.tileSet.tileSourceImages[tile.TileID] = tilemapToFill.tileSet.tileSetTexture.SubImage(
						image.Rect(
							tile.TileTextureSrcStart.x,
							tile.TileTextureSrcStart.y,
							tile.TileTextureSrcStart.x+layer.GridSize,
							tile.TileTextureSrcStart.y+layer.GridSize)).(*ebiten.Image)
				}

				currentLevel.layers[currentLayerIndex].tiles[currentTileIndex].position = Add_Vector2(&currentLevel.levelPositionInWorld, &tile.TilePos)
				currentLevel.layers[currentLayerIndex].tiles[currentTileIndex].TileSetTilesID = tile.TileID
			}
		}
	}
	return nil
}
