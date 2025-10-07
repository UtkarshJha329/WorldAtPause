package Minigames

import (
	"WorldAtPause/Engine"
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SnekGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	worldBoard Engine.Board[int]

	playerBoardID int
	snake         Snake

	fruitBoardID int

	redTileEntityID    int
	blueTileEntityID   int
	greenTileEntityID  int
	yellowTileEntityID int

	snekHeadTileEntityID int
	snekTailTileEntityID int

	showCurrentUITree bool
	currentUITree     *Engine.UITree

	lastInputDirection Engine.Vector2Int

	totalFruitsConsumed   int
	fruitsToConsumeForWin int
	invalidMoveMade       bool
}

func (snekGameMode *SnekGameMode) Init() {

	snekGameMode.showCurrentUITree = false

	snekGameMode.worldBoard.InitBoard(14, 14, 16, true)

	for y := range snekGameMode.worldBoard.Board_num_rows {
		for x := range snekGameMode.worldBoard.Board_num_cols {
			// snekGameMode.worldBoard.BoardData[y][x] = rand.IntN(4)
			snekGameMode.worldBoard.BoardData[y][x] = 3
		}
	}

	snekGameMode.redTileEntityID = snekGameMode.SceneRef.EntityIDsByName["Snek Red Grid Tile"]
	snekGameMode.blueTileEntityID = snekGameMode.SceneRef.EntityIDsByName["Snek Blue Grid Tile"]
	snekGameMode.greenTileEntityID = snekGameMode.SceneRef.EntityIDsByName["Snek Green Grid Tile"]
	snekGameMode.yellowTileEntityID = snekGameMode.SceneRef.EntityIDsByName["Snek Yellow Grid Tile"]

	snekGameMode.snekHeadTileEntityID = snekGameMode.SceneRef.EntityIDsByName["Snek Head Tile"]
	snekGameMode.snekTailTileEntityID = snekGameMode.SceneRef.EntityIDsByName["Snek Tail Tile"]

	snekGameMode.playerBoardID = -1
	snekGameMode.snake.InitSnakeWithSegments(8, Engine.Vector2Int{X: 7, Y: 7}, Engine.Vector2Int{X: -1, Y: 0})
	for _, snakeSegmentPos := range snekGameMode.snake.segmentPositions {
		snekGameMode.worldBoard.BoardData[snakeSegmentPos.Y][snakeSegmentPos.X] = snekGameMode.playerBoardID
	}

	snekGameMode.fruitBoardID = 0
	snekGameMode.SetFruitToRandomPosOnBoard()

	snekGameMode.totalFruitsConsumed = 0
	snekGameMode.invalidMoveMade = false
}

