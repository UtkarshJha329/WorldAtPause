package main

import (
	"encoding/json"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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
	LevelName string          `json:"identifier"`
	Layers    []LayerDataJSON `json:"layerInstances"`
}

type TilemapDataJSON struct {
	Levels []LevelDataJSON `json:"levels"`
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
	levelName string
	layers    []Layer
}

type Tilemap struct {
	tileSet  TileSet
	gridSize int
	levels   []Level
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
	tilemapToFill.levels = make([]Level, len(tilemapDataJSON.Levels))

	tilemapToFill.tileSet.tileSourceImages = make(map[int]*ebiten.Image)

	for currentLevelIndex, level := range tilemapDataJSON.Levels {

		tilemapToFill.levels[currentLevelIndex].levelName = level.LevelName
		tilemapToFill.levels[currentLevelIndex].layers = make([]Layer, len(level.Layers))

		for currentLayerIndex, layer := range level.Layers {

			tilemapToFill.levels[currentLevelIndex].layers[currentLayerIndex].layerName = layer.LayerName
			tilemapToFill.levels[currentLevelIndex].layers[currentLayerIndex].tiles = make([]Tile, len(layer.TilesData))

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

				tilemapToFill.levels[currentLevelIndex].layers[currentLayerIndex].tiles[currentTileIndex].position = tile.TilePos
				tilemapToFill.levels[currentLevelIndex].layers[currentLayerIndex].tiles[currentTileIndex].TileSetTilesID = tile.TileID
			}
		}
	}
	return nil
}
