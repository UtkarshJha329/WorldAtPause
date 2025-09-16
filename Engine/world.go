package Engine

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
