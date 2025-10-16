package CrossyRoad

import (
	"WorldAtPause/Engine"

	"github.com/hajimehoshi/ebiten/v2"
)

type Row struct {
	ObjectsAreHarmful     bool
	ObjectsMoveDirectionX int
	InRowObjectsPool      Engine.Pool[Object]

	ObjectsMovementInterpolationOffset   int
	ObjectsMovementInterpolationPoolItem *Engine.PoolItem[Engine.Interpolation[int]]
}

type Chunk struct {
	ChunkIndex int
	Filled     bool
	ChunkBoard Engine.Board[int]

	RowsAndObjects []Row
}

func (chunk *Chunk) IsChunkCurrentlyVisibleOnScreen(worldPositionOffsetY float64) bool {

	chunkPosOffsetY := float64(worldPositionOffsetY) - float64(chunk.ChunkIndex*chunk.ChunkBoard.Board_num_rows*chunk.ChunkBoard.TileGridTileSizeInPixels)

	genericChunkTopPosOnScreen := chunk.ChunkBoard.GetPositionOfTileOnScreen(Engine.Vector2Int{X: 0, Y: 0})
	topChunkPosY := float64(genericChunkTopPosOnScreen.Y) + chunkPosOffsetY

	genericChunkBottomPosOnScreen := chunk.ChunkBoard.GetPositionOfTileOnScreen(Engine.Vector2Int{X: 0, Y: chunk.ChunkBoard.Board_num_rows - 1})
	bottomChunkPosY := float64(genericChunkBottomPosOnScreen.Y) + chunkPosOffsetY

	return topChunkPosY >= 0 && topChunkPosY <= float64(chunk.ChunkBoard.Board_num_rows*chunk.ChunkBoard.TileGridTileSizeInPixels) ||
		bottomChunkPosY >= 0 && bottomChunkPosY <= float64(chunk.ChunkBoard.Board_num_rows*chunk.ChunkBoard.TileGridTileSizeInPixels)
}

func (chunk *Chunk) HasChunkBeenPassed(worldPositionOffsetY float64) bool {

	chunkPosOffsetY := float64(worldPositionOffsetY) - float64(chunk.ChunkIndex*chunk.ChunkBoard.Board_num_rows*chunk.ChunkBoard.TileGridTileSizeInPixels)

	genericChunkTopPosOnScreen := chunk.ChunkBoard.GetPositionOfTileOnScreen(Engine.Vector2Int{X: 0, Y: 0})
	topChunkPosY := float64(genericChunkTopPosOnScreen.Y) + chunkPosOffsetY

	return topChunkPosY > float64(chunk.ChunkBoard.Board_num_rows*chunk.ChunkBoard.TileGridTileSizeInPixels)
}

func (chunk *Chunk) EmptyAllRowsAndObjects() {
	for i := 0; i < len(chunk.RowsAndObjects); i++ {
		chunk.RowsAndObjects[i].ObjectsAreHarmful = false
		chunk.RowsAndObjects[i].ObjectsMoveDirectionX = 0
		chunk.RowsAndObjects[i].InRowObjectsPool.KillAllItemsInPool()
	}
}

func (chunk *Chunk) SlideAllRowsInAppropriateDirectionsByOneTileOnBoard() {
	for i := 0; i < len(chunk.RowsAndObjects); i++ {
		// for i := 0; i < 0; i++ {

		// Move all the object head positions by 1 ?? VVV What is this copy of the data even for?
		chunk.RowsAndObjects[i].InRowObjectsPool.PerformOperationOnAlivePoolItems(func(curObjectPoolItem *Engine.PoolItem[Object]) {
			curObjectPoolItem.Item.HeadTilePositionIndex.X += chunk.RowsAndObjects[i].ObjectsMoveDirectionX
			if curObjectPoolItem.Item.HeadTilePositionIndex.X <= -curObjectPoolItem.Item.ObjectSizeInTiles.X {
				// if curObjectPoolItem.Item.HeadTilePositionIndex.X < 0 {
				curObjectPoolItem.Item.HeadTilePositionIndex.X = chunk.ChunkBoard.Board_num_cols - curObjectPoolItem.Item.ObjectSizeInTiles.X
			} else if curObjectPoolItem.Item.HeadTilePositionIndex.X > chunk.ChunkBoard.Board_num_cols-1 {
				curObjectPoolItem.Item.HeadTilePositionIndex.X = 0
			}
		})

		// Move everything on grid by one in the correct direction
		loopTileFromIndex := 0
		loopTileToIndex := chunk.ChunkBoard.Board_num_cols - 1
		if chunk.RowsAndObjects[i].ObjectsMoveDirectionX > 0 {
			loopTileFromIndex = chunk.ChunkBoard.Board_num_cols - 1
			loopTileToIndex = 0
		}

		loopTileFromValue := chunk.ChunkBoard.BoardData[i][loopTileFromIndex]

		if loopTileFromIndex < loopTileToIndex {
			for j := loopTileFromIndex; j < loopTileToIndex; j++ {
				chunk.ChunkBoard.BoardData[i][j] = chunk.ChunkBoard.BoardData[i][j+1]
			}
		} else {
			for j := loopTileFromIndex; j > loopTileToIndex; j-- {
				chunk.ChunkBoard.BoardData[i][j] = chunk.ChunkBoard.BoardData[i][j-1]
			}
		}

		chunk.ChunkBoard.BoardData[i][loopTileToIndex] = loopTileFromValue
	}
}

