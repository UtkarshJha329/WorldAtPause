package main

type Room struct {
	entitiesInThisRoom []int
	enemyEntityIDs     []int
	obstacleEntityIDs  []int
	itemEntityIDs      []int
}

type Scene struct {
	entityComponentsForScene *EntityComponents
	roomsData                map[Vector2Int]*Room
}

func CreateSceneWithNumEntities(numEntitiesToCreateInScene int) *Scene {
	return &Scene{
		entityComponentsForScene: CreateEntityComponents(numEntitiesToCreateInScene),
		roomsData:                make(map[Vector2Int]*Room),
	}
}
