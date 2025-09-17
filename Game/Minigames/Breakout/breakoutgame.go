package Games

import (
	"WorldAtPause/Engine"

	"github.com/hajimehoshi/ebiten/v2"
)

type BreakoutGameMode struct {
	SceneRef *Engine.Scene
}

func (breakoutGameMode *BreakoutGameMode) Init() {

}

func (breakoutGameMode *BreakoutGameMode) Update() {

}

func (breakoutGameMode *BreakoutGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(breakoutGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(breakoutGameMode.SceneRef, screenRef)

}
