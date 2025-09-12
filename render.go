package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func (sprite Sprite) DrawSprite(gRef *Game, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions, drawPositionRef *Vector2) {
	drawImgOptionsRef.GeoM.Translate(drawPositionRef.x, drawPositionRef.y)
	screenRef.DrawImage(sprite.image.SubImage(image.Rect(sprite.renderRectStart.x, sprite.renderRectStart.y, sprite.renderRectEnd.x, sprite.renderRectEnd.y)).(*ebiten.Image), drawImgOptionsRef)
	drawImgOptionsRef.GeoM.Reset()
}

func DrawEntityID(gRef *Game, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions, entityID int) {
	drawPositionWithCameraOffser := Add_Vector2(&gRef.entityComponentsRef.positions[entityID], &gRef.entityComponentsRef.cameraData.targetFollowOffset)
	gRef.entityComponentsRef.sprites[entityID].DrawSprite(gRef, screenRef, drawImgOptionsRef, &drawPositionWithCameraOffser)
}

func DrawTileMapLevel(levelToDrawRef *Level, gRef *Game, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions) {
	for _, layers := range levelToDrawRef.layers {
		for _, tile := range layers.tiles {
			drawTilePosition := Add_Vector2(&tile.position, &gRef.entityComponentsRef.cameraData.targetFollowOffset)
			drawImgOptionsRef.GeoM.Translate(drawTilePosition.x, drawTilePosition.y)
			screenRef.DrawImage(gRef.entityComponentsRef.tilemap.tileSet.tileSourceImages[tile.TileSetTilesID], drawImgOptionsRef)
			drawImgOptionsRef.GeoM.Reset()
		}
	}
}
