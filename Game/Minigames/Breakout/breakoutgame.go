package Games

import (
	"WorldAtPause/Engine"

	"github.com/hajimehoshi/ebiten/v2"
)

type BreakoutGameMode struct {
	SceneRef *Engine.Scene

	ballInputDirection Engine.Vector2
}

func (breakoutGameMode *BreakoutGameMode) Init() {
	breakoutGameMode.ballInputDirection = Engine.Vector2{X: 0.0, Y: 1.0}
}

func (breakoutGameMode *BreakoutGameMode) Update() {

	curSceneRef := breakoutGameMode.SceneRef
	entityComponentsRef := breakoutGameMode.SceneRef.EntityComponentsForScene
	curRoomIndex := Engine.Vector2Int{X: 0, Y: 0}

	paddleMoveAmountPerFrame := 2.0

	inputX := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputX += 1.0
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputX -= 1.0
	}

	paddleCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
		CollideWithTiles:     true,
		CollideWithObstacles: false,
		SlideWhenCollide:     false,
	}

	paddleInputDirection := Engine.Vector2{X: inputX, Y: 0.0}
	paddleTotalMoveAmount := Engine.Multiply_Float_Vector2(paddleMoveAmountPerFrame, &paddleInputDirection)

	breakoutGameMode.SceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, entityComponentsRef.PlayerEntityID, paddleInputDirection, paddleTotalMoveAmount, paddleMoveAmountPerFrame, paddleCollideAndMoveParameters)

	ballEntityID := curSceneRef.RoomsData[curRoomIndex].ObstacleEntityIDs[0]
	ballMoveAmountPerFrame := 1.0

	ballInputDirection := breakoutGameMode.ballInputDirection
	ballTotalMoveAmount := Engine.Multiply_Float_Vector2(ballMoveAmountPerFrame, &ballInputDirection)

	ballCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
		CollideWithTiles:     true,
		CollideWithObstacles: true,
		SlideWhenCollide:     false,
	}

	ballCollisions := breakoutGameMode.SceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, ballEntityID, ballInputDirection, ballTotalMoveAmount, ballMoveAmountPerFrame, ballCollideAndMoveParameters)

	ballCollidedWithTilemap := ballCollisions.TilemapXMoveCollisionResult.Collided || ballCollisions.TilemapYMoveCollisionResult.Collided
	ballCollidedWithObstacles := ballCollisions.ObstacleXMoveCollisionResult.CollidedWithObstacle || ballCollisions.ObstacleYMoveCollisionResult.CollidedWithObstacle

	if ballCollidedWithTilemap || ballCollidedWithObstacles {

		ballCollidedTileNormal := Engine.Add_Vector2(&ballCollisions.TilemapXMoveCollisionResult.CollisionTileNormal, &ballCollisions.TilemapYMoveCollisionResult.CollisionTileNormal)
		ballCollidedObstacleNormal := Engine.Add_Vector2(&ballCollisions.ObstacleXMoveCollisionResult.CollisionNormal, &ballCollisions.ObstacleYMoveCollisionResult.CollisionNormal)

		useNormal := ballCollidedTileNormal
		if !ballCollidedWithTilemap {
			useNormal = ballCollidedObstacleNormal
		}
		useNormal = Engine.Normalise_Vector2(&useNormal)

		breakoutGameMode.ballInputDirection = Engine.Reflect_Vector2(&breakoutGameMode.ballInputDirection, &useNormal)
	}

}

func (breakoutGameMode *BreakoutGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(breakoutGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(breakoutGameMode.SceneRef, screenRef)

}
