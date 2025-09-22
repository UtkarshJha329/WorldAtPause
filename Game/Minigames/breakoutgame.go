package Minigames

import (
	"WorldAtPause/Engine"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type BreakoutGameMode struct {
	World    *Engine.World
	SceneRef *Engine.Scene

	ballInputDirection      Engine.Vector2
	ballMoveAmountPerFrame  float64
	numBricks               int
	failureTriggerTriggered bool
}

func (breakoutGameMode *BreakoutGameMode) Init() {

	sceneIndex := breakoutGameMode.World.SceneIndexByName[breakoutGameMode.SceneRef.SceneName]
	Engine.ReloadSceneWithSceneData(breakoutGameMode.SceneRef, &Engine.ScenesData[sceneIndex])

	breakoutGameMode.ballInputDirection = Engine.Vector2{X: -1.0, Y: -1.0}
	breakoutGameMode.ballMoveAmountPerFrame = 2.0
	breakoutGameMode.numBricks = 5.0

}

func (breakoutGameMode *BreakoutGameMode) Update() {

	curSceneRef := breakoutGameMode.SceneRef
	entityComponentsRef := breakoutGameMode.SceneRef.EntityComponentsForScene
	curRoomIndex := Engine.Vector2Int{X: 0, Y: 0}

	paddleMoveAmountPerFrame := 4.0
	ballEntityID := curSceneRef.EntityIDsByName["Breakout Ball"]
	ballPos := curSceneRef.EntityComponentsForScene.Positions[ballEntityID]
	ballCollisionShapeRef := curSceneRef.EntityComponentsForScene.CollisionShapes[ballEntityID]

	entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID].Y = 240 - 20

	inputX := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputX += 1.0
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputX -= 1.0
	}

	// playerBottomPos := curSceneRef.EntityComponentsForScene.Positions[entityComponentsRef.PlayerEntityID].Y + curSceneRef.EntityComponentsForScene.CollisionShapes[entityComponentsRef.PlayerEntityID].GetBoundingBoxDims().Y
	paddleCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
		CollideWithTiles:     true,
		CollideWithObstacles: false,
		SlideWhenCollide:     false,
		CollideWithPlayer:    false,
		MovementLock:         Engine.Vector2{X: 1.0, Y: 0.0},
	}

	// fmt.Println("Paddle collides with obstacles : ", paddleCollideAndMoveParameters.CollideWithObstacles)

	paddleInputDirection := Engine.Vector2{X: inputX, Y: 0.0}
	paddleTotalMoveAmount := Engine.Multiply_Float_Vector2(paddleMoveAmountPerFrame, &paddleInputDirection)

	breakoutGameMode.SceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, entityComponentsRef.PlayerEntityID, paddleInputDirection, paddleTotalMoveAmount, paddleMoveAmountPerFrame, paddleCollideAndMoveParameters)

	if breakoutGameMode.numBricks <= 0 {
		// breakoutGameMode.ballInputDirection = Engine.Vector2{X: 0.0, Y: 0.0}
	}

	ballInputDirection := breakoutGameMode.ballInputDirection
	ballTotalMoveAmount := Engine.Multiply_Float_Vector2(breakoutGameMode.ballMoveAmountPerFrame, &ballInputDirection)

	ballCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
		CollideWithTiles:     true,
		CollideWithObstacles: true,
		SlideWhenCollide:     false,
		CollideWithPlayer:    true,
		MovementLock:         Engine.Vector2{X: 1.0, Y: 1.0},
	}

	failureTriggerEntityID := curSceneRef.EntityIDsByName["Breakout Ball Failure Trigger"]
	failureTriggerPos := curSceneRef.EntityComponentsForScene.Positions[failureTriggerEntityID]
	failureTriggerColliderRef := curSceneRef.EntityComponentsForScene.CollisionShapes[failureTriggerEntityID]

	if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(&ballPos, ballCollisionShapeRef, &failureTriggerPos, failureTriggerColliderRef); collided {

		breakoutGameMode.failureTriggerTriggered = true
		ballCollideAndMoveParameters.MovementLock = Engine.Vector2{X: 0.0, Y: 0.0}
		curSceneRef.EntityComponentsForScene.Positions[ballEntityID] = Engine.Vector2{X: 160, Y: 120}
	}

	if !breakoutGameMode.failureTriggerTriggered {
		ballCollisions := breakoutGameMode.SceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, ballEntityID, ballInputDirection, ballTotalMoveAmount, breakoutGameMode.ballMoveAmountPerFrame, ballCollideAndMoveParameters)

		ballCollidedWithTilemap := ballCollisions.TilemapXMoveCollisionResult.Collided || ballCollisions.TilemapYMoveCollisionResult.Collided

		ballXCollidedWithObstacle := ballCollisions.ObstacleXMoveCollisionResult.CollidedWithObstacle && ballCollisions.ObstacleXMoveCollisionResult.CollisionPenetrationAmount != 0.0
		ballYCollidedWithObstacle := ballCollisions.ObstacleYMoveCollisionResult.CollidedWithObstacle && ballCollisions.ObstacleYMoveCollisionResult.CollisionPenetrationAmount != 0.0
		ballCollidedWithObstacles := ballXCollidedWithObstacle || ballYCollidedWithObstacle

		// if breakoutGameMode.numBricks > 0 && (ballCollidedWithTilemap || ballCollidedWithObstacles) {
		if ballCollidedWithTilemap || ballCollidedWithObstacles {

			ballCollidedTileNormal := Engine.Add_Vector2(&ballCollisions.TilemapXMoveCollisionResult.CollisionTileNormal, &ballCollisions.TilemapYMoveCollisionResult.CollisionTileNormal)
			ballCollidedObstacleNormal := Engine.Add_Vector2(&ballCollisions.ObstacleXMoveCollisionResult.CollisionNormal, &ballCollisions.ObstacleYMoveCollisionResult.CollisionNormal)

			useNormal := ballCollidedTileNormal
			if !ballCollidedWithTilemap {
				useNormal = ballCollidedObstacleNormal
			}
			useNormal = Engine.Normalise_Vector2(&useNormal)

			breakoutGameMode.ballInputDirection = Engine.Reflect_Vector2(&breakoutGameMode.ballInputDirection, &useNormal)

			if ballCollidedWithObstacles {
				if ballXCollidedWithObstacle {
					obstacleEntityID := ballCollisions.ObstacleXMoveCollisionResult.CollidedWithObstacleEntityID
					if entityComponentsRef.PlayerEntityID != obstacleEntityID {
						if !entityComponentsRef.EntityDead[obstacleEntityID] {
							breakoutGameMode.numBricks -= 1
						}
						entityComponentsRef.EntityDead[obstacleEntityID] = true
					}
				}
				if ballYCollidedWithObstacle {
					obstacleEntityID := ballCollisions.ObstacleYMoveCollisionResult.CollidedWithObstacleEntityID
					if entityComponentsRef.PlayerEntityID != obstacleEntityID {
						if !entityComponentsRef.EntityDead[obstacleEntityID] {
							breakoutGameMode.numBricks -= 1
						}
						entityComponentsRef.EntityDead[obstacleEntityID] = true
					}
				}
			}

		}
	}

	// fmt.Println("Finished breakout update.")

	if !breakoutGameMode.failureTriggerTriggered && breakoutGameMode.numBricks > 0 {
		curSceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_IN_PROGRESS
	} else if breakoutGameMode.numBricks <= 0 {
		curSceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_WON
		curSceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	} else if breakoutGameMode.failureTriggerTriggered {
		curSceneRef.SceneGameStateData.GameState = Engine.GAMEMODE_LOST
		curSceneRef.SceneGameStateData.SceneChangeData.SceneChangeMode = Engine.SCENE_CHANGE_TO_PARENT
	}
}

