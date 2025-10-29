package Minigames

import (
	"WorldAtPause/Engine"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TicTacToeGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	boardNumCols              int
	boardNumRows              int
	boardGridTileSizeInPixels int
	ticTacToeBoard            Engine.Board[int]

	currentSelectedTileCoords Engine.Vector2
	currentSelectedCellCoord  Engine.Vector2Int

	X_tile_entity_id        int
	O_tile_entity_id        int
	selected_tile_entity_id int

	curPlayerXOID int
}

func (ticTacToeGameMode *TicTacToeGameMode) Init() {

	sceneIndex := ticTacToeGameMode.World.SceneIndexByName[ticTacToeGameMode.SceneRef.SceneName]
	Engine.ReloadSceneWithSceneData(ticTacToeGameMode.SceneRef, &Engine.ScenesData[sceneIndex])

	ticTacToeGameMode.boardNumCols = 3
	ticTacToeGameMode.boardNumRows = 3
	ticTacToeGameMode.boardGridTileSizeInPixels = 16
	ticTacToeGameMode.ticTacToeBoard.InitBoard(ticTacToeGameMode.boardNumCols, ticTacToeGameMode.boardNumRows, ticTacToeGameMode.boardGridTileSizeInPixels, true)

	for y := 0; y < ticTacToeGameMode.boardNumRows; y++ {
		for x := 0; x < ticTacToeGameMode.boardNumCols; x++ {
			ticTacToeGameMode.ticTacToeBoard.BoardData[y][x] = -1
		}
	}

	ticTacToeGameMode.X_tile_entity_id = ticTacToeGameMode.SceneRef.EntityIDsByName["Tic Tac Toe X Tile"]
	ticTacToeGameMode.O_tile_entity_id = ticTacToeGameMode.SceneRef.EntityIDsByName["Tic Tac Toe O Tile"]
	ticTacToeGameMode.selected_tile_entity_id = ticTacToeGameMode.SceneRef.EntityIDsByName["Tic Tac Toe Selected Tile"]

	ticTacToeGameMode.curPlayerXOID = 1
}

func (ticTacToeGameMode *TicTacToeGameMode) Update() {

	x, y := ebiten.CursorPosition()
	mousePixelPos := Engine.Vector2{X: float64(x), Y: float64(y)}

	tileGridStartPos := Engine.Vector2{X: float64(ticTacToeGameMode.ticTacToeBoard.TileGridTotalXOffsetInPixels), Y: float64(ticTacToeGameMode.ticTacToeBoard.TileGridTotalYOffsetInPixels)}
	mousePosRelToTileGrid := Engine.Subtract_Vector2(&mousePixelPos, &tileGridStartPos)

	if mousePosRelToTileGrid.X < 0 ||
		mousePosRelToTileGrid.Y < 0 ||
		mousePosRelToTileGrid.X >= float64(ticTacToeGameMode.ticTacToeBoard.TileGridWidthInPixels) ||
		mousePosRelToTileGrid.Y >= float64(ticTacToeGameMode.ticTacToeBoard.TileGridHeightInPixels) {
		ticTacToeGameMode.currentSelectedTileCoords = Engine.Vector2{X: -1, Y: -1}
	} else {
		gridCellXExcess := int(mousePosRelToTileGrid.X) % ticTacToeGameMode.ticTacToeBoard.TileGridTileSizeInPixels
		gridCellYExcess := int(mousePosRelToTileGrid.Y) % ticTacToeGameMode.ticTacToeBoard.TileGridTileSizeInPixels

		ticTacToeGameMode.currentSelectedTileCoords = Engine.Vector2{X: (mousePosRelToTileGrid.X - float64(gridCellXExcess)), Y: (mousePosRelToTileGrid.Y - float64(gridCellYExcess))}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if ticTacToeGameMode.currentSelectedTileCoords.X < 0 || ticTacToeGameMode.currentSelectedTileCoords.Y < 0 {
			ticTacToeGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: -1, Y: -1}
		} else {
			ticTacToeGameMode.currentSelectedCellCoord = Engine.Vector2Int{X: int(ticTacToeGameMode.currentSelectedTileCoords.X) / ticTacToeGameMode.ticTacToeBoard.TileGridTileSizeInPixels, Y: int(ticTacToeGameMode.currentSelectedTileCoords.Y) / ticTacToeGameMode.ticTacToeBoard.TileGridTileSizeInPixels}
		}

		if ticTacToeGameMode.currentSelectedCellCoord.X >= 0 && ticTacToeGameMode.currentSelectedCellCoord.X < ticTacToeGameMode.boardNumCols &&
			ticTacToeGameMode.currentSelectedCellCoord.Y >= 0 && ticTacToeGameMode.currentSelectedCellCoord.Y < ticTacToeGameMode.boardNumRows {

			ticTacToeGameMode.ticTacToeBoard.BoardData[ticTacToeGameMode.currentSelectedCellCoord.Y][ticTacToeGameMode.currentSelectedCellCoord.X] = ticTacToeGameMode.curPlayerXOID

			if ticTacToeGameMode.curPlayerXOID == 1 {
				ticTacToeGameMode.curPlayerXOID = 2
			} else if ticTacToeGameMode.curPlayerXOID == 2 {
				ticTacToeGameMode.curPlayerXOID = 1
			}

			winningPlayerID := ticTacToeGameMode.CheckIfEitherPlayerIsWinning()

			if winningPlayerID == 1 {
				ticTacToeGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WON
				ticTacToeGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
			} else if winningPlayerID == 2 {
				ticTacToeGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_LOST
				ticTacToeGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
			}
		}
	} else {
		ticTacToeGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS
	}
}

