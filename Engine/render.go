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

func DrawEntityIDIfAlive(entityComponentsRef *EntityComponents, screenRef *ebiten.Image, drawImgOptionsRef *ebiten.DrawImageOptions, entityID int) {
	if entityComponentsRef.IsEntityAlive(entityID) {
		drawPositionWithCameraOffser := Add_Vector2(&entityComponentsRef.Positions[entityID], &entityComponentsRef.CameraData.TargetFollowOffset)
		entityComponentsRef.Sprites[entityID].DrawSprite(screenRef, drawImgOptionsRef, &drawPositionWithCameraOffser)
	}
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

	curTilemapRoomRef, ok := entityComponentsRef.Tilemap.TilemapRooms[curRoomIndex]
	if ok {
		DrawTileMapLevel(curTilemapRoomRef, entityComponentsRef, screenRef, &drawImgOptions)
	}

	if curRoomHasSomeData {

		for _, triggerEntityID := range curRoom.VisibleTriggerEntityIDs {
			DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, triggerEntityID)
		}

		for _, itemEntityID := range curRoom.ItemEntityIDs {
			DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, itemEntityID)
		}

		for _, obstacleEntityID := range curRoom.ObstacleEntityIDs {
			DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, obstacleEntityID)
		}

		for _, enemyEntityID := range curRoom.EnemyEntityIDs {
			DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, enemyEntityID)
		}
	}

	DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, entityComponentsRef.PlayerEntityID)
}

func DrawActiveTilemapInRoomInScene(curSceneRef *Scene, screenRef *ebiten.Image) {

	drawImgOptions := ebiten.DrawImageOptions{}

	entityComponentsRef := curSceneRef.EntityComponentsForScene

	playerPosRef := &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)

	curTilemapRoomRef, ok := entityComponentsRef.Tilemap.TilemapRooms[curRoomIndex]
	if ok {
		DrawTileMapLevel(curTilemapRoomRef, entityComponentsRef, screenRef, &drawImgOptions)
	}
}

func DrawActiveRoomObjectsInScene(curSceneRef *Scene, screenRef *ebiten.Image) {

	drawImgOptions := ebiten.DrawImageOptions{}

	entityComponentsRef := curSceneRef.EntityComponentsForScene

	playerPosRef := &entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]
	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(playerPosRef)
	curRoom, curRoomHasSomeData := curSceneRef.RoomsData[curRoomIndex]

	if curRoomHasSomeData {

		for _, triggerEntityID := range curRoom.VisibleTriggerEntityIDs {
			DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, triggerEntityID)
		}

		for _, itemEntityID := range curRoom.ItemEntityIDs {
			DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, itemEntityID)
		}

		for _, obstacleEntityID := range curRoom.ObstacleEntityIDs {
			DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, obstacleEntityID)
		}

		for _, enemyEntityID := range curRoom.EnemyEntityIDs {
			DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, enemyEntityID)
		}
	}
}

func DrawPlayerInScene(curSceneRef *Scene, screenRef *ebiten.Image) {

	drawImgOptions := ebiten.DrawImageOptions{}

	entityComponentsRef := curSceneRef.EntityComponentsForScene

	DrawEntityIDIfAlive(entityComponentsRef, screenRef, &drawImgOptions, entityComponentsRef.PlayerEntityID)
}

func DrawColliderForEntityIfAlive(entityComponentsRef *EntityComponents, screenRef *ebiten.Image, entityID int) {
	if entityComponentsRef.IsEntityAlive(entityID) {
		entityPos := entityComponentsRef.Positions[entityID]
		entityCollisionShapeRef := entityComponentsRef.CollisionShapes[entityID]

		if entityCollisionShapeRef.CollisionShapeType() == Box {
			entityCollisionPoints := entityCollisionShapeRef.GetCollisionPoints()
			topLeftCollider := (*entityCollisionPoints)[0]
			topLeftCollider = Add_Vector2(&entityPos, &topLeftCollider)
			vector.StrokeRect(screenRef, float32(topLeftCollider.X), float32(topLeftCollider.Y), float32(entityCollisionShapeRef.(*BoxCollider).size.X), float32(entityCollisionShapeRef.(*BoxCollider).size.Y), 1.0, color.White, false)
		} else if entityCollisionShapeRef.CollisionShapeType() == Circle {
			offsetCentre := entityCollisionShapeRef.GetOffsetOrigin(&entityPos)
			vector.StrokeCircle(screenRef, float32(offsetCentre.X), float32(offsetCentre.Y), float32(entityCollisionShapeRef.GetBoundingBoxDims().X*0.5), 1.0, color.White, false)
		}
	}
}

func DrawActiveRoomInSceneColliders(curSceneRef *Scene, screenRef *ebiten.Image) {

	entityComponentsRef := curSceneRef.EntityComponentsForScene
	playerPos := entityComponentsRef.Positions[entityComponentsRef.PlayerEntityID]

	DrawColliderForEntityIfAlive(entityComponentsRef, screenRef, entityComponentsRef.PlayerEntityID)

	curRoomIndex := entityComponentsRef.Tilemap.GetRoomIndexOfPosition(&playerPos)
	curRoom, curRoomHasSomeData := curSceneRef.RoomsData[curRoomIndex]

	if curRoomHasSomeData {

		for _, enemyEntityID := range curRoom.EnemyEntityIDs {
			DrawColliderForEntityIfAlive(entityComponentsRef, screenRef, enemyEntityID)
		}

		for _, obstacleEntityID := range curRoom.ObstacleEntityIDs {
			DrawColliderForEntityIfAlive(entityComponentsRef, screenRef, obstacleEntityID)
		}

		for _, itemEntityID := range curRoom.ItemEntityIDs {
			DrawColliderForEntityIfAlive(entityComponentsRef, screenRef, itemEntityID)
		}
	}

}
