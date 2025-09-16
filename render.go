package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

func DrawTileMapLevel(roomToDrawRef *TilemapRoom, entityComponentsRef *EntityComponents, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions) {
	for _, layers := range roomToDrawRef.layers {
		for _, tile := range layers.tiles {
			drawTilePosition := Add_Vector2(&tile.position, &entityComponentsRef.cameraData.targetFollowOffset)
			drawImgOptionsRef.GeoM.Translate(drawTilePosition.x, drawTilePosition.y)
			screenRef.DrawImage(entityComponentsRef.tilemap.tileSet.tileSourceImages[tile.TileSetTilesID], drawImgOptionsRef)
			drawImgOptionsRef.GeoM.Reset()
		}
	}
}

func DrawActiveRoomInScene(curSceneRef *Scene, screenRef *ebiten.Image) {

	drawImgOptions := ebiten.DrawImageOptions{}

	entityComponentsRef := curSceneRef.entityComponentsForScene

	playerPosRef := &entityComponentsRef.positions[entityComponentsRef.playerEntityID]
	curRoomIndex := entityComponentsRef.tilemap.GetRoomIndexOfPosition(playerPosRef)
	curRoom, curRoomHasSomeData := curSceneRef.roomsData[curRoomIndex]

	DrawTileMapLevel(entityComponentsRef.tilemap.tilemapRooms[curRoomIndex], entityComponentsRef, screenRef, &drawImgOptions)

	if curRoomHasSomeData {

		for _, itemEntityID := range curRoom.itemEntityIDs {
			DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, itemEntityID)
		}

		for _, enemyEntityID := range curRoom.enemyEntityIDs {
			DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, enemyEntityID)
		}

		for _, obstacleEntityID := range curRoom.obstacleEntityIDs {
			DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, obstacleEntityID)
		}
	}

	DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, entityComponentsRef.playerEntityID)
}

func DrawActiveRoomInSceneColliders(curSceneRef *Scene, screenRef *ebiten.Image) {

	entityComponentsRef := curSceneRef.entityComponentsForScene
	playerPos := entityComponentsRef.positions[entityComponentsRef.playerEntityID]

	playerCollisionShapeRef := entityComponentsRef.collisionShapes[entityComponentsRef.playerEntityID]
	topLeft := playerCollisionShapeRef.GetOffsetOrigin(&playerPos)
	size := playerCollisionShapeRef.GetBoundingBoxDims()
	vector.StrokeRect(screenRef, float32(topLeft.x), float32(topLeft.y), float32(size.x), float32(size.y), 1.0, color.Black, false)

	curRoomIndex := entityComponentsRef.tilemap.GetRoomIndexOfPosition(&playerPos)
	curRoom, curRoomHasSomeData := curSceneRef.roomsData[curRoomIndex]

	if curRoomHasSomeData {

		for _, enemyEntityID := range curRoom.enemyEntityIDs {
			skelePosRef := &entityComponentsRef.positions[enemyEntityID]

			skeletonCollisionShapeRef := entityComponentsRef.collisionShapes[enemyEntityID]
			vector.StrokeRect(screenRef, float32(skelePosRef.x), float32(skelePosRef.y), float32(skeletonCollisionShapeRef.(*BoxCollider).size.x), float32(skeletonCollisionShapeRef.(*BoxCollider).size.y), 1.0, color.Black, false)
		}

		for _, obstacleIndex := range curRoom.obstacleEntityIDs {

			obstaclePosRef := &entityComponentsRef.positions[obstacleIndex]
			obstacleCollisionShapeRef := entityComponentsRef.collisionShapes[obstacleIndex]

			if obstacleCollisionShapeRef.CollisionShapeType() == Box {
				vector.StrokeRect(screenRef, float32(obstaclePosRef.x), float32(obstaclePosRef.y), float32(obstacleCollisionShapeRef.(*BoxCollider).size.x), float32(obstacleCollisionShapeRef.(*BoxCollider).size.y), 1.0, color.Black, false)
			} else if obstacleCollisionShapeRef.CollisionShapeType() == Circle {
				offsetCentre := obstacleCollisionShapeRef.GetOffsetOrigin(obstaclePosRef)
				vector.StrokeCircle(screenRef, float32(offsetCentre.x), float32(offsetCentre.y), float32(obstacleCollisionShapeRef.GetBoundingBoxDims().x*0.5), 1.0, color.Black, false)
			}
		}

		for _, itemEntityID := range curRoom.itemEntityIDs {

			itemPosRef := &entityComponentsRef.positions[itemEntityID]
			itemCollisionShapeRef := entityComponentsRef.collisionShapes[itemEntityID]
			itemColliderOrigin := itemCollisionShapeRef.GetOffsetOrigin(itemPosRef)

			vector.StrokeCircle(screenRef, float32(itemColliderOrigin.x), float32(itemColliderOrigin.y), float32(itemCollisionShapeRef.(*CircleCollider).radius), 1.0, color.Black, false)
		}
	}

}
