package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func (sprite Sprite) DrawSprite(screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions, drawPositionRef *Vector2) {
	drawImgOptionsRef.GeoM.Translate(drawPositionRef.x, drawPositionRef.y)
	screenRef.DrawImage(sprite.image.SubImage(image.Rect(sprite.renderRectStart.x, sprite.renderRectStart.y, sprite.renderRectEnd.x, sprite.renderRectEnd.y)).(*ebiten.Image), drawImgOptionsRef)
	drawImgOptionsRef.GeoM.Reset()
}

func DrawEntityID(entityComponentsRef *EntityComponents, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions, entityID int) {
	drawPositionWithCameraOffser := Add_Vector2(&entityComponentsRef.positions[entityID], &entityComponentsRef.cameraData.targetFollowOffset)
	entityComponentsRef.sprites[entityID].DrawSprite(screenRef, drawImgOptionsRef, &drawPositionWithCameraOffser)
}

func DrawTileMapLevel(levelToDrawRef *Level, entityComponentsRef *EntityComponents, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions) {
	for _, layers := range levelToDrawRef.layers {
		for _, tile := range layers.tiles {
			drawTilePosition := Add_Vector2(&tile.position, &entityComponentsRef.cameraData.targetFollowOffset)
			drawImgOptionsRef.GeoM.Translate(drawTilePosition.x, drawTilePosition.y)
			screenRef.DrawImage(entityComponentsRef.tilemap.tileSet.tileSourceImages[tile.TileSetTilesID], drawImgOptionsRef)
			drawImgOptionsRef.GeoM.Reset()
		}
	}
}
