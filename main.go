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
	tilemap             Tilemap
}

func (g *Game) Update() error {

	playerPosRef := &g.entityComponentsRef.positions[PlayerEntityID]

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		playerPosRef.x += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		playerPosRef.x -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		playerPosRef.y += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		playerPosRef.y -= 2
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

	CameraFollowTarget(g.entityComponentsRef.positions[g.entityComponentsRef.playerEntityID], g.entityComponentsRef)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{120, 180, 255, 255})

	// TODO : Sort drawings based on Layers order and sort order index and then draw all at once.

	drawImgOptions := ebiten.DrawImageOptions{}

	for _, levels := range g.tilemap.levels {
		for _, layers := range levels.layers {
			for _, tile := range layers.tiles {
				drawTilePosition := Add_Vector2(&tile.position, &g.entityComponentsRef.cameraData.targetFollowOffset)
				drawImgOptions.GeoM.Translate(drawTilePosition.x, drawTilePosition.y)
				screen.DrawImage(g.tilemap.tileSet.tileSourceImages[tile.TileSetTilesID], &drawImgOptions)
				drawImgOptions.GeoM.Reset()
			}
		}
	}

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

const PlayerEntityID = 0

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ninja!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := Game{
		entityComponentsRef: CreateAndPopulateEntitiesAndComponentsFromGameData("Assets/AssetsData.json"),
		tilemap:             Tilemap{},
	}

	game.entityComponentsRef.cameraData.screenSize = Vector2{320, 240}

	err := NewTilemap("Assets/Maps/TilesetFloor.png", "Assets/Maps/WorldMap.json", 16, &game.tilemap)
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
