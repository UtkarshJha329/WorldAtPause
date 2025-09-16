package main

import (
	"image"

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
	LayerName     string         `json:"__identifier"`
	LayerType     string         `json:"__type"`
	GridSize      int            `json:"__gridSize"`
	CollisionData []int          `json:"intGridCsv"`
	TilesData     []TileDataJSON `json:"gridTiles"`
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

type Tile struct {
	TileSetTilesID int
	position       Vector2
}

type TileSet struct {
	tileSetTexture   *ebiten.Image
	tileSourceImages map[int]*ebiten.Image
}

const TILEMAP_COLLISION_LAYER = 0

type Layer struct {
	layerName     string
	tiles         []Tile
	collisionData *[]int
}

type TilemapRoom struct {
	roomName            string
	roomPositionInWorld Vector2
	layers              []Layer
}

type Tilemap struct {
	tileSet       TileSet
	gridSize      int
	worldGridSize Vector2
	tilemapRooms  map[Vector2Int]*TilemapRoom
}

func (t *Tilemap) GetTilesPerLevel() Vector2Int {
	return Vector2Int{int(t.worldGridSize.x) / t.gridSize, int(t.worldGridSize.y) / t.gridSize}
}

func (t *Tilemap) GetFlattenedTileIndex(tileXInLevel int, tileYInLevel int) int {
	return tileXInLevel + t.GetTilesPerLevel().x*tileYInLevel
}

func (t *Tilemap) GetLevelPos(levelIndex *Vector2Int) Vector2 {
	return Vector2{float64(levelIndex.x) * t.worldGridSize.x, float64(levelIndex.y) * t.worldGridSize.y}
}

func (t *Tilemap) GetTileOfPointInLevel(levelIndex *Vector2Int, const_pointPosition *Vector2) Vector2 {

	currentLevelPos := t.GetLevelPos(levelIndex)
	pointLevelPos := Subtract_Vector2(const_pointPosition, &currentLevelPos)
	return Vector2{float64(int(pointLevelPos.x) / t.gridSize), float64(int(pointLevelPos.y) / t.gridSize)}
}

func (t *Tilemap) PointCollidesWithTilemapCollisionLayerInLevel(levelIndex *Vector2Int, pointPosition *Vector2) bool {

	pointCurrentTileInLevel := t.GetTileOfPointInLevel(levelIndex, pointPosition)
	flattenedTileIndexInLevel := t.GetFlattenedTileIndex(int(pointCurrentTileInLevel.x), int(pointCurrentTileInLevel.y))

	return (*t.tilemapRooms[*levelIndex].layers[TILEMAP_COLLISION_LAYER].collisionData)[flattenedTileIndexInLevel] == 1
}

func (t *Tilemap) GetRoomIndexOfPosition(const_pointPosition *Vector2) Vector2Int {
	return Vector2Int{int(const_pointPosition.x) / int(t.worldGridSize.x), int(const_pointPosition.y) / int(t.worldGridSize.y)}
}

func (t *Tilemap) PointCollidesWithTilemapCollisionLayer(const_pointPosition *Vector2) bool {
	pointLevelIndex := t.GetRoomIndexOfPosition(const_pointPosition)
	return t.PointCollidesWithTilemapCollisionLayerInLevel(&pointLevelIndex, const_pointPosition)
}

func NewTilemap(tileSetTexturePath string, tilemapPath string, gridSize int, tilemapToFill *Tilemap, tilemapDataJSON *TilemapDataJSON) error {

	// Will not work with android directly because of lack of filesystem, use go:embed!
	var err error
	tilemapToFill.tileSet.tileSetTexture, _, err = ebitenutil.NewImageFromFile(tileSetTexturePath)
	if err != nil {
		return err
	}

	tilemapToFill.gridSize = gridSize
	tilemapToFill.tilemapRooms = make(map[Vector2Int]*TilemapRoom)

	tilemapToFill.tileSet.tileSourceImages = make(map[int]*ebiten.Image)

	tilemapToFill.worldGridSize = Vector2{tilemapDataJSON.WorldGridWidth, tilemapDataJSON.WorldGridHeight}

	for _, level := range tilemapDataJSON.Levels {

		levelIndex := Vector2Int{int(level.LevelPosInWorldX) / int(tilemapDataJSON.WorldGridWidth), int(level.LevelPosInWorldY) / int(tilemapDataJSON.WorldGridHeight)}

		currentRoom := TilemapRoom{
			roomName:            level.LevelName,
			roomPositionInWorld: Vector2{level.LevelPosInWorldX, level.LevelPosInWorldY},
			layers:              make([]Layer, len(level.Layers)),
		}
		tilemapToFill.tilemapRooms[levelIndex] = &currentRoom

		for currentLayerIndex, layer := range level.Layers {

			currentRoom.layers[currentLayerIndex].layerName = layer.LayerName

			switch layer.LayerType {
			case "IntGrid":
				currentRoom.layers[currentLayerIndex].collisionData = &layer.CollisionData

			case "Tiles":

				currentRoom.layers[currentLayerIndex].tiles = make([]Tile, len(layer.TilesData))

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

					currentRoom.layers[currentLayerIndex].tiles[currentTileIndex].position = Add_Vector2(&currentRoom.roomPositionInWorld, &tile.TilePos)
					currentRoom.layers[currentLayerIndex].tiles[currentTileIndex].TileSetTilesID = tile.TileID
				}
			}
		}
	}
	return nil
}
