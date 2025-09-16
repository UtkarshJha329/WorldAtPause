package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	world *World
}

func (g *Game) Update() error {

	curSceneRef := g.world.scenes[g.world.currentSceneIndex]
	entityComponentsRef := curSceneRef.entityComponentsForScene

	playerMoveAmountPerFrame := 2.0

	playerPosRef := &entityComponentsRef.positions[entityComponentsRef.playerEntityID]
	curRoomIndex := entityComponentsRef.tilemap.GetRoomIndexOfPosition(playerPosRef)
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
	currentLevelIndex := entityComponentsRef.tilemap.GetRoomIndexOfPosition(playerPosRef)
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

	DrawActiveRoomInScene(g.world.scenes[g.world.currentSceneIndex], screen)
	DrawActiveRoomInSceneColliders(g.world.scenes[g.world.currentSceneIndex], screen)

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
