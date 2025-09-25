package Minigames

import (
	"WorldAtPause/Engine"
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TILE_RED = iota
	TILE_BLUE
	TILE_GREEN
	TILE_YELLOW
)

type MatchThreeGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	board        Engine.Board[int]
	boardMatches [][]bool

	red_tile_entity_id    int
	blue_tile_entity_id   int
	green_tile_entity_id  int
	yellow_tile_entity_id int

	selected_tile_entity_id int

	currentSelectedTileCoords Engine.Vector2

	currentSelectedCellCoord Engine.Vector2Int
	lastSelectedCellCoord    Engine.Vector2Int

	score             int
	scoreNeededForWin int
}

func (matchThreeGameMode *MatchThreeGameMode) Init() {

	sceneIndex := matchThreeGameMode.World.SceneIndexByName[matchThreeGameMode.SceneRef.SceneName]
	Engine.ReloadSceneWithSceneData(matchThreeGameMode.SceneRef, &Engine.ScenesData[sceneIndex])

	matchThreeGameMode.board.InitBoard(12, 12, 16, true)

	matchThreeGameMode.boardMatches = make([][]bool, matchThreeGameMode.board.Board_num_rows)
	for i := range matchThreeGameMode.board.Board_num_rows {
		matchThreeGameMode.boardMatches[i] = make([]bool, matchThreeGameMode.board.Board_num_cols)
	}

	for y := range matchThreeGameMode.board.Board_num_rows {
		for x := range matchThreeGameMode.board.Board_num_cols {
			matchThreeGameMode.board.BoardData[y][x] = rand.IntN(4)
			matchThreeGameMode.boardMatches[y][x] = false
		}
	}

	matchThreeGameMode.red_tile_entity_id = matchThreeGameMode.SceneRef.EntityIDsByName["Match Three Red Tile"]
	matchThreeGameMode.blue_tile_entity_id = matchThreeGameMode.SceneRef.EntityIDsByName["Match Three Blue Tile"]
	matchThreeGameMode.green_tile_entity_id = matchThreeGameMode.SceneRef.EntityIDsByName["Match Three Green Tile"]
	matchThreeGameMode.yellow_tile_entity_id = matchThreeGameMode.SceneRef.EntityIDsByName["Match Three Yellow Tile"]

	matchThreeGameMode.selected_tile_entity_id = matchThreeGameMode.SceneRef.EntityIDsByName["Match Three Selected Tile"]

	matchThreeGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}
	matchThreeGameMode.lastSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}
	matchThreeGameMode.score = 0

}

