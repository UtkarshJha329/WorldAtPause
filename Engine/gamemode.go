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

type SceneChangeData struct {
	SceneChangeMode              int
	SceneChangeToIndex           int
	QuestIssuedDuringSceneChange QuestData
}

type GameStateData struct {
	GameState       int
	SceneChangeData SceneChangeData
	GameStatsData   GameStatsData
}

type GameMode interface {
	Init()
	Update()
	Draw(screenRef *ebiten.Image)
	SceneTransitionHandler(previousGameStateData GameStateData)
}

//func (GameMode *GameMode)	Init() {}
//func (GameMode *GameMode)	Update() {}
//func (GameMode *GameMode)	Draw(screenRef *ebiten.Image) {}
//func (GameMode *GameMode)	SceneTransitionHandler(previousGameStateData GameStateData) {}
