package main

type EntityComponents struct {
	totalNumEntities int
	playerEntityID   int
	cameraEntityID   int
	cameraData       CameraData
	tilemap          Tilemap
	enemyEntityIDs   []int
	itemEntityIDs    []int
	sprites          []Sprite
	positions        []Vector2
}

func CreateEntityComponents(totalNumEntitiesToCreate int) *EntityComponents {
	return &EntityComponents{
		totalNumEntities: totalNumEntitiesToCreate,
		sprites:          make([]Sprite, totalNumEntitiesToCreate),
		positions:        make([]Vector2, totalNumEntitiesToCreate)}
}