func (matchThreeGameMode *MatchThreeGameMode) Update() {

	x, y := ebiten.CursorPosition()
	mousePixelPos := Engine.Vector2{X: float64(x), Y: float64(y)}

	tileGridStartPos := Engine.Vector2{X: float64(matchThreeGameMode.board.TileGridTotalXOffsetInPixels), Y: float64(matchThreeGameMode.board.TileGridTotalYOffsetInPixels)}
	mousePosRelToTileGrid := Engine.Subtract_Vector2(&mousePixelPos, &tileGridStartPos)

	if mousePosRelToTileGrid.X < 0 ||
		mousePosRelToTileGrid.Y < 0 ||
		mousePosRelToTileGrid.X >= float64(matchThreeGameMode.board.TileGridWidthInPixels) ||
		mousePosRelToTileGrid.Y >= float64(matchThreeGameMode.board.TileGridHeightInPixels) {
		matchThreeGameMode.currentSelectedTileCoords = Engine.Vector2{X: -1, Y: -1}
	} else {
		gridCellXExcess := int(mousePosRelToTileGrid.X) % matchThreeGameMode.board.TileGridTileSizeInPixels
		gridCellYExcess := int(mousePosRelToTileGrid.Y) % matchThreeGameMode.board.TileGridTileSizeInPixels

		matchThreeGameMode.currentSelectedTileCoords = Engine.Vector2{X: (mousePosRelToTileGrid.X - float64(gridCellXExcess)), Y: (mousePosRelToTileGrid.Y - float64(gridCellYExcess))}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		matchThreeGameMode.lastSelectedCellCoord = matchThreeGameMode.currentSelectedCellCoord
		if matchThreeGameMode.currentSelectedTileCoords.X < 0 || matchThreeGameMode.currentSelectedTileCoords.Y < 0 {
			matchThreeGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}
		} else {
			matchThreeGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: int(matchThreeGameMode.currentSelectedTileCoords.X) / matchThreeGameMode.board.TileGridTileSizeInPixels, Y: int(matchThreeGameMode.currentSelectedTileCoords.Y) / matchThreeGameMode.board.TileGridTileSizeInPixels}
		}
	}

	if matchThreeGameMode.lastSelectedCellCoord.X >= 0 && matchThreeGameMode.lastSelectedCellCoord.Y >= 0 &&
		matchThreeGameMode.currentSelectedCellCoord.X >= 0 && matchThreeGameMode.currentSelectedCellCoord.Y >= 0 &&
		matchThreeGameMode.lastSelectedCellCoord.X < matchThreeGameMode.board.Board_num_cols && matchThreeGameMode.lastSelectedCellCoord.Y < matchThreeGameMode.board.Board_num_rows &&
		matchThreeGameMode.currentSelectedCellCoord.X < matchThreeGameMode.board.Board_num_cols && matchThreeGameMode.currentSelectedCellCoord.Y < matchThreeGameMode.board.Board_num_rows &&
		!(matchThreeGameMode.lastSelectedCellCoord.X == matchThreeGameMode.currentSelectedCellCoord.X &&
			matchThreeGameMode.lastSelectedCellCoord.Y == matchThreeGameMode.currentSelectedCellCoord.Y) {

		distX := math.Abs(float64(matchThreeGameMode.lastSelectedCellCoord.X - matchThreeGameMode.currentSelectedCellCoord.X))
		distY := math.Abs(float64(matchThreeGameMode.lastSelectedCellCoord.Y - matchThreeGameMode.currentSelectedCellCoord.Y))

		if distX <= 1 && distY <= 1 && !(distX == 1 && distY == 1) {
			SwapValuesInBoard(&matchThreeGameMode.board.BoardData, matchThreeGameMode.lastSelectedCellCoord, matchThreeGameMode.currentSelectedCellCoord, false)

			matchThreeGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}
			matchThreeGameMode.lastSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}

			for CheckBoardForMatches(&matchThreeGameMode.board.BoardData, matchThreeGameMode.board.Board_num_cols, matchThreeGameMode.board.Board_num_rows, 3, &matchThreeGameMode.boardMatches) {
				matchThreeGameMode.score += DeleteMatchesOnBoard(&matchThreeGameMode.board.BoardData, matchThreeGameMode.board.Board_num_cols, matchThreeGameMode.board.Board_num_rows, &matchThreeGameMode.boardMatches)
				FallTilesIntoEmptyCells(&matchThreeGameMode.board.BoardData, matchThreeGameMode.board.Board_num_cols, matchThreeGameMode.board.Board_num_rows)
			}

		} else {

			matchThreeGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}
			matchThreeGameMode.lastSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}
		}
	}

	matchThreeGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS

	if matchThreeGameMode.score >= matchThreeGameMode.scoreNeededForWin {
		fmt.Println("Won with score : ", matchThreeGameMode.score)
		matchThreeGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WON
		matchThreeGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	}

}

