package Minigames

import (
	"WorldAtPause/Engine"

	"github.com/hajimehoshi/ebiten/v2"
)

type BreakoutGameMode struct {
	SceneRef *Engine.Scene

	ballVelocity Engine.Vector2
}

func (breakoutGameMode *BreakoutGameMode) Init() {
	breakoutGameMode.ballVelocity = Engine.Vector2{X: -1.0, Y: -1.0}
}

func (breakoutGameMode *BreakoutGameMode) Update() {

	roomIndex := Engine.Vector2Int{X: 0, Y: 0}
	ballEntityID := breakoutGameMode.SceneRef.RoomsData[roomIndex].ObstacleEntityIDs[0]

	curSceneRef := breakoutGameMode.SceneRef
	entityComponentsRef := curSceneRef.EntityComponentsForScene

	paddleMoveAmountPerFrame := 2.0

	inputX := 0.0

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputX += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputX -= 1
	}

	paddleTotalMoveAmount := Engine.Vector2{X: inputX * paddleMoveAmountPerFrame, Y: 0.0}
	curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(roomIndex, entityComponentsRef.PlayerEntityID, Engine.Vector2{X: inputX, Y: 0.0}, paddleTotalMoveAmount, paddleMoveAmountPerFrame)
	// paddleMovementCollisionResult := curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(roomIndex, entityComponentsRef.PlayerEntityID, Engine.Vector2{X: inputX, Y: 0.0}, paddleTotalMoveAmount, paddleMoveAmountPerFrame)
	// paddleCollidedWithObstacle := paddleMovementCollisionResult.ObstacleXMoveCollisionResult.CollidedWithObstacle || paddleMovementCollisionResult.ObstacleYMoveCollisionResult.CollidedWithObstacle

	ballInputDirection := Engine.Vector2{X: breakoutGameMode.ballVelocity.X, Y: breakoutGameMode.ballVelocity.Y}
	collisionResult := curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(roomIndex, ballEntityID, ballInputDirection, breakoutGameMode.ballVelocity, 1.0)

	collidedObstacle := collisionResult.ObstacleXMoveCollisionResult.CollidedWithObstacle || collisionResult.ObstacleYMoveCollisionResult.CollidedWithObstacle
	collidedTilemap := collisionResult.TilemapXMoveCollisionResult.Collided || collisionResult.TilemapYMoveCollisionResult.Collided

	if collidedObstacle || collidedTilemap {

		obstacleCollisionNormalTotal := Engine.Add_Vector2(&collisionResult.ObstacleXMoveCollisionResult.CollisionNormalFromObstacle, &collisionResult.ObstacleYMoveCollisionResult.CollisionNormalFromObstacle)
		tilemapCollisionNormalTotal := Engine.Add_Vector2(&collisionResult.TilemapXMoveCollisionResult.CollisionTileNormal, &collisionResult.TilemapYMoveCollisionResult.CollisionTileNormal)

		usageNormal := Engine.Normalise_Vector2(&tilemapCollisionNormalTotal)
		if collidedObstacle {
			usageNormal = Engine.Normalise_Vector2(&obstacleCollisionNormalTotal)
		}

		breakoutGameMode.ballVelocity = Engine.Reflect_Vector2(&breakoutGameMode.ballVelocity, &usageNormal)
	}
	// else if paddleCollidedWithObstacle {
	// 	paddleNormal := Engine.Vector2{X: 0.0, Y: -1.0}
	// 	breakoutGameMode.ballVelocity = Engine.Reflect_Vector2(&breakoutGameMode.ballVelocity, &paddleNormal)
	// }

}

func (breakoutGameMode *BreakoutGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(breakoutGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(breakoutGameMode.SceneRef, screenRef)

}