func (ticTacToeGameMode *TicTacToeGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(ticTacToeGameMode.SceneRef, screenRef)

	XSprite := &ticTacToeGameMode.SceneRef.EntityComponentsForScene.Sprites[ticTacToeGameMode.X_tile_entity_id]
	OSprite := &ticTacToeGameMode.SceneRef.EntityComponentsForScene.Sprites[ticTacToeGameMode.O_tile_entity_id]

	drawImgOptions := ebiten.DrawImageOptions{}
	for y := range ticTacToeGameMode.ticTacToeBoard.Board_num_rows {
		for x := range ticTacToeGameMode.ticTacToeBoard.Board_num_cols {
			drawPos := Engine.Vector2{X: float64(ticTacToeGameMode.ticTacToeBoard.TileGridTotalXOffsetInPixels + (x * ticTacToeGameMode.ticTacToeBoard.TileGridTileSizeInPixels)), Y: float64(ticTacToeGameMode.ticTacToeBoard.TileGridTotalYOffsetInPixels + (y * ticTacToeGameMode.ticTacToeBoard.TileGridTileSizeInPixels))}
			switch ticTacToeGameMode.ticTacToeBoard.BoardData[y][x] {
			case 1:
				XSprite.DrawSprite(screenRef, &drawImgOptions, &drawPos)
			case 2:
				OSprite.DrawSprite(screenRef, &drawImgOptions, &drawPos)
			default:
				continue
			}
		}
	}

	if ticTacToeGameMode.currentSelectedTileCoords.X >= 0 && ticTacToeGameMode.currentSelectedTileCoords.Y >= 0 {
		selectionSprite := &ticTacToeGameMode.SceneRef.EntityComponentsForScene.Sprites[ticTacToeGameMode.selected_tile_entity_id]
		drawPos := Engine.Vector2{X: float64(ticTacToeGameMode.ticTacToeBoard.TileGridTotalXOffsetInPixels) + ticTacToeGameMode.currentSelectedTileCoords.X, Y: float64(ticTacToeGameMode.ticTacToeBoard.TileGridTotalYOffsetInPixels) + ticTacToeGameMode.currentSelectedTileCoords.Y}
		selectionSprite.DrawSprite(screenRef, &drawImgOptions, &drawPos)
	}

}

func (ticTacToeGameMode *TicTacToeGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

}

func (ticTacToeGameMode *TicTacToeGameMode) CheckIfEitherPlayerIsWinning() int {

	winningPlayerID := -1

	// Check rows for win
	for y := 0; y < ticTacToeGameMode.boardNumRows; y++ {

		rowIsWinning := true

		currentPlayerID := ticTacToeGameMode.ticTacToeBoard.BoardData[y][0]
		for x := 1; x < ticTacToeGameMode.boardNumCols; x++ {
			if currentPlayerID != -1 && currentPlayerID != ticTacToeGameMode.ticTacToeBoard.BoardData[y][x] {
				rowIsWinning = false
				break
			}
		}

		if rowIsWinning {
			winningPlayerID = currentPlayerID
			break
		}
	}

	// Check cols for win
	if winningPlayerID == -1 {
		for x := 0; x < ticTacToeGameMode.boardNumCols; x++ {

			colIsWinning := true

			currentPlayerID := ticTacToeGameMode.ticTacToeBoard.BoardData[0][x]
			for y := 1; y < ticTacToeGameMode.boardNumRows; y++ {
				if currentPlayerID != -1 && currentPlayerID != ticTacToeGameMode.ticTacToeBoard.BoardData[y][x] {
					colIsWinning = false
					break
				}
			}

			if colIsWinning {
				winningPlayerID = currentPlayerID
				break
			}
		}
	}

	// Check left diagonal for win
	if winningPlayerID == -1 {
		currentPlayerID := ticTacToeGameMode.ticTacToeBoard.BoardData[0][0]

		if currentPlayerID != -1 && currentPlayerID == ticTacToeGameMode.ticTacToeBoard.BoardData[1][1] &&
			currentPlayerID == ticTacToeGameMode.ticTacToeBoard.BoardData[2][2] {
			winningPlayerID = currentPlayerID
		}
	}

	// Check right diagonal for win
	if winningPlayerID == -1 {
		currentPlayerID := ticTacToeGameMode.ticTacToeBoard.BoardData[0][2]

		if currentPlayerID != -1 && currentPlayerID == ticTacToeGameMode.ticTacToeBoard.BoardData[1][1] &&
			currentPlayerID == ticTacToeGameMode.ticTacToeBoard.BoardData[2][0] {
			winningPlayerID = currentPlayerID
		}
	}

	return winningPlayerID
}
