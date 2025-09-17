package Engine

import "github.com/hajimehoshi/ebiten/v2"

type World struct {
	CurrentSceneIndex int
	Scenes            []*Scene
}

func CreateWorldWithNumScenes(totalNumScenesToCreate int) *World {
	return &World{
		CurrentSceneIndex: 0,
		Scenes:            make([]*Scene, totalNumScenesToCreate),
	}
}

func (w *World) InitCurrentSceneGameMode() {
	w.Scenes[w.CurrentSceneIndex].GameMode.Init()
}

func (w *World) UpdateCurrentSceneGameMode() {
	w.Scenes[w.CurrentSceneIndex].GameMode.Update()
}

func (w *World) DrawCurrentSceneGameMode(screenRef *ebiten.Image) {
	w.Scenes[w.CurrentSceneIndex].GameMode.Draw(screenRef)
}
