package Engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type TilemapAssetData struct {
	TilemapTilesetFilePath   string `json:"tilemapTilesetsSource"`
	TilemapWorldJSONFilePath string `json:"tilemapSourceJSON"`
	TilemapGridSize          int    `json:"GridSize"`
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
	LevelPxWidth     float64         `json:"pxWid"`
	LevelPxHeight    float64         `json:"pxHei"`
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
	roomDimsInTiles     Vector2Int
}

type Tilemap struct {
	TileSet            TileSet
	GridSize           int
	WorldGridSize      Vector2
	TilemapRooms       map[Vector2Int]*TilemapRoom
	TileCollisionShape BoxCollider
}

func (t *Tilemap) GetRoomDimsInPixels(roomIndex *Vector2Int) Vector2 {
	return Vector2{float64(t.TilemapRooms[*roomIndex].roomDimsInTiles.X * t.GridSize), float64(t.TilemapRooms[*roomIndex].roomDimsInTiles.X * t.GridSize)}
}

func (t *Tilemap) GetFlattenedTileIndexForRoom(roomIndex *Vector2Int, tileXInRoom int, tileYInRoom int) int {

	functionFlattenedTile := tileXInRoom + t.TilemapRooms[*roomIndex].roomDimsInTiles.X*tileYInRoom
	return functionFlattenedTile
}

func (t *Tilemap) GetRoomPos(roomIndex *Vector2Int) Vector2 {
	return t.TilemapRooms[*roomIndex].roomPositionInWorld
}

func (t *Tilemap) GetTilePosInWorld(roomIndex *Vector2Int, tileIndex *Vector2) Vector2 {
	tilePos := Vector2{tileIndex.X * float64(t.GridSize), tileIndex.Y * float64(t.GridSize)}
	roomPosInWorld := t.GetRoomPos(roomIndex)
	return Add_Vector2(&tilePos, &roomPosInWorld)
}

func (t *Tilemap) GetTileOfPointInRoom(roomIndex *Vector2Int, const_pointPosition *Vector2) Vector2 {

	currentRoomPos := t.GetRoomPos(roomIndex)
	// fmt.Println(currentRoomPos)
	pointRoomPos := Subtract_Vector2(const_pointPosition, &currentRoomPos)
	return Vector2{float64(int(pointRoomPos.X) / t.GridSize), float64(int(pointRoomPos.Y) / t.GridSize)}
}

func (t *Tilemap) GetRoomIndexOfPosition(const_pointPosition *Vector2) Vector2Int {
	return Vector2Int{int(const_pointPosition.X) / int(t.WorldGridSize.X), int(const_pointPosition.Y) / int(t.WorldGridSize.Y)}
}

func (t *Tilemap) PointCollidesWithTilemapCollisionLayerInRoom(roomIndex *Vector2Int, pointPosition *Vector2) bool {

	roomDimsInPixels := t.GetRoomDimsInPixels(roomIndex)
	roomPosMaxInPixels := Add_Vector2(&t.TilemapRooms[*roomIndex].roomPositionInWorld, &roomDimsInPixels)
	if pointPosition.X >= roomPosMaxInPixels.X || pointPosition.Y >= roomPosMaxInPixels.Y {
		return false
	}

	pointCurrentTileInRoom := t.GetTileOfPointInRoom(roomIndex, pointPosition)

	flattenedTileIndexInRoom := t.GetFlattenedTileIndexForRoom(roomIndex, int(pointCurrentTileInRoom.X), int(pointCurrentTileInRoom.Y))

	collisionData := t.TilemapRooms[*roomIndex].layers[TILEMAP_COLLISION_LAYER].collisionData
	tileIsSolid := false
	if flattenedTileIndexInRoom >= 0 && flattenedTileIndexInRoom < len((*collisionData)) {
		tileIsSolid = (*collisionData)[flattenedTileIndexInRoom] == 1
	}

	return tileIsSolid
}