func (breakoutGameMode *BreakoutGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(breakoutGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(breakoutGameMode.SceneRef, screenRef)

	if breakoutGameMode.numBricks <= 0 {
		textEntity := breakoutGameMode.SceneRef.EntityIDsByName["Breakout You Win Text"]
		op := &text.DrawOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Translate(125.0, 100.0)
		currentText := breakoutGameMode.SceneRef.Texts[breakoutGameMode.SceneRef.EntityComponentsForScene.Texts[textEntity].TextName]
		text.Draw(screenRef, currentText, breakoutGameMode.SceneRef.EntityComponentsForScene.Texts[textEntity].Font.Face, op)
	} else if breakoutGameMode.failureTriggerTriggered {
		textEntity := breakoutGameMode.SceneRef.EntityIDsByName["Breakout You Lose Text"]
		op := &text.DrawOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Translate(35.0, 100.0)
		currentText := breakoutGameMode.SceneRef.Texts[breakoutGameMode.SceneRef.EntityComponentsForScene.Texts[textEntity].TextName]
		text.Draw(screenRef, currentText, breakoutGameMode.SceneRef.EntityComponentsForScene.Texts[textEntity].Font.Face, op)

	}

}

func (breakoutGameMode *BreakoutGameMode) SceneTransitionHandler(previousGameStateData Engine.GameStateData) {

}