func (snekGameMode *SnekGameMode) Update() {

	// curSceneRef := snekGameMode.SceneRef
	// entityComponentsRef := curSceneRef.EntityComponentsForScene

	originalSnakeHeadPos := snekGameMode.snake.GetHeadPos()
	modifiedSnakeHeadPos := snekGameMode.snake.GetHeadPos()
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		modifiedSnakeHeadPos.X++
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) {
		modifiedSnakeHeadPos.X--
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		modifiedSnakeHeadPos.Y--
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		modifiedSnakeHeadPos.Y++
	}

	if modifiedSnakeHeadPos.X < 0 {
		modifiedSnakeHeadPos.X = 0
	}
	if modifiedSnakeHeadPos.X >= snekGameMode.worldBoard.Board_num_cols {
		modifiedSnakeHeadPos.X = snekGameMode.worldBoard.Board_num_cols - 1
	}
	if modifiedSnakeHeadPos.Y < 0 {
		modifiedSnakeHeadPos.Y = 0
	}
	if modifiedSnakeHeadPos.Y >= snekGameMode.worldBoard.Board_num_rows {
		modifiedSnakeHeadPos.Y = snekGameMode.worldBoard.Board_num_rows - 1
	}

	snakeMoveDelta := Engine.Subtract_Vector2Int(&modifiedSnakeHeadPos, &originalSnakeHeadPos)

	if snakeMoveDelta.X != 0 || snakeMoveDelta.Y != 0 {
		snekGameMode.lastInputDirection = snakeMoveDelta

		for _, snakeSegmentPos := range snekGameMode.snake.segmentPositions {
			snekGameMode.worldBoard.BoardData[snakeSegmentPos.Y][snakeSegmentPos.X] = 3
		}

		curSnekHeadBoardPos := snekGameMode.snake.GetHeadPos()
		predictedHeadPos := Engine.Add_Vector2Int(&curSnekHeadBoardPos, &snakeMoveDelta)
		consumeFruitAndGrow := snekGameMode.worldBoard.BoardData[predictedHeadPos.Y][predictedHeadPos.X] == snekGameMode.fruitBoardID

		if snekGameMode.snake.MoveHeadOnBoardWithBody(snakeMoveDelta, consumeFruitAndGrow) {
			if consumeFruitAndGrow {
				snekGameMode.SetFruitToRandomPosOnBoard()
				snekGameMode.totalFruitsConsumed++
			}
		} else {
			fmt.Println("Invalid move.")
			snekGameMode.invalidMoveMade = true
		}
		for _, snakeSegmentPos := range snekGameMode.snake.segmentPositions {
			snekGameMode.worldBoard.BoardData[snakeSegmentPos.Y][snakeSegmentPos.X] = snekGameMode.playerBoardID
		}

		// if snekGameMode.snake.SnakeFormsClosedLoop() {
		// 	fmt.Println("Snake forms closed loop.")
		// }
	}

	if !snekGameMode.invalidMoveMade && snekGameMode.totalFruitsConsumed < snekGameMode.fruitsToConsumeForWin {
		snekGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS
	} else if snekGameMode.totalFruitsConsumed >= snekGameMode.fruitsToConsumeForWin {
		snekGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WON
		snekGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	} else if snekGameMode.invalidMoveMade {
		snekGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_LOST
		snekGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	}

}

