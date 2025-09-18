package MainGame

import (
	"WorldAtPause/Engine"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type MainGameMode struct {
	SceneRef *Engine.Scene
}

func (mainGameMode *MainGameMode) Init() {

}

func (mainGameMode *MainGameMode) Update() {

	curSceneRef := mainGameMode.SceneRef
	entityComponentsRef := curSceneRef.EntityComponentsForScene

	playerMoveAmountPerFrame := 2.0

	// collideAndSlide := true

	playerPosRef := &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)
	curRoom, curRoomHasSomeData := curSceneRef.RoomsData[curRoomIndex]

	inputDirection := Engine.Vector2{X: 0.0, Y: 0.0}

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

		curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, curSceneRef.EntityComponentsForScene.PlayerEntityID, inputDirection, totalMoveAmount, playerMoveAmountPerFrame, playerCollideAndMoveParameters)
	}

	if curRoomHasSomeData {
		playerCollisionShapeRef := entityComponentsRef.CollisionShapes[entityComponentsRef.PlayerEntityID]
		stoppageDistanceFromPlayer := 32.0
		skeleMoveAmountPerFrame := 1.0
		for _, enemyEntityID := range curRoom.EnemyEntityIDs {
			skelePosRef := &entityComponentsRef.Positions[enemyEntityID]

			if Engine.DistanceSquare_Vector2(skelePosRef, playerPosRef) > math.Pow(stoppageDistanceFromPlayer, 2) {
				directionToPlayer := Engine.Subtract_Vector2(playerPosRef, skelePosRef)
				directionToPlayerNormalised := Engine.Normalise_Vector2(&directionToPlayer)

				totalDisplacement := Engine.Multiply_Float_Vector2(skeleMoveAmountPerFrame, &directionToPlayerNormalised)

				skeletonCollideAndMoveParameters := Engine.CollideAndMoveCollisionParameters{
					CollideWithTiles:     true,
					CollideWithObstacles: true,
					SlideWhenCollide:     true,
					CollideWithPlayer:    true,
					MovementLock:         Engine.Vector2{X: 1.0, Y: 1.0},
				}

				curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, enemyEntityID, directionToPlayerNormalised, totalDisplacement, skeleMoveAmountPerFrame, skeletonCollideAndMoveParameters)
			}

			skeletonCollisionShapeRef := entityComponentsRef.CollisionShapes[enemyEntityID]
			if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(skelePosRef, skeletonCollisionShapeRef, playerPosRef, playerCollisionShapeRef); collided {
				// fmt.Println("Skeleton is colliding with player!")
			}
		}

		for index, itemEntityID := range curRoom.ItemEntityIDs {

			itemPosRef := &entityComponentsRef.Positions[itemEntityID]
			itemCollisionShapeRef := entityComponentsRef.CollisionShapes[itemEntityID]

			if _, _, _, collided := Engine.CollisionShapeOverlapsWithCollisionShape(itemPosRef, itemCollisionShapeRef, &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID], entityComponentsRef.CollisionShapes[entityComponentsRef.PlayerEntityID]); collided {
				curRoom.ItemEntityIDs = append(curRoom.ItemEntityIDs[:index], curRoom.ItemEntityIDs[index+1:]...)
				fmt.Printf("Picked up item entity ID : %d\n", itemEntityID)
			}

		}
	}
	currentLevelIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)
	playerInLevel := Engine.Vector2{X: float64(currentLevelIndex.X), Y: float64(currentLevelIndex.Y)}

	levelHalfSize := Engine.Multiply_Float_Vector2(0.5, &entityComponentsRef.Tilemap.WorldGridSize)
	currentLevelPos := Engine.Vector2{X: playerInLevel.X * entityComponentsRef.Tilemap.WorldGridSize.X, Y: playerInLevel.Y * entityComponentsRef.Tilemap.WorldGridSize.Y}
	currentLevelCentre := Engine.Add_Vector2(&currentLevelPos, &levelHalfSize)

	Engine.CameraFollowTarget(currentLevelCentre, entityComponentsRef)
}

func (mainGameMode *MainGameMode) Draw(screenRef *ebiten.Image) {

	Engine.DrawActiveRoomInScene(mainGameMode.SceneRef, screenRef)
	Engine.DrawActiveRoomInSceneColliders(mainGameMode.SceneRef, screenRef)
}
