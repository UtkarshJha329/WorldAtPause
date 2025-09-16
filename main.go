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
	world *World
}

func (g *Game) Update() error {

	// playerPos := entityComponentsRef.positions[entityComponentsRef.playerEntityID]

	curSceneRef := g.world.scenes[g.world.currentSceneIndex]
	entityComponentsRef := curSceneRef.entityComponentsForScene

	playerMoveAmountPerFrame := 2.0

	playerPosRef := &entityComponentsRef.positions[entityComponentsRef.playerEntityID]
	curRoomIndex := entityComponentsRef.tilemap.GetLevelIndexOfPosition(playerPosRef)
	curRoom, curRoomHasSomeData := curSceneRef.roomsData[curRoomIndex]

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

		curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, curSceneRef.entityComponentsForScene.playerEntityID, inputDirection, totalMoveAmount, playerMoveAmountPerFrame)
	}

	if curRoomHasSomeData {
		playerCollisionShapeRef := entityComponentsRef.collisionShapes[entityComponentsRef.playerEntityID]
		stoppageDistanceFromPlayer := 32.0
		skeleMoveAmountPerFrame := 1.0
		for _, enemyEntityID := range curRoom.enemyEntityIDs {
			skelePosRef := &entityComponentsRef.positions[enemyEntityID]

			if DistanceSquare_Vector2(skelePosRef, playerPosRef) > math.Pow(stoppageDistanceFromPlayer, 2) {
				directionToPlayer := Subtract_Vector2(playerPosRef, skelePosRef)
				directionToPlayerNormalised := Normalise_Vector2(&directionToPlayer)

				totalDisplacement := Multiply_Float_Vector2(skeleMoveAmountPerFrame, &directionToPlayerNormalised)

				curSceneRef.CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex, enemyEntityID, directionToPlayerNormalised, totalDisplacement, skeleMoveAmountPerFrame)
			}

			skeletonCollisionShapeRef := entityComponentsRef.collisionShapes[enemyEntityID]
			if CollisionShapeOverlapsWithCollisionShape(skelePosRef, skeletonCollisionShapeRef, playerPosRef, playerCollisionShapeRef) {
				// fmt.Println("Skeleton is colliding with player!")
			}
		}

		for index, itemEntityID := range curRoom.itemEntityIDs {

			itemPosRef := &entityComponentsRef.positions[itemEntityID]
			itemCollisionShapeRef := entityComponentsRef.collisionShapes[itemEntityID]

			if CollisionShapeOverlapsWithCollisionShape(itemPosRef, itemCollisionShapeRef, &entityComponentsRef.positions[entityComponentsRef.playerEntityID], entityComponentsRef.collisionShapes[entityComponentsRef.playerEntityID]) {
				curRoom.itemEntityIDs = append(curRoom.itemEntityIDs[:index], curRoom.itemEntityIDs[index+1:]...)
				fmt.Printf("Picked up item entity ID : %d\n", itemEntityID)
			}

		}
	}
	currentLevelIndex := entityComponentsRef.tilemap.GetLevelIndexOfPosition(playerPosRef)
	playerInLevel := Vector2{float64(currentLevelIndex.x), float64(currentLevelIndex.y)}

	levelHalfSize := Multiply_Float_Vector2(0.5, &entityComponentsRef.tilemap.worldGridSize)
	currentLevelPos := Vector2{playerInLevel.x * entityComponentsRef.tilemap.worldGridSize.x, playerInLevel.y * entityComponentsRef.tilemap.worldGridSize.y}
	currentLevelCentre := Add_Vector2(&currentLevelPos, &levelHalfSize)

	CameraFollowTarget(currentLevelCentre, entityComponentsRef)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	// TODO : Sort drawings based on Layers order and sort order index and then draw all at once.

	drawImgOptions := ebiten.DrawImageOptions{}

	curSceneRef := g.world.scenes[g.world.currentSceneIndex]
	entityComponentsRef := curSceneRef.entityComponentsForScene

	playerPosRef := &entityComponentsRef.positions[entityComponentsRef.playerEntityID]
	curRoomIndex := entityComponentsRef.tilemap.GetLevelIndexOfPosition(playerPosRef)
	curRoom, curRoomHasSomeData := curSceneRef.roomsData[curRoomIndex]

	playerPos := entityComponentsRef.positions[entityComponentsRef.playerEntityID]
	playerInLevel := Vector2Int{int(playerPos.x) / int(entityComponentsRef.tilemap.worldGridSize.x), int(playerPos.y) / int(entityComponentsRef.tilemap.worldGridSize.y)}

	DrawTileMapLevel(entityComponentsRef.tilemap.levels[playerInLevel], entityComponentsRef, screen, &drawImgOptions)

	if curRoomHasSomeData {
		for _, itemEntityID := range curRoom.itemEntityIDs {
			DrawEntityID(entityComponentsRef, screen, &drawImgOptions, itemEntityID)
		}

		for _, enemyEntityID := range curRoom.enemyEntityIDs {
			DrawEntityID(entityComponentsRef, screen, &drawImgOptions, enemyEntityID)
		}

		for _, obstacleEntityID := range curRoom.obstacleEntityIDs {
			DrawEntityID(entityComponentsRef, screen, &drawImgOptions, obstacleEntityID)
		}
	}

	DrawEntityID(entityComponentsRef, screen, &drawImgOptions, entityComponentsRef.playerEntityID)

	playerCollisionShapeRef := entityComponentsRef.collisionShapes[entityComponentsRef.playerEntityID]
	playerCollisionShapePoints := playerCollisionShapeRef.GetCollisionPoints()
	topLeft := Add_Vector2(&(*playerCollisionShapePoints)[0], &playerPos)
	bottomRight := Add_Vector2(&(*playerCollisionShapePoints)[2], &playerPos)
	size := Vector2{bottomRight.x - topLeft.x, bottomRight.y - topLeft.y}

	vector.StrokeRect(screen, float32(topLeft.x), float32(topLeft.y), float32(size.x), float32(size.y), 1.0, color.Black, false)

	if curRoomHasSomeData {

		for _, enemyEntityID := range curRoom.enemyEntityIDs {
			skelePosRef := &entityComponentsRef.positions[enemyEntityID]

			skeletonCollisionShapeRef := entityComponentsRef.collisionShapes[enemyEntityID]
			vector.StrokeRect(screen, float32(skelePosRef.x), float32(skelePosRef.y), float32(skeletonCollisionShapeRef.(*BoxCollider).size.x), float32(skeletonCollisionShapeRef.(*BoxCollider).size.y), 1.0, color.Black, false)
		}

		for _, obstacleIndex := range curRoom.obstacleEntityIDs {

			obstaclePosRef := &entityComponentsRef.positions[obstacleIndex]
			obstacleCollisionShapeRef := entityComponentsRef.collisionShapes[obstacleIndex]

			if obstacleCollisionShapeRef.CollisionShapeType() == Box {
				vector.StrokeRect(screen, float32(obstaclePosRef.x), float32(obstaclePosRef.y), float32(obstacleCollisionShapeRef.(*BoxCollider).size.x), float32(obstacleCollisionShapeRef.(*BoxCollider).size.y), 1.0, color.Black, false)
			} else if obstacleCollisionShapeRef.CollisionShapeType() == Circle {
				offsetCentre := obstacleCollisionShapeRef.GetOffsetOrigin(obstaclePosRef)
				vector.StrokeCircle(screen, float32(offsetCentre.x), float32(offsetCentre.y), float32(obstacleCollisionShapeRef.GetBoundingBoxDims().x*0.5), 1.0, color.Black, false)
			}
		}

		for _, itemEntityID := range curRoom.itemEntityIDs {

			itemPosRef := &entityComponentsRef.positions[itemEntityID]
			itemCollisionShapeRef := entityComponentsRef.collisionShapes[itemEntityID]
			itemColliderOrigin := itemCollisionShapeRef.GetOffsetOrigin(itemPosRef)

			vector.StrokeCircle(screen, float32(itemColliderOrigin.x), float32(itemColliderOrigin.y), float32(itemCollisionShapeRef.(*CircleCollider).radius), 1.0, color.Black, false)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	entityComponentsRef := g.world.scenes[g.world.currentSceneIndex].entityComponentsForScene
	return int(entityComponentsRef.cameraData.screenSize.x), int(entityComponentsRef.cameraData.screenSize.y)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ninja!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := Game{
		world: CreateAndPopulateWorldScenesAndEntitiesAndComponentsFromGameData("Assets/AssetsData.json"),
	}

	entityComponentsRef := game.world.scenes[game.world.currentSceneIndex].entityComponentsForScene
	entityComponentsRef.cameraData.screenSize = Vector2{320, 240}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