func (chunk *Chunk) DrawChunkObjectsInPosition(screenRef *ebiten.Image, obstacleSpriteRef *Engine.Sprite, drawImageOptionsRef *ebiten.DrawImageOptions, worldPositionOffset *Engine.Vector2Int, worldGridSize *Engine.Vector2Int) {

	// numRowsInBoard := chunk.ChunkBoard.Board_num_rows
	// drawImageOptionsRef.ColorScale.SetA(0.15)

	for i := 0; i < len(chunk.RowsAndObjects); i++ {
		chunk.RowsAndObjects[i].InRowObjectsPool.PerformOperationOnAlivePoolItems(func(curObjectPoolItem *Engine.PoolItem[Object]) {

			x := curObjectPoolItem.Item.HeadTilePositionIndex.X
			y := curObjectPoolItem.Item.HeadTilePositionIndex.Y

			// fmt.Println(curObjectPoolItem.Item.HeadTilePositionIndex)

			// drawPos := Engine.Vector2{X: float64((x * chunk.ChunkBoard.TileGridTileSizeInPixels) + chunk.ChunkBoard.TileGridTotalXOffsetInPixels), Y: float64((y*chunk.ChunkBoard.TileGridTileSizeInPixels)+chunk.ChunkBoard.TileGridTotalYOffsetInPixels) + float64(worldPositionOffset.Y) - float64(chunk.ChunkIndex*worldGridSize.Y*chunk.ChunkBoard.TileGridTileSizeInPixels)}
			// if chunk.ChunkBoard.BoardData[numRowsInBoard-y-1][x] >= 5 {
			// 	obstacleSpriteRef.DrawSprite(screenRef, drawImageOptionsRef, &drawPos)
			// }

			tileY := chunk.ChunkBoard.Board_num_rows - y - 1
			for j := 0; j < curObjectPoolItem.Item.ObjectSizeInTiles.X; j++ {

				tileX := x + j

				if tileX < 0 {
					tileX = chunk.ChunkBoard.Board_num_cols + (tileX)
				} else if tileX > chunk.ChunkBoard.Board_num_cols-1 {
					tileX = tileX - (chunk.ChunkBoard.Board_num_cols)
				}

				drawPos := Engine.Vector2{X: float64((tileX*chunk.ChunkBoard.TileGridTileSizeInPixels)+chunk.ChunkBoard.TileGridTotalXOffsetInPixels) + float64(chunk.RowsAndObjects[i].ObjectsMovementInterpolationOffset*chunk.RowsAndObjects[i].ObjectsMoveDirectionX), Y: float64((tileY*chunk.ChunkBoard.TileGridTileSizeInPixels)+chunk.ChunkBoard.TileGridTotalYOffsetInPixels) + float64(worldPositionOffset.Y) - float64(chunk.ChunkIndex*worldGridSize.Y*chunk.ChunkBoard.TileGridTileSizeInPixels)}
				obstacleSpriteRef.DrawSprite(screenRef, drawImageOptionsRef, &drawPos)
			}

		})
	}

	// drawImageOptionsRef.ColorScale.Reset()

}

func (chunk *Chunk) TileIsThreeFourthWayCoveredByObject(tileIndex Engine.Vector2Int) bool {

	if chunk.ChunkBoard.BoardData[tileIndex.Y][tileIndex.X] >= 5 {
		if chunk.RowsAndObjects[tileIndex.Y].ObjectsMovementInterpolationOffset <= 16-(16/4) {
			return true
		}
	}
	if chunk.UseCircularIndexToAccessBoardData(Engine.Vector2Int{X: tileIndex.X - (chunk.RowsAndObjects[tileIndex.Y].ObjectsMoveDirectionX), Y: tileIndex.Y}) >= 5 {
		if chunk.RowsAndObjects[tileIndex.Y].ObjectsMovementInterpolationOffset >= (16 / 4) {
			return true
		}
	}

	return false
}

func (chunk *Chunk) UseCircularIndexToAccessBoardData(tileIndex Engine.Vector2Int) int {
	x := tileIndex.X
	y := tileIndex.Y

	if y < 0 {
		y = chunk.ChunkBoard.Board_num_rows + y
	}
	if y >= chunk.ChunkBoard.Board_num_rows {
		y = y - chunk.ChunkBoard.Board_num_rows
	}

	if x < 0 {
		x = chunk.ChunkBoard.Board_num_cols + x
	}
	if x >= chunk.ChunkBoard.Board_num_cols {
		x = x - chunk.ChunkBoard.Board_num_cols
	}

	return chunk.ChunkBoard.BoardData[y][x]
}
