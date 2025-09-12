package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	entityComponentsRef *EntityComponents
}

func (g *Game) Update() error {

	playerPosRef := &g.entityComponentsRef.positions[g.entityComponentsRef.playerEntityID]

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
		totalMoveAmount := Multiply_Float_Vector2(2.0, &normalisedInputDir)

		playerPosRef.x += totalMoveAmount.x
		playerPosRef.y += totalMoveAmount.y
	}

	itemPickupDistance := 16.0
	stoppageDistanceFromPlayer := 32.0
	skeleMoveAmountPerFrame := 2.0
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
	}

	for index, itemEntityID := range g.entityComponentsRef.itemEntityIDs {
		itemPosRef := &g.entityComponentsRef.positions[itemEntityID]

		if DistanceSquare_Vector2(itemPosRef, playerPosRef) < math.Pow(itemPickupDistance, 2) {
			g.entityComponentsRef.itemEntityIDs = append(g.entityComponentsRef.itemEntityIDs[:index], g.entityComponentsRef.itemEntityIDs[index+1:]...)
			fmt.Printf("Picked up item entity ID : %d\n", itemEntityID)
		}
	}

	playerPos := g.entityComponentsRef.positions[g.entityComponentsRef.playerEntityID]
	playerInLevel := Vector2{float64(int(playerPos.x) / int(g.entityComponentsRef.tilemap.worldGridSize.x)), float64(int(playerPos.y) / int(g.entityComponentsRef.tilemap.worldGridSize.y))}
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

	DrawEntityID(g, screen, &drawImgOptions, g.entityComponentsRef.playerEntityID)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.entityComponentsRef.cameraData.screenSize.x), int(g.entityComponentsRef.cameraData.screenSize.y)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ninja!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := Game{
		entityComponentsRef: CreateAndPopulateEntitiesAndComponentsFromGameData(&EmbeddedAssetsFS, "Assets/AssetsData.json"),
	}

	game.entityComponentsRef.cameraData.screenSize = Vector2{320, 240}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