func (t *Tilemap) PointCollidesWithTilemapCollisionLayer(const_pointPosition *Vector2) bool {
	pointLevelIndex := t.GetRoomIndexOfPosition(const_pointPosition)
	return t.PointCollidesWithTilemapCollisionLayerInRoom(&pointLevelIndex, const_pointPosition)
}

func NewTilemap(tileSetTexturePath string, tilemapPath string, GridSize int, tilemapToFill *Tilemap, tilemapDataJSON *TilemapDataJSON) error {

	// Will not work with android directly because of lack of filesystem, use go:embed!
	tilemapToFill.TileSet.tileSetTexture = LoadImageFromFileSystem(tileSetTexturePath)

	tilemapToFill.GridSize = GridSize
	tilemapToFill.TilemapRooms = make(map[Vector2Int]*TilemapRoom)

	tilemapToFill.TileCollisionShape = BoxCollider{
		Collider: Collider{
			ColliderOriginOffset: Vector2{0.0, 0.0},
		},
		size: Vector2{float64(tilemapToFill.GridSize), float64(tilemapToFill.GridSize)},
	}
	tilemapToFill.TileCollisionShape.CreateCollisionPoints()

	tilemapToFill.TileSet.tileSourceImages = make(map[int]*ebiten.Image)

	tilemapToFill.WorldGridSize = Vector2{tilemapDataJSON.WorldGridWidth, tilemapDataJSON.WorldGridHeight}

	for _, level := range tilemapDataJSON.Levels {

		levelIndex := Vector2Int{int(level.LevelPosInWorldX) / int(tilemapDataJSON.WorldGridWidth), int(level.LevelPosInWorldY) / int(tilemapDataJSON.WorldGridHeight)}

		currentRoom := TilemapRoom{
			roomName:            level.LevelName,
			roomPositionInWorld: Vector2{level.LevelPosInWorldX, level.LevelPosInWorldY},
			layers:              make([]Layer, len(level.Layers)),
			roomDimsInTiles:     Vector2Int{int(level.LevelPxWidth) / GridSize, int(level.LevelPxHeight) / GridSize},
		}
		tilemapToFill.TilemapRooms[levelIndex] = &currentRoom

		for currentLayerIndex, layer := range level.Layers {

			currentRoom.layers[currentLayerIndex].layerName = layer.LayerName

			switch layer.LayerType {
			case "IntGrid":
				currentRoom.layers[currentLayerIndex].collisionData = &layer.CollisionData

				// if len(tilemapDataJSON.Levels) == 1 {
				// 	for y := range int(currentRoom.roomDimsInTiles.Y) {
				// 		for x := range int(currentRoom.roomDimsInTiles.X) {
				// 			flattenedTileIndex := x + int(currentRoom.roomDimsInTiles.X)*y
				// 			fmt.Print((*currentRoom.layers[currentLayerIndex].collisionData)[flattenedTileIndex], ", ")
				// 		}
				// 		fmt.Println("")
				// 	}
				// }

			case "Tiles":

				currentRoom.layers[currentLayerIndex].tiles = make([]Tile, len(layer.TilesData))

				for currentTileIndex, tile := range layer.TilesData {

					_, ok := tilemapToFill.TileSet.tileSourceImages[tile.TileID]
					if !ok {

						tilemapToFill.TileSet.tileSourceImages[tile.TileID] = tilemapToFill.TileSet.tileSetTexture.SubImage(
							image.Rect(
								tile.TileTextureSrcStart.X,
								tile.TileTextureSrcStart.Y,
								tile.TileTextureSrcStart.X+layer.GridSize,
								tile.TileTextureSrcStart.Y+layer.GridSize)).(*ebiten.Image)
					}

					currentRoom.layers[currentLayerIndex].tiles[currentTileIndex].position = Add_Vector2(&currentRoom.roomPositionInWorld, &tile.TilePos)
					currentRoom.layers[currentLayerIndex].tiles[currentTileIndex].TileSetTilesID = tile.TileID
				}
			}
		}
	}
	return nil
}
