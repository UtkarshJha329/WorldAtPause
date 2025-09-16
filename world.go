package main

type World struct {
	currentSceneIndex int
	scenes            []*Scene
}

func CreateWorldWithNumScenes(totalNumScenesToCreate int) *World {
	return &World{
		currentSceneIndex: 0,
		scenes:            make([]*Scene, totalNumScenesToCreate),
	}
}