func (matchThreeGameMode *MatchThreeGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(matchThreeGameMode.SceneRef, screenRef)

	redSprite := &matchThreeGameMode.SceneRef.EntityComponentsForScene.Sprites[matchThreeGameMode.red_tile_entity_id]
	blueSprite := &matchThreeGameMode.SceneRef.EntityComponentsForScene.Sprites[matchThreeGameMode.blue_tile_entity_id]
	greenSprite := &matchThreeGameMode.SceneRef.EntityComponentsForScene.Sprites[matchThreeGameMode.green_tile_entity_id]
	yellowSprite := &matchThreeGameMode.SceneRef.EntityComponentsForScene.Sprites[matchThreeGameMode.yellow_tile_entity_id]

	drawImgOptions := ebiten.DrawImageOptions{}
	for y := range matchThreeGameMode.board.Board_num_rows {
		for x := range matchThreeGameMode.board.Board_num_cols {
			drawPos := Engine.Vector2{X: float64(matchThreeGameMode.board.TileGridTotalXOffsetInPixels + (x * matchThreeGameMode.board.TileGridTileSizeInPixels)), Y: float64(matchThreeGameMode.board.TileGridTotalYOffsetInPixels + (y * matchThreeGameMode.board.TileGridTileSizeInPixels))}
			switch matchThreeGameMode.board.BoardData[y][x] {
			case 0:
				redSprite.DrawSprite(screenRef, &drawImgOptions, &drawPos)
			case 1:
				blueSprite.DrawSprite(screenRef, &drawImgOptions, &drawPos)
			case 2:
				greenSprite.DrawSprite(screenRef, &drawImgOptions, &drawPos)
			case 3:
				yellowSprite.DrawSprite(screenRef, &drawImgOptions, &drawPos)
			default:
				continue
			}
		}
	}

	if matchThreeGameMode.currentSelectedTileCoords.X >= 0 && matchThreeGameMode.currentSelectedTileCoords.Y >= 0 {
		selectionSprite := &matchThreeGameMode.SceneRef.EntityComponentsForScene.Sprites[matchThreeGameMode.selected_tile_entity_id]
		drawPos := Engine.Vector2{X: float64(matchThreeGameMode.board.TileGridTotalXOffsetInPixels) + matchThreeGameMode.currentSelectedTileCoords.X, Y: float64(matchThreeGameMode.board.TileGridTotalYOffsetInPixels) + matchThreeGameMode.currentSelectedTileCoords.Y}
		selectionSprite.DrawSprite(screenRef, &drawImgOptions, &drawPos)
	}

}

func (matchThreeGameMode *MatchThreeGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

	if previousGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestType == Engine.QUEST_TYPE_SCORE_LIMIT {
		matchThreeGameMode.scoreNeededForWin = int(previousGameStateData.SceneChangeData.QuestIssuedDuringSceneChange.QuestValues[Engine.QUEST_TYPE_SCORE_LIMIT])
	}

}

func SwapValuesInBoard(board *[][]int, coordA Engine.Vector2Int, coordB Engine.Vector2Int, swapWithEmptyTile bool) {

	if ((*board)[coordA.Y][coordA.X] >= 0 && (*board)[coordB.Y][coordB.X] >= 0) || swapWithEmptyTile {
		temp := (*board)[coordA.Y][coordA.X]
		(*board)[coordA.Y][coordA.X] = (*board)[coordB.Y][coordB.X]
		(*board)[coordB.Y][coordB.X] = temp
	}
}

func CheckBoardForMatches(board *[][]int, boardWidth int, boardHeight int, matchRange int, matchedTiles *[][]bool) bool {

	matchedSomething := false
	for y := range boardHeight {
		for x := range boardWidth {

			currentTileID := (*board)[y][x]

			if currentTileID == -1 {
				continue
			}

			matchedHorizontal := true
			for i := range matchRange {
				if x+i >= boardWidth {
					matchedHorizontal = false
					break
				} else if (*board)[y][x+i] != currentTileID {
					matchedHorizontal = false
					break
				}
			}
			matchedVertical := true
			for i := range matchRange {
				if y+i >= boardHeight {
					matchedVertical = false
					break
				} else if (*board)[y+i][x] != currentTileID {
					matchedVertical = false
					break
				}
			}

			if matchedHorizontal {
				for i := range matchRange {
					(*matchedTiles)[y][x+i] = true
					matchedSomething = true
				}
			}
			if matchedVertical {
				for i := range matchRange {
					(*matchedTiles)[y+i][x] = true
					matchedSomething = true
				}
			}

		}
	}

	return matchedSomething
}

func DeleteMatchesOnBoard(board *[][]int, boardWidth int, boardHeight int, matchedTiles *[][]bool) int {

	deleted := 0

	for y := range boardHeight {
		for x := range boardWidth {
			if (*matchedTiles)[y][x] {
				(*board)[y][x] = -1
				(*matchedTiles)[y][x] = false
				deleted++
			}
		}
	}

	return deleted
}

func FallTilesIntoEmptyCells(board *[][]int, boardWidth int, boardHeight int) {

	for reverseY := boardHeight - 2; reverseY >= 0; reverseY-- {
		for y := reverseY; y < boardHeight-1; y++ {
			for x := range boardWidth {
				if (*board)[y+1][x] == -1 {
					curTileCoord := Engine.Vector2Int{X: x, Y: y}
					bottomTileCoord := Engine.Vector2Int{X: x, Y: y + 1}

					SwapValuesInBoard(board, curTileCoord, bottomTileCoord, true)
				}
			}
		}
	}
}
