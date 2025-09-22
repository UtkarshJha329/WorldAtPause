package Engine

import "time"

type Room struct {
	EntitiesInThisRoom []int
	EnemyEntityIDs     []int
	TriggerEntityIDs   []int
	ObstacleEntityIDs  []int
	ItemEntityIDs      []int
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
	UIRectSprites            map[string]*Sprite

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
		UIRectSprites:            make(map[string]*Sprite),
		lastUpdatedTime:          time.Now(),
	}
}
