package Engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	CurrentSceneIndex int
	Scenes            []*Scene
	SceneStack        SceneStack
	changedScenes     bool
}

func CreateWorldWithNumScenes(totalNumScenesToCreate int) *World {
	return &World{
		CurrentSceneIndex: 0,
		Scenes:            make([]*Scene, totalNumScenesToCreate),
		SceneStack:        SceneStack{SceneIndexStack: make([]int, 0)},
		changedScenes:     false,
	}
}

func (w *World) InitCurrentSceneGameMode() {

	entityComponentsRef := w.Scenes[w.CurrentSceneIndex].EntityComponentsForScene
	entityComponentsRef.CameraData.ScreenSize = Vector2{X: 320, Y: 240}

	w.Scenes[w.CurrentSceneIndex].GameMode.Init()

	w.Scenes[w.CurrentSceneIndex].SceneGameStateData.GameStatsData.TimeSinceLaunch = 0.0
}

func (w *World) UpdateCurrentSceneGameMode() {

	w.Scenes[w.CurrentSceneIndex].GameMode.Update()

	currentSceneGameState := w.Scenes[w.CurrentSceneIndex].SceneGameStateData

	if currentSceneGameState.GameState != GAMEMODE_IN_PROGRESS {

		switch currentSceneGameState.SceneChangeData.SceneChangeMode {

		case SCENE_CHANGE_TO_CHILD:
			w.SceneStack.Push(w.CurrentSceneIndex)
			w.SceneStack.Push(currentSceneGameState.SceneChangeData.SceneChangeToIndex)
			w.CurrentSceneIndex = currentSceneGameState.SceneChangeData.SceneChangeToIndex
			w.Scenes[w.CurrentSceneIndex].lastUpdatedTime = time.Now()
			w.InitCurrentSceneGameMode()
			w.Scenes[w.CurrentSceneIndex].GameMode.SceneTransitionHandler(currentSceneGameState)
			w.changedScenes = true

		case SCENE_CHANGE_TO_PARENT:

			// Remove previous scene index from top of stack
			w.SceneStack.Pop()
			// Now the current scene is the previous scene that spawned the current scene
			w.CurrentSceneIndex, _ = w.SceneStack.Pop()

			w.Scenes[w.CurrentSceneIndex].GameMode.SceneTransitionHandler(currentSceneGameState)
			w.changedScenes = true
		}
	}

	now := time.Now()
	dt := now.Sub(w.Scenes[w.CurrentSceneIndex].lastUpdatedTime)
	w.Scenes[w.CurrentSceneIndex].elapsed += dt
	w.Scenes[w.CurrentSceneIndex].SceneGameStateData.GameStatsData.TimeSinceLaunch += dt
	w.Scenes[w.CurrentSceneIndex].SceneGameStateData.GameStatsData.TimeLastFrame = dt

	w.Scenes[w.CurrentSceneIndex].lastUpdatedTime = time.Now()
}

func (w *World) DrawCurrentSceneGameMode(screenRef *ebiten.Image) {
	if !w.changedScenes {
		w.Scenes[w.CurrentSceneIndex].GameMode.Draw(screenRef)
	} else {
		w.changedScenes = false
	}
}

// SCENE STACK

type SceneStack struct {
	SceneIndexStack []int
}

func (sceneStack *SceneStack) IsEmpty() bool {
	return len(sceneStack.SceneIndexStack) == 0
}

func (sceneStack *SceneStack) Push(sceneIndex int) {
	sceneStack.SceneIndexStack = append(sceneStack.SceneIndexStack, sceneIndex)
}

func (sceneStack *SceneStack) Pop() (int, bool) {

	if sceneStack.IsEmpty() {
		return -1, false
	}

	index := len(sceneStack.SceneIndexStack) - 1
	sceneIndexToReturn := sceneStack.SceneIndexStack[index]
	sceneStack.SceneIndexStack = sceneStack.SceneIndexStack[:index]

	return sceneIndexToReturn, true
}

func (sceneStack *SceneStack) Peek() (int, bool) {
	if sceneStack.IsEmpty() {
		return -1, false
	}
	return sceneStack.SceneIndexStack[len(sceneStack.SceneIndexStack)-1], true
}

func (sceneStack *SceneStack) Flush() {
	sceneStack.SceneIndexStack = sceneStack.SceneIndexStack[:0]
}
