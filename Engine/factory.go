package Engine

type FactoryType int

const (
	FactoryType_EnemyEntities FactoryType = iota
	FactoryType_InvisibleTriggerEntities
	FactoryType_VisibleTriggerEntities
	FactoryType_ObstacleEntities
	FactoryType_ItemEntities
)

var FactoryStateFromString = map[string]FactoryType{
	"FactoryType_EnemyEntities":            FactoryType_EnemyEntities,
	"FactoryType_InvisibleTriggerEntities": FactoryType_InvisibleTriggerEntities,
	"FactoryType_VisibleTriggerEntities":   FactoryType_VisibleTriggerEntities,
	"FactoryType_ObstacleEntities":         FactoryType_ObstacleEntities,
	"FactoryType_ItemEntities":             FactoryType_ItemEntities,
}

type Factory struct {
	EntitiesPool Pool[int]
	EntityBegin  int
	EntityEnd    int
}

type FactoryAssetData struct {
	FactoryName                     string `json:"FactoryName"`
	FactoryType                     string `json:"FactoryType"`
	FactoryTemplateEntityPrefabName string `json:"FactoryTemplateEntityPrefabName"`
	FactoryNumEntities              int    `json:"FactoryNumEntities"`
}

func (factory *Factory) InitFactory(entityComponentsRef *EntityComponents, curSceneRef *Scene, currentRoom *Room, factoryName string, totalNumEntitiesInFactory int, baseEntityToCopy int, factoryType FactoryType) {

	factoryEntitiesDead := make([]bool, totalNumEntitiesInFactory)
	factoryEntitiesPositions := make([]Vector2, totalNumEntitiesInFactory)
	factoryEntitiesSprites := make([]Sprite, totalNumEntitiesInFactory)
	factoryEntitiesCollisionShapes := make([]CollisionShape, totalNumEntitiesInFactory)
	factoryEntitiesTexts := make([]Text, totalNumEntitiesInFactory)

	factory.EntityBegin = entityComponentsRef.TotalNumEntities
	entityComponentsRef.TotalNumEntities += totalNumEntitiesInFactory
	factory.EntityEnd = entityComponentsRef.TotalNumEntities - 1

	entityComponentsRef.EntityDead = append(entityComponentsRef.EntityDead, factoryEntitiesDead...)
	entityComponentsRef.Positions = append(entityComponentsRef.Positions, factoryEntitiesPositions...)
	entityComponentsRef.Sprites = append(entityComponentsRef.Sprites, factoryEntitiesSprites...)
	entityComponentsRef.CollisionShapes = append(entityComponentsRef.CollisionShapes, factoryEntitiesCollisionShapes...)
	entityComponentsRef.Texts = append(entityComponentsRef.Texts, factoryEntitiesTexts...)

	factory.EntitiesPool.InitPool(factoryName+" Entities Pool", totalNumEntitiesInFactory)

	for i := factory.EntityBegin; i <= factory.EntityEnd; i++ {
		currentEntity := i
		factory.EntitiesPool.GetAnUnusedItemFromPool().Item = currentEntity

		entityComponentsRef.EntityDead[currentEntity] = true
		entityComponentsRef.Positions[currentEntity] = entityComponentsRef.Positions[baseEntityToCopy]

		entityComponentsRef.Sprites[currentEntity] = Sprite{
			Image:                 entityComponentsRef.Sprites[baseEntityToCopy].Image,
			RenderRectStart:       entityComponentsRef.Sprites[baseEntityToCopy].RenderRectStart,
			RenderRectEnd:         entityComponentsRef.Sprites[baseEntityToCopy].RenderRectEnd,
			CurrentAnimationIndex: 0,
		}
		for i := 0; i < len(entityComponentsRef.Sprites[baseEntityToCopy].Animations); i++ {
			curAnimation := Animation{
				currentFrameCounter: 0,
				animationFramesData: entityComponentsRef.Sprites[baseEntityToCopy].Animations[i].animationFramesData,
				AnimationTimer: curSceneRef.AnimationsTimerSystem.SetTimerFromPoolWithDurationLoopAndFunc(entityComponentsRef.Sprites[baseEntityToCopy].Animations[i].AnimationTimer.Item.TimerTotalDuration, true, func() {

					curSpriteRef := &entityComponentsRef.Sprites[currentEntity]

					curAnimationFrameIndex := int(curSpriteRef.Animations[curSpriteRef.CurrentAnimationIndex].currentFrameCounter) % len(curSpriteRef.Animations[curSpriteRef.CurrentAnimationIndex].animationFramesData)
					curAnimationFrameData := curSpriteRef.Animations[curSpriteRef.CurrentAnimationIndex].animationFramesData[curAnimationFrameIndex]

					curSpriteRef.RenderRectStart = curAnimationFrameData.StartFramePos
					curSpriteRef.RenderRectEnd = curAnimationFrameData.EndFramePos

					curSpriteRef.Animations[curSpriteRef.CurrentAnimationIndex].currentFrameCounter++
				}),
			}
			entityComponentsRef.Sprites[currentEntity].Animations = append(entityComponentsRef.Sprites[currentEntity].Animations, curAnimation)
			entityComponentsRef.Sprites[currentEntity].Animations[i].AnimationTimer.Item.PauseTimer()
		}

		entityComponentsRef.CollisionShapes[currentEntity] = entityComponentsRef.CollisionShapes[baseEntityToCopy]
		entityComponentsRef.Texts[currentEntity] = entityComponentsRef.Texts[baseEntityToCopy]

		// Current Room Entities
		currentRoom.EntitiesInThisRoom = append(currentRoom.EntitiesInThisRoom, currentEntity)

		switch factoryType {
		case FactoryType_EnemyEntities:

			currentRoom.EnemyEntityIDs = append(currentRoom.EnemyEntityIDs, currentEntity)

		case FactoryType_InvisibleTriggerEntities:

			currentRoom.InvisibleTriggerEntityIDs = append(currentRoom.InvisibleTriggerEntityIDs, currentEntity)

		case FactoryType_VisibleTriggerEntities:

			currentRoom.VisibleTriggerEntityIDs = append(currentRoom.VisibleTriggerEntityIDs, currentEntity)

		case FactoryType_ObstacleEntities:

			currentRoom.ObstacleEntityIDs = append(currentRoom.ObstacleEntityIDs, currentEntity)

		case FactoryType_ItemEntities:

			currentRoom.ItemEntityIDs = append(currentRoom.ItemEntityIDs, currentEntity)

		}

	}
}
