package Engine

type Room struct {
	EntitiesInThisRoom []int
	PhysicsEntityIDs   []int
	EnemyEntityIDs     []int
	ObstacleEntityIDs  []int
	ItemEntityIDs      []int
}

type Scene struct {
	SceneType                string
	GameMode                 GameMode
	EntityComponentsForScene *EntityComponents
	RoomsData                map[Vector2Int]*Room
}

func CreateSceneWithNumEntities(numEntitiesToCreateInScene int) *Scene {
	return &Scene{
		EntityComponentsForScene: CreateEntityComponents(numEntitiesToCreateInScene),
		RoomsData:                make(map[Vector2Int]*Room),
	}
}
