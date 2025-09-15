package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	entityComponentsRef *EntityComponents
}

func (g *Game) Update() error {

	// playerPos := g.entityComponentsRef.positions[g.entityComponentsRef.playerEntityID]

	playerPosRef := &g.entityComponentsRef.positions[g.entityComponentsRef.playerEntityID]
	playerCollisionShapeRef := g.entityComponentsRef.collisionShapes[g.entityComponentsRef.playerEntityID]
	playerMoveAmountPerFrame := 2.0

	inputDirection := Vector2{0.0, 0.0}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		inputDirection.x += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		inputDirection.x -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		inputDirection.y += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		inputDirection.y -= 1
	}

	if inputDirection.x != 0 || inputDirection.y != 0 {

		normalisedInputDir := Normalise_Vector2(&inputDirection)
		totalMoveAmount := Multiply_Float_Vector2(playerMoveAmountPerFrame, &normalisedInputDir)

		g.entityComponentsRef.MoveAndCollideEntityWithTilemapAndObstacles(g.entityComponentsRef.playerEntityID, inputDirection, totalMoveAmount, playerMoveAmountPerFrame)

		// playerToMoveXPos := playerPos.x + totalMoveAmount.x
		// moveDirectionCheckOffsetX := 0.0
		// if inputDirection.x > 0 {
		// 	moveDirectionCheckOffsetX = playerCollisionShapeRef.(*BoxCollider).size.x
		// }

		// collidingOnXWithTilemap := g.entityComponentsRef.tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{playerToMoveXPos + moveDirectionCheckOffsetX, playerPos.y})
		// collidingOnX := collidingOnXWithTilemap

		// maxMoveAmtOnX := 0.0
		// if !collidingOnXWithTilemap {
		// 	obstacleCollisionResult := g.entityComponentsRef.DoesEntityCollideWithObstacles(&Vector2{playerToMoveXPos, playerPos.y}, playerCollisionShapeRef)
		// 	if obstacleCollisionResult.collidedWithObstacle {
		// 		collidingOnX = true
		// 		if inputDirection.x > 0 {
		// 			maxMoveAmtOnX = obstacleCollisionResult.collidedWithObstaclePosition.x - (playerPos.x + playerCollisionShapeRef.(*BoxCollider).size.x)
		// 		} else if inputDirection.x < 0 {
		// 			maxMoveAmtOnX = playerPos.x - (obstacleCollisionResult.collidedWithObstaclePosition.x + obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.(*BoxCollider).size.x)
		// 		}
		// 	}
		// }

		// playerToMoveYPos := playerPos.y + totalMoveAmount.y
		// moveDirectionCheckOffsetY := 0.0
		// if inputDirection.y > 0 {
		// 	moveDirectionCheckOffsetY = playerCollisionShapeRef.(*BoxCollider).size.y
		// }

		// collidingOnYWithTilemap := g.entityComponentsRef.tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{playerPos.x, playerToMoveYPos + moveDirectionCheckOffsetY})
		// collidingOnY := collidingOnYWithTilemap

		// maxMoveAmtOnY := 0.0
		// if !collidingOnYWithTilemap {
		// 	obstacleCollisionResult := g.entityComponentsRef.DoesEntityCollideWithObstacles(&Vector2{playerPos.x, playerToMoveYPos}, playerCollisionShapeRef)
		// 	if obstacleCollisionResult.collidedWithObstacle {
		// 		collidingOnY = true
		// 		if inputDirection.y > 0 {
		// 			maxMoveAmtOnY = obstacleCollisionResult.collidedWithObstaclePosition.y - (playerPos.y + playerCollisionShapeRef.(*BoxCollider).size.y)
		// 		} else if inputDirection.y < 0 {
		// 			maxMoveAmtOnY = playerPos.y - (obstacleCollisionResult.collidedWithObstaclePosition.y + obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.(*BoxCollider).size.y)
		// 		}
		// 	}
		// }

		// if !collidingOnX && !collidingOnY {
		// 	playerPosRef.x += totalMoveAmount.x
		// 	playerPosRef.y += totalMoveAmount.y
		// } else if !collidingOnX && collidingOnY {
		// 	playerPosRef.x += inputDirection.x * playerMoveAmountPerFrame
		// 	playerPosRef.y += inputDirection.y * maxMoveAmtOnY
		// } else if collidingOnX && !collidingOnY {
		// 	playerPosRef.y += inputDirection.y * playerMoveAmountPerFrame
		// 	playerPosRef.x += inputDirection.x * maxMoveAmtOnX
		// }

	}

	stoppageDistanceFromPlayer := 32.0
	skeleMoveAmountPerFrame := 1.0
	for _, enemyEntityID := range g.entityComponentsRef.enemyEntityIDs {
		skelePosRef := &g.entityComponentsRef.positions[enemyEntityID]

		if DistanceSquare_Vector2(skelePosRef, playerPosRef) > math.Pow(stoppageDistanceFromPlayer, 2) {
			directionToPlayer := Subtract_Vector2(playerPosRef, skelePosRef)
			directionToPlayer = Normalise_Vector2(&directionToPlayer)

			totalDisplacement := Multiply_Float_Vector2(skeleMoveAmountPerFrame, &directionToPlayer)

			finalPosition := Add_Vector2(skelePosRef, &totalDisplacement)

			skelePosRef.x = finalPosition.x
			skelePosRef.y = finalPosition.y
		}

		skeletonCollisionShapeRef := g.entityComponentsRef.collisionShapes[enemyEntityID]
		if CollisionShapeOverlapsWithCollisionShape(skelePosRef, skeletonCollisionShapeRef, playerPosRef, playerCollisionShapeRef) {
			// fmt.Println("Skeleton is colliding with player!")
		}
	}

	for index, itemEntityID := range g.entityComponentsRef.itemEntityIDs {
		itemPosRef := &g.entityComponentsRef.positions[itemEntityID]

		itemCentre := Vector2{itemPosRef.x + 5.0, itemPosRef.y + 5.0}
		itemCollisionShapeRef := g.entityComponentsRef.collisionShapes[itemEntityID]

		if CollisionShapeOverlapsWithCollisionShape(&itemCentre, itemCollisionShapeRef, &g.entityComponentsRef.positions[g.entityComponentsRef.playerEntityID], g.entityComponentsRef.collisionShapes[g.entityComponentsRef.playerEntityID]) {
			g.entityComponentsRef.itemEntityIDs = append(g.entityComponentsRef.itemEntityIDs[:index], g.entityComponentsRef.itemEntityIDs[index+1:]...)
			fmt.Printf("Picked up item entity ID : %d\n", itemEntityID)
		}

	}

	currentLevelIndex := g.entityComponentsRef.tilemap.GetLevelIndexOfPosition(playerPosRef)
	playerInLevel := Vector2{float64(currentLevelIndex.x), float64(currentLevelIndex.y)}

	levelHalfSize := Multiply_Float_Vector2(0.5, &g.entityComponentsRef.tilemap.worldGridSize)
	currentLevelPos := Vector2{playerInLevel.x * g.entityComponentsRef.tilemap.worldGridSize.x, playerInLevel.y * g.entityComponentsRef.tilemap.worldGridSize.y}
	currentLevelCentre := Add_Vector2(&currentLevelPos, &levelHalfSize)

	CameraFollowTarget(currentLevelCentre, g.entityComponentsRef)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	// TODO : Sort drawings based on Layers order and sort order index and then draw all at once.

	drawImgOptions := ebiten.DrawImageOptions{}

	playerPos := g.entityComponentsRef.positions[g.entityComponentsRef.playerEntityID]
	playerInLevel := Vector2Int{int(playerPos.x) / int(g.entityComponentsRef.tilemap.worldGridSize.x), int(playerPos.y) / int(g.entityComponentsRef.tilemap.worldGridSize.y)}

	DrawTileMapLevel(g.entityComponentsRef.tilemap.levels[playerInLevel], g, screen, &drawImgOptions)

	for _, itemEntityID := range g.entityComponentsRef.itemEntityIDs {
		DrawEntityID(g, screen, &drawImgOptions, itemEntityID)
	}

	for _, enemyEntityID := range g.entityComponentsRef.enemyEntityIDs {
		DrawEntityID(g, screen, &drawImgOptions, enemyEntityID)
	}

	for _, obstacleEntityID := range g.entityComponentsRef.obstacleEntityIDs {
		DrawEntityID(g, screen, &drawImgOptions, obstacleEntityID)
	}

	DrawEntityID(g, screen, &drawImgOptions, g.entityComponentsRef.playerEntityID)

	playerCollisionShapeRef := g.entityComponentsRef.collisionShapes[g.entityComponentsRef.playerEntityID]
	playerCollisionShapePoints := playerCollisionShapeRef.GetCollisionPoints()
	topLeft := Add_Vector2(&(*playerCollisionShapePoints)[0], &playerPos)
	bottomRight := Add_Vector2(&(*playerCollisionShapePoints)[2], &playerPos)
	size := Vector2{bottomRight.x - topLeft.x, bottomRight.y - topLeft.y}

	vector.StrokeRect(screen, float32(topLeft.x), float32(topLeft.y), float32(size.x), float32(size.y), 1.0, color.Black, false)

	for _, enemyEntityID := range g.entityComponentsRef.enemyEntityIDs {
		skelePosRef := &g.entityComponentsRef.positions[enemyEntityID]

		skeletonCollisionShapeRef := g.entityComponentsRef.collisionShapes[enemyEntityID]
		vector.StrokeRect(screen, float32(skelePosRef.x), float32(skelePosRef.y), float32(skeletonCollisionShapeRef.(*BoxCollider).size.x), float32(skeletonCollisionShapeRef.(*BoxCollider).size.y), 1.0, color.Black, false)
	}

	for _, obstacleIndex := range g.entityComponentsRef.obstacleEntityIDs {

		obstaclePosRef := &g.entityComponentsRef.positions[obstacleIndex]
		obstacleCollisionShapeRef := g.entityComponentsRef.collisionShapes[obstacleIndex]

		vector.StrokeRect(screen, float32(obstaclePosRef.x), float32(obstaclePosRef.y), float32(obstacleCollisionShapeRef.(*BoxCollider).size.x), float32(obstacleCollisionShapeRef.(*BoxCollider).size.y), 1.0, color.Black, false)
	}

	for _, itemEntityID := range g.entityComponentsRef.itemEntityIDs {

		itemPosRef := &g.entityComponentsRef.positions[itemEntityID]
		itemCentre := Vector2{itemPosRef.x + 5.0, itemPosRef.y + 5.0}
		itemCollisionShapeRef := g.entityComponentsRef.collisionShapes[itemEntityID]

		vector.StrokeCircle(screen, float32(itemCentre.x), float32(itemCentre.y), float32(itemCollisionShapeRef.(*CircleCollider).radius), 1.0, color.Black, false)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.entityComponentsRef.cameraData.screenSize.x), int(g.entityComponentsRef.cameraData.screenSize.y)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ninja!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := Game{
		entityComponentsRef: CreateAndPopulateEntitiesAndComponentsFromGameData("Assets/AssetsData.json"),
	}

	game.entityComponentsRef.cameraData.screenSize = Vector2{320, 240}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
