package Engine

import "github.com/hajimehoshi/ebiten/v2"

const (
	GAMEMODE_IN_PROGRESS = iota
	GAMEMODE_WAITING_FOR_CHILD
	GAMEMODE_WON
	GAMEMODE_LOST
	GAMEMODE_DRAW
)

const (
	SCENE_CHANGE_TO_INDEX = iota
	SCENE_CHANGE_TO_PARENT
	SCENE_CHANGE_TO_CHILD
)

type GameStateData struct {
	GameState          int
	SceneChangeMode    int
	SceneChangeToIndex int
}

type GameMode interface {
	Init()
	Update() GameStateData
	Draw(screenRef *ebiten.Image)
	SceneTransitionHandler(previousGameStateData GameStateData)
}
