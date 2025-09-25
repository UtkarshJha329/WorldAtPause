package Minigames

import (
	"WorldAtPause/Engine"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type SokobanGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	playerMoveAmountPerFrame float64

	boxCollectionTriggerEntityID int
	boxesLeftToBeCollected       int
}

func (sokobanGameMode *SokobanGameMode) Init() {

	sceneIndex := sokobanGameMode.World.SceneIndexByName[sokobanGameMode.SceneRef.SceneName]
	Engine.ReloadSceneWithSceneData(sokobanGameMode.SceneRef, &Engine.ScenesData[sceneIndex])

	sokobanGameMode.playerMoveAmountPerFrame = 4.0
	sokobanGameMode.boxCollectionTriggerEntityID = sokobanGameMode.SceneRef.EntityIDsByName["Sokoban Box Collection Trigger"]
	sokobanGameMode.boxesLeftToBeCollected = 5
}

func (sokobanGameMode *SokobanGameMode) Update() {

	curSceneRef := sokobanGameMode.SceneRef
	entityComponentsRef := curSceneRef.EntityComponentsForScene

	playerMoveAmountPerFrame := sokobanGameMode.playerMoveAmountPerFrame

	playerPosRef := &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)

	inputDirection := Engine.Vector2{X: 0.0, Y: 0.0}

	curSceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputDirection.X += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputDirection.X -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		inputDirection.Y += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		inputDirection.Y -= 1
	}

	if inputDirection.X != 0 || inputDirection.Y != 0 {

		normalisedInputDir := Engine.Normalise_Vector2(&inputDirection)
		totalMoveAmount := Engine.Multiply_Float_Vector2(playerMoveAmountPerFrame, &normalisedInputDir)

		playerCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
			CollideWithTiles:     true,
			CollideWithObstacles: true,
			SlideWhenCollide:     true,
			CollideWithPlayer:    false,
			MovementLock:         Engine.Vector2{X: 1.0, Y: 1.0},
		}

		playerMoveCollisionResult := curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, curSceneRef.EntityComponentsForScene.PlayerEntityID, inputDirection, totalMoveAmount, playerMoveAmountPerFrame, playerCollideAndMoveParameters)

		collidedWithObstacleEntityID := -1
		penetrationNormal := Engine.Vector2{X: 0.0, Y: 0.0}

		if playerMoveCollisionResult.ObstacleXMoveCollisionResult.CollidedWithObstacle {

			collidedWithObstacleEntityID = playerMoveCollisionResult.ObstacleXMoveCollisionResult.CollidedWithObstacleEntityID
			penetrationNormal = playerMoveCollisionResult.ObstacleXMoveCollisionResult.CollisionNormal

		} else if playerMoveCollisionResult.ObstacleYMoveCollisionResult.CollidedWithObstacle {

			collidedWithObstacleEntityID = playerMoveCollisionResult.ObstacleYMoveCollisionResult.CollidedWithObstacleEntityID
			penetrationNormal = playerMoveCollisionResult.ObstacleYMoveCollisionResult.CollisionNormal

		}

		if collidedWithObstacleEntityID != -1 {

			obstaclePos := entityComponentsRef.Positions[collidedWithObstacleEntityID]
			obstacleCollisionShapeRef := entityComponentsRef.CollisionShapes[collidedWithObstacleEntityID]

			obstacleCentre := Engine.Multiply_Float_Vector2(0.5, obstacleCollisionShapeRef.GetBoundingBoxDims())
			obstacleCentre = Engine.Add_Vector2(&obstaclePos, &obstacleCentre)

			obstacleCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
				CollideWithTiles:     true,
				CollideWithObstacles: true,
				SlideWhenCollide:     false,
				CollideWithPlayer:    true,
				MovementLock:         Engine.Vector2{X: math.Abs(inputDirection.X), Y: math.Abs(inputDirection.Y)},
			}

			inputX := 0.0
			if penetrationNormal.X != 0.0 {
				inputX = math.Copysign(1.0, penetrationNormal.X)
			}
			inputY := 0.0
			if penetrationNormal.Y != 0.0 {
				inputY = math.Copysign(1.0, penetrationNormal.Y)
			}
			obstacleInputdirection := Engine.Vector2{X: inputX, Y: inputY}
			obstacleInputdirection = Engine.Multiply_Float_Vector2(-1.0, &obstacleInputdirection)
			obstacleMoveAmountPerFrame := 1.0
			obstacleTotalMoveAmount := Engine.Multiply_Float_Vector2(obstacleMoveAmountPerFrame, &obstacleInputdirection)

			curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, collidedWithObstacleEntityID, obstacleInputdirection, obstacleTotalMoveAmount, obstacleMoveAmountPerFrame, obstacleCollideAndMoveParameters)

			triggerPosRef := &entityComponentsRef.Positions[sokobanGameMode.boxCollectionTriggerEntityID]
			triggerCollisionShapeRef := entityComponentsRef.CollisionShapes[sokobanGameMode.boxCollectionTriggerEntityID]
			if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(&obstaclePos, obstacleCollisionShapeRef, triggerPosRef, triggerCollisionShapeRef); collided {

				entityComponentsRef.EntityDead[collidedWithObstacleEntityID] = true
				sokobanGameMode.boxesLeftToBeCollected--
			}
		}
	}

	if sokobanGameMode.boxesLeftToBeCollected <= 0 {
		sokobanGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WON
		sokobanGameMode.SceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	} else {
		sokobanGameMode.SceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS
	}

}

func (sokobanGameMode *SokobanGameMode) Draw(screenRef *ebiten.Image) {
	Engine.DrawActiveRoomInScene(sokobanGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(sokobanGameMode.SceneRef, screenRef)
}

func (sokobanGameMode *SokobanGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

}