func (snekGameMode *SnekGameMode) Draw(screenRef *ebiten.Image) {

	// Engine.DrawActiveRoomInScene(snekGameMode.SceneRef, screenRef)
	// Engine.DrawActiveRoomInSceneColliders(snekGameMode.SceneRef, screenRef)

	drawOptions := ebiten.DrawImageOptions{}

	redTileSpriteRef := &snekGameMode.SceneRef.EntityComponentsForScene.Sprites[snekGameMode.redTileEntityID]
	blueTileSpriteRef := &snekGameMode.SceneRef.EntityComponentsForScene.Sprites[snekGameMode.blueTileEntityID]
	greenTileSpriteRef := &snekGameMode.SceneRef.EntityComponentsForScene.Sprites[snekGameMode.greenTileEntityID]
	yellowTileSpriteRef := &snekGameMode.SceneRef.EntityComponentsForScene.Sprites[snekGameMode.yellowTileEntityID]

	snakeBodySpriteRef := &snekGameMode.SceneRef.EntityComponentsForScene.Sprites[snekGameMode.SceneRef.EntityComponentsForScene.PlayerEntityID]

	snekHeadPosOnGrid := snekGameMode.snake.GetHeadPos()
	snekTailSegmentPosOnGrid, snakeDirectionFromPenultimateSegment := snekGameMode.snake.GetTailPosAndDirectionFromBody()

	for y := range snekGameMode.worldBoard.Board_num_rows {
		for x := range snekGameMode.worldBoard.Board_num_cols {
			drawPos := Engine.Vector2{X: float64((x * snekGameMode.worldBoard.TileGridTileSizeInPixels) + snekGameMode.worldBoard.TileGridTotalXOffsetInPixels), Y: float64((y * snekGameMode.worldBoard.TileGridTileSizeInPixels) + snekGameMode.worldBoard.TileGridTotalYOffsetInPixels)}
			switch snekGameMode.worldBoard.BoardData[y][x] {
			case 0:
				redTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			case 1:
				blueTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			case 2:
				greenTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			case 3:
				yellowTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			case snekGameMode.playerBoardID:
				if (x == snekHeadPosOnGrid.X && y == snekHeadPosOnGrid.Y) || (x == snekTailSegmentPosOnGrid.X && y == snekTailSegmentPosOnGrid.Y) {
					continue
				}
				snakeBodySpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
			}
		}
	}

	snekHeadSpriteRef := &snekGameMode.SceneRef.EntityComponentsForScene.Sprites[snekGameMode.snekHeadTileEntityID]
	snekTailSpriteRef := &snekGameMode.SceneRef.EntityComponentsForScene.Sprites[snekGameMode.snekTailTileEntityID]

	snekHeadPos := Engine.Vector2{X: float64((snekHeadPosOnGrid.X * snekGameMode.worldBoard.TileGridTileSizeInPixels) + snekGameMode.worldBoard.TileGridTotalXOffsetInPixels), Y: float64((snekHeadPosOnGrid.Y * snekGameMode.worldBoard.TileGridTileSizeInPixels) + +snekGameMode.worldBoard.TileGridTotalYOffsetInPixels)}

	if snekGameMode.lastInputDirection.X < 0 {
		drawOptions.GeoM.Scale(-1, 1)
		drawOptions.GeoM.Translate(16, 0)
	} else if snekGameMode.lastInputDirection.Y < 0 {
		drawOptions.GeoM.Rotate(3 * math.Pi / 2)
		drawOptions.GeoM.Translate(0, 16)
	} else if snekGameMode.lastInputDirection.Y > 0 {
		drawOptions.GeoM.Rotate(math.Pi / 2)
		drawOptions.GeoM.Translate(16, 0)
	}
	snekHeadSpriteRef.DrawSprite(screenRef, &drawOptions, &snekHeadPos)

	snekTailPos := Engine.Vector2{X: float64((snekTailSegmentPosOnGrid.X * snekGameMode.worldBoard.TileGridTileSizeInPixels) + snekGameMode.worldBoard.TileGridTotalXOffsetInPixels), Y: float64((snekTailSegmentPosOnGrid.Y * snekGameMode.worldBoard.TileGridTileSizeInPixels) + snekGameMode.worldBoard.TileGridTotalYOffsetInPixels)}

	if snakeDirectionFromPenultimateSegment.X < 0 {
		drawOptions.GeoM.Scale(-1, 1)
		drawOptions.GeoM.Translate(16, 0)
	} else if snakeDirectionFromPenultimateSegment.Y < 0 {
		drawOptions.GeoM.Rotate(3 * math.Pi / 2)
		drawOptions.GeoM.Translate(0, 16)
	} else if snakeDirectionFromPenultimateSegment.Y > 0 {
		drawOptions.GeoM.Rotate(math.Pi / 2)
		drawOptions.GeoM.Translate(16, 0)
	}

	snekTailSpriteRef.DrawSprite(screenRef, &drawOptions, &snekTailPos)

	// topLeft, bottomRight := snekGameMode.snake.GetBoundingBox()

	// for y := topLeft.Y; y <= bottomRight.Y; y++ {
	// 	for x := topLeft.X; x <= bottomRight.X; x++ {
	// 		drawPos := Engine.Vector2{X: float64((x * snekGameMode.worldBoard.TileGridTileSizeInPixels) + snekGameMode.worldBoard.TileGridTotalXOffsetInPixels), Y: float64((y * snekGameMode.worldBoard.TileGridTileSizeInPixels) + snekGameMode.worldBoard.TileGridTotalYOffsetInPixels)}
	// 		greenTileSpriteRef.DrawSprite(screenRef, &drawOptions, &drawPos)
	// 	}
	// }

	if snekGameMode.showCurrentUITree {
		snekGameMode.currentUITree.RenderUITree(0, Engine.Vector2{X: 0.0, Y: 0.0}, screenRef, &drawOptions)
	}

}

func (snekGameMode *SnekGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

	if previousGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestType == Engine.QUEST_TYPE_SCORE_LIMIT {
		snekGameMode.fruitsToConsumeForWin = int(previousGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestValues[Engine.QUEST_TYPE_SCORE_LIMIT])
	}
}

func (snekGameMode *SnekGameMode) SetFruitToRandomPosOnBoard() {
	randomPosOnBoard := Engine.Vector2Int{X: rand.IntN(snekGameMode.worldBoard.Board_num_cols), Y: rand.IntN(snekGameMode.worldBoard.Board_num_rows)}
	snekGameMode.worldBoard.BoardData[randomPosOnBoard.Y][randomPosOnBoard.X] = snekGameMode.fruitBoardID
}
