package Minigames

import (
	"WorldAtPause/Engine"

	"github.com/hajimehoshi/ebiten/v2"
)

type MemoryGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	boardNumCols              int
	boardNumRows              int
	boardGridTileSizeInPixels int
	memoryGameBoard           Engine.Board[int]
}

func (memoryGameMode *MemoryGameMode) Init() {
	sceneIndex := memoryGameMode.World.SceneIndexByName[memoryGameMode.SceneRef.SceneName]
	Engine.ReloadSceneWithSceneData(memoryGameMode.SceneRef, &Engine.ScenesData[sceneIndex])

	memoryGameMode.boardNumCols = 10
	memoryGameMode.boardNumRows = 10
	memoryGameMode.boardGridTileSizeInPixels = 16
	memoryGameMode.memoryGameBoard.InitBoard(memoryGameMode.boardNumCols, memoryGameMode.boardNumRows, memoryGameMode.boardGridTileSizeInPixels, true)

	for y := 0; y < memoryGameMode.boardNumRows; y++ {
		for x := 0; x < memoryGameMode.boardNumCols; x++ {
			memoryGameMode.memoryGameBoard.BoardData[y][x] = -1
		}
	}
}

func (memoryGameMode *MemoryGameMode) Update() {

}

func (memoryGameMode *MemoryGameMode) Draw(screenRef *ebiten.Image) {

}

func (memoryGameMode *MemoryGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

}
