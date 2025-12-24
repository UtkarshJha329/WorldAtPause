package Engine

import (
	"time"
)

type Room struct {
	EntitiesInThisRoom        []int
	EnemyEntityIDs            []int
	InvisibleTriggerEntityIDs []int
	VisibleTriggerEntityIDs   []int
	ObstacleEntityIDs         []int
	ItemEntityIDs             []int

	Factories []*Factory
}

type Scene struct {
	SceneName                string
	GameMode                 GameMode
	SceneGameStateData       GameStateData
	EntityComponentsForScene *EntityComponents
	RoomsData                map[Vector2Int]*Room
	Fonts                    map[string]*Font
	Texts                    map[string]string
	EntityIDsByName          map[string]int

	AnimationsTimerSystem TimerSystem

	lastUpdatedTime time.Time
	elapsed         time.Duration
}

func CreateSceneWithNumEntities(numEntitiesToCreateInScene int) *Scene {
	return &Scene{
		EntityComponentsForScene: CreateEntityComponents(numEntitiesToCreateInScene),
		RoomsData:                make(map[Vector2Int]*Room),
		Fonts:                    make(map[string]*Font),
		Texts:                    make(map[string]string),
		EntityIDsByName:          make(map[string]int),
		lastUpdatedTime:          time.Now(),
	}
}

func (curSceneRef *Scene) SpawnChildSceneWithQuest(childSceneIndex int, questToIssue QuestData) {

	curSceneRef.SceneGameStateData.GameState = GAMEMODE_WAITING_FOR_CHILD
	curSceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = SCENE_CHANGE_TO_CHILD
	curSceneRef.SceneGameStateData.SceneChangeData.SceneChangeToIndex = childSceneIndex

	curSceneRef.SceneGameStateData.SceneChangeData.QuestIssuedDuringSceneChange = questToIssue

}
