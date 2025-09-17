package Engine

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (sprite Sprite) DrawSprite(screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions, drawPositionRef *Vector2) {
	drawImgOptionsRef.GeoM.Translate(drawPositionRef.X, drawPositionRef.Y)
	screenRef.DrawImage(sprite.Image.SubImage(image.Rect(sprite.RenderRectStart.X, sprite.RenderRectStart.Y, sprite.RenderRectEnd.X, sprite.RenderRectEnd.Y)).(*ebiten.Image), drawImgOptionsRef)
	drawImgOptionsRef.GeoM.Reset()
}

func DrawEntityID(entityComponentsRef *EntityComponents, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions, entityID int) {
	drawPositionWithCameraOffser := Add_Vector2(&entityComponentsRef.Positions[entityID], &entityComponentsRef.CameraData.TargetFollowOffset)
	entityComponentsRef.Sprites[entityID].DrawSprite(screenRef, drawImgOptionsRef, &drawPositionWithCameraOffser)
}

func DrawTileMapLevel(roomToDrawRef *TilemapRoom, entityComponentsRef *EntityComponents, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions) {
	for _, layers := range roomToDrawRef.layers {
		for _, tile := range layers.tiles {
			drawTilePosition := Add_Vector2(&tile.position, &entityComponentsRef.CameraData.TargetFollowOffset)
			drawImgOptionsRef.GeoM.Translate(drawTilePosition.X, drawTilePosition.Y)
			screenRef.DrawImage(entityComponentsRef.Tilemap.TileSet.tileSourceImages[tile.TileSetTilesID], drawImgOptionsRef)
			drawImgOptionsRef.GeoM.Reset()
		}
	}
}

func DrawActiveRoomInScene(curSceneRef *Scene, screenRef *ebiten.Image) {

	drawImgOptions := ebiten.DrawImageOptions{}

	entityComponentsRef := curSceneRef.EntityComponentsForScene

	playerPosRef := &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)
	curRoom, curRoomHasSomeData := curSceneRef.RoomsData[curRoomIndex]

	DrawTileMapLevel(entityComponentsRef.Tilemap.TilemapRooms[curRoomIndex], entityComponentsRef, screenRef, &drawImgOptions)

	if curRoomHasSomeData {

		for _, itemEntityID := range curRoom.ItemEntityIDs {
			DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, itemEntityID)
		}

		for _, enemyEntityID := range curRoom.EnemyEntityIDs {
			DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, enemyEntityID)
		}

		for _, obstacleEntityID := range curRoom.ObstacleEntityIDs {
			DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, obstacleEntityID)
		}

		for _, physicsEntityID := range curRoom.PhysicsEntityIDs {
			DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, physicsEntityID)
		}
	}

	DrawEntityID(entityComponentsRef, screenRef, &drawImgOptions, entityComponentsRef.PlayerEntityID)
}

func DrawActiveRoomInSceneColliders(curSceneRef *Scene, screenRef *ebiten.Image) {

	entityComponentsRef := curSceneRef.EntityComponentsForScene
	playerPos := entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]

	playerCollisionShapeRef := entityComponentsRef.CollisionShapes[entityComponentsRef.PlayerEntityID]
	topLeft := playerCollisionShapeRef.GetOffsetOrigin(&playerPos)
	size := playerCollisionShapeRef.GetBoundingBoxDims()
	vector.StrokeRect(screenRef, float32(topLeft.X), float32(topLeft.Y), float32(size.X), float32(size.Y), 1.0, color.Black, false)

	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(&playerPos)
	curRoom, curRoomHasSomeData := curSceneRef.RoomsData[curRoomIndex]

	if curRoomHasSomeData {

		for _, enemyEntityID := range curRoom.EnemyEntityIDs {
			skelePosRef := &entityComponentsRef.Positions[enemyEntityID]

			skeletonCollisionShapeRef := entityComponentsRef.CollisionShapes[enemyEntityID]
			vector.StrokeRect(screenRef, float32(skelePosRef.X), float32(skelePosRef.Y), float32(skeletonCollisionShapeRef.(*BoxCollider).size.X), float32(skeletonCollisionShapeRef.(*BoxCollider).size.Y), 1.0, color.Black, false)
		}

		for _, obstacleIndex := range curRoom.ObstacleEntityIDs {

			obstaclePosRef := &entityComponentsRef.Positions[obstacleIndex]
			obstacleCollisionShapeRef := entityComponentsRef.CollisionShapes[obstacleIndex]

			if obstacleCollisionShapeRef.CollisionShapeType() == Box {
				vector.StrokeRect(screenRef, float32(obstaclePosRef.X), float32(obstaclePosRef.Y), float32(obstacleCollisionShapeRef.(*BoxCollider).size.X), float32(obstacleCollisionShapeRef.(*BoxCollider).size.Y), 1.0, color.Black, false)
			} else if obstacleCollisionShapeRef.CollisionShapeType() == Circle {
				offsetCentre := obstacleCollisionShapeRef.GetOffsetOrigin(obstaclePosRef)
				vector.StrokeCircle(screenRef, float32(offsetCentre.X), float32(offsetCentre.Y), float32(obstacleCollisionShapeRef.GetBoundingBoxDims().X*0.5), 1.0, color.Black, false)
			}
		}

		for _, itemEntityID := range curRoom.ItemEntityIDs {

			itemPosRef := &entityComponentsRef.Positions[itemEntityID]
			itemCollisionShapeRef := entityComponentsRef.CollisionShapes[itemEntityID]
			itemColliderOrigin := itemCollisionShapeRef.GetOffsetOrigin(itemPosRef)

			vector.StrokeCircle(screenRef, float32(itemColliderOrigin.X), float32(itemColliderOrigin.Y), float32(itemCollisionShapeRef.(*CircleCollider).radius), 1.0, color.Black, false)
		}

		for _, physicsEntityID := range curRoom.PhysicsEntityIDs {

			physicsEntityPosRef := &entityComponentsRef.Positions[physicsEntityID]
			physicsEntityCollisionShapeRef := entityComponentsRef.CollisionShapes[physicsEntityID]
			physicsEntityColliderOrigin := physicsEntityCollisionShapeRef.GetOffsetOrigin(physicsEntityPosRef)

			vector.StrokeCircle(screenRef, float32(physicsEntityColliderOrigin.X), float32(physicsEntityColliderOrigin.Y), float32(physicsEntityCollisionShapeRef.(*CircleCollider).radius), 1.0, color.Black, false)
		}
	}

}
