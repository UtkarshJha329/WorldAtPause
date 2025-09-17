package Engine

import (
	"math"
)

type EntityObstacleCollisionData struct {
	CollidedWithObstacle                  bool
	CollidedWithObstacleEntityID          int
	PositionOfCollidedWithObstacle        Vector2
	ColliderOffsetOriginPosition          Vector2
	CollisionPenetrationAmount            float64
	CollisionPenetrationNormal            Vector2
	CollisionNormal                       Vector2
	CollidedWithObstacleCollisionShapeRef CollisionShape
}

func (curScene *Scene) DoesEntityCollideWithObstacles(roomIndex Vector2Int, entityID int, const_entityPositionRef *Vector2, const_entityCollisionShapeRef CollisionShape) EntityObstacleCollisionData {

	entityComponents := curScene.EntityComponentsForScene
	curRoom, curRoomHasSomeData := curScene.RoomsData[roomIndex]

	curCollisionData := EntityObstacleCollisionData{
		CollidedWithObstacle: false,
	}

	if !curRoomHasSomeData {
		return curCollisionData
	}

	for _, obstacleEntityID := range curRoom.ObstacleEntityIDs {

		if obstacleEntityID == entityID {
			continue
		}

		obstaclePos := entityComponents.Positions[obstacleEntityID]
		obstacleCollisionShapeRef := entityComponents.CollisionShapes[obstacleEntityID]
		curCollisionData.CollisionNormal, curCollisionData.CollisionPenetrationNormal, curCollisionData.CollisionPenetrationAmount, curCollisionData.CollidedWithObstacle = CollisionShapeOverlapsWithCollisionShape(const_entityPositionRef, const_entityCollisionShapeRef, &obstaclePos, obstacleCollisionShapeRef)
		if curCollisionData.CollidedWithObstacle {

			curCollisionData.CollidedWithObstacleEntityID = obstacleEntityID
			curCollisionData.PositionOfCollidedWithObstacle = entityComponents.Positions[obstacleEntityID]
			curCollisionData.CollidedWithObstacleCollisionShapeRef = obstacleCollisionShapeRef
			curCollisionData.ColliderOffsetOriginPosition = *obstacleCollisionShapeRef.GetOffsetOrigin(&obstaclePos)

			return curCollisionData
		}
	}

	return curCollisionData
}

func MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirectionX Vector2, circleColliderOriginPosition Vector2, entityPos Vector2, inputDirection Vector2, entityCollisionShapeRef CollisionShape, entityMoveAmountPerFrame float64) Vector2 {

	normal := circleObstaclePenetrationEscapeDirectionX
	entityBBDims := entityCollisionShapeRef.GetBoundingBoxDims()
	entityColliderCentrePos := Vector2{entityPos.X + entityBBDims.X*0.5, entityPos.Y + entityBBDims.Y*0.5}
	relPos := Subtract_Vector2(&entityColliderCentrePos, &circleColliderOriginPosition)
	relPos = Vector2{math.Copysign(1.0, relPos.X), math.Copysign(1.0, relPos.Y)}

	absNormal := Vector2{math.Abs(normal.X), math.Abs(normal.Y)}

	tangent := normal
	if inputDirection.X != 0 {
		tangent = Vector2{math.Copysign(1.0, inputDirection.X) * absNormal.Y, relPos.Y * absNormal.X}
	}
	if inputDirection.Y != 0 {
		tangent = Vector2{relPos.X * absNormal.Y, math.Copysign(1.0, inputDirection.Y) * absNormal.X}
	}

	tangent = Normalise_Vector2(&tangent)
	return Multiply_Float_Vector2(entityMoveAmountPerFrame, &tangent)
}

type ShapeCollisionWithTilemapResult struct {
	Collided                  bool
	CollidedTileIndex         Vector2Int
	CollidedTilePos           Vector2
	CollisionTileNormal       Vector2
	PenetrationAmount         float64
	CollisionSeparationNormal Vector2
}

func (curScene *Scene) BoxColliderCollidesWithTilemapCollisionLayerInRoom(boxCollider CollisionShape, colliderEntityWorldPos Vector2, roomIndex Vector2Int) ShapeCollisionWithTilemapResult {

	curTilemap := curScene.EntityComponentsForScene.Tilemap

	collisionResult := ShapeCollisionWithTilemapResult{
		Collided:        false,
		CollidedTilePos: Vector2{0.0, 0.0},
	}

	colliderCollisionPoints := boxCollider.GetCollisionPoints()
	for _, collisionPoint := range *colliderCollisionPoints {
		collisionPointWorldPos := Add_Vector2(&collisionPoint, &colliderEntityWorldPos)

		if curTilemap.PointCollidesWithTilemapCollisionLayerInRoom(&roomIndex, &collisionPointWorldPos) {
			collisionResult.Collided = true

			curTile := curScene.EntityComponentsForScene.Tilemap.GetTileOfPointInRoom(&roomIndex, &collisionPointWorldPos)
			curTilePos := curScene.EntityComponentsForScene.Tilemap.GetTilePosInWorld(&roomIndex, &curTile)
			collisionResult.CollidedTilePos = curTilePos

			collisionResult.CollisionTileNormal, collisionResult.CollisionSeparationNormal, collisionResult.PenetrationAmount, collisionResult.Collided = CollisionShapeOverlapsWithCollisionShape(&colliderEntityWorldPos, boxCollider, &curTilePos, &curTilemap.TileCollisionShape)

			if collisionResult.Collided {
				break
			}
		}
	}

	return collisionResult
}

func (curScene *Scene) CircleColliderCollidesWithTilemapCollisionLayerInRoom(circleCollider CollisionShape, colliderEntityWorldPos Vector2, roomIndex Vector2Int) ShapeCollisionWithTilemapResult {

	collisionResult := ShapeCollisionWithTilemapResult{
		Collided:                  false,
		CollisionSeparationNormal: Vector2{0.0, 0.0},
	}

	circleBoundingBoxCol := BoxCollider{
		Collider: Collider{
			ColliderOriginOffset: Vector2{0.0, 0.0},
		},
		size: *circleCollider.GetBoundingBoxDims(),
	}
	circleBoundingBoxCol.CreateCollisionPoints()

	boundingBoxCollisionResult := curScene.BoxColliderCollidesWithTilemapCollisionLayerInRoom(&circleBoundingBoxCol, colliderEntityWorldPos, roomIndex)
	if boundingBoxCollisionResult.Collided {
		collisionTilePosition := boundingBoxCollisionResult.CollidedTilePos
		tileBoundingBox := curScene.EntityComponentsForScene.Tilemap.TileCollisionShape

		collisionResult.CollisionTileNormal, collisionResult.CollisionSeparationNormal, collisionResult.PenetrationAmount, collisionResult.Collided = CollisionShapeOverlapsWithCollisionShape(&colliderEntityWorldPos, circleCollider, &collisionTilePosition, &tileBoundingBox)
		if collisionResult.Collided {
			circleColliderOrigin := circleCollider.GetOffsetOrigin(&colliderEntityWorldPos)
			normalFromTileCollider := tileBoundingBox.GetNormalFromPoint(&collisionTilePosition, circleColliderOrigin)
			collisionResult.CollisionTileNormal = *normalFromTileCollider
		}
	}

	return collisionResult
}

func (curScene *Scene) ShapeCollidesWithTilemapCollisionLayer(curRoomIndex Vector2Int, entityID int, inputDirection Vector2, totalMoveAmount Vector2, entityMoveAmountPerFrame float64) ShapeCollisionWithTilemapResult {

	entityCompnentsRef := curScene.EntityComponentsForScene

	curEntityPos := entityCompnentsRef.Positions[entityID]
	entityCollisionShapeRef := entityCompnentsRef.CollisionShapes[entityID]

	predictedEntityPos := Add_Vector2(&curEntityPos, &totalMoveAmount)

	if entityCollisionShapeRef.CollisionShapeType() == Box {

		return curScene.BoxColliderCollidesWithTilemapCollisionLayerInRoom(entityCollisionShapeRef, predictedEntityPos, curRoomIndex)

	} else if entityCollisionShapeRef.CollisionShapeType() == Circle {

		return curScene.CircleColliderCollidesWithTilemapCollisionLayerInRoom(entityCollisionShapeRef, predictedEntityPos, curRoomIndex)
	}

	return ShapeCollisionWithTilemapResult{Collided: false}
}

type CollideAndMoveCollisionResultSteps struct {
	TilemapXMoveCollisionResult ShapeCollisionWithTilemapResult
	TilemapYMoveCollisionResult ShapeCollisionWithTilemapResult

	ObstacleXMoveCollisionResult EntityObstacleCollisionData
	ObstacleYMoveCollisionResult EntityObstacleCollisionData
}

type CollideAndMoveCollisionParameters struct {
	CollideWithTiles     bool
	CollideWithObstacles bool
	SlideWhenCollide     bool
}

func (curScene *Scene) CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex Vector2Int, entityID int, inputDirection Vector2, totalMoveAmount Vector2, entityMoveAmountPerFrame float64, collideAndMoveParameters CollideAndMoveCollisionParameters) CollideAndMoveCollisionResultSteps {

	collideAndMoveResultData := CollideAndMoveCollisionResultSteps{
		TilemapXMoveCollisionResult: ShapeCollisionWithTilemapResult{Collided: false},
		TilemapYMoveCollisionResult: ShapeCollisionWithTilemapResult{Collided: false},

		ObstacleXMoveCollisionResult: EntityObstacleCollisionData{CollidedWithObstacle: false},
		ObstacleYMoveCollisionResult: EntityObstacleCollisionData{CollidedWithObstacle: false},
	}

	entityComponentsRef := curScene.EntityComponentsForScene

	entityPos := entityComponentsRef.Positions[entityID]
	entityPosRef := &entityComponentsRef.Positions[entityID]

	entityCollisionShapeRef := entityComponentsRef.CollisionShapes[entityID]

	playerToMoveXPos := entityPos.X + totalMoveAmount.X

	collidingOnXWithTilemap := ShapeCollisionWithTilemapResult{
		Collided: false,
	}
	if collideAndMoveParameters.CollideWithTiles {
		collidingOnXWithTilemap = curScene.ShapeCollidesWithTilemapCollisionLayer(curRoomIndex, entityID, Vector2{inputDirection.X, 0.0}, Vector2{totalMoveAmount.X, 0.0}, entityMoveAmountPerFrame)
		collideAndMoveResultData.TilemapXMoveCollisionResult = collidingOnXWithTilemap
	}
	collidedOnXWithObstacle := false

	if !collidingOnXWithTilemap.Collided && collideAndMoveParameters.CollideWithObstacles {
		obstacleCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, entityID, &Vector2{playerToMoveXPos, entityPos.Y}, entityCollisionShapeRef)
		collideAndMoveResultData.ObstacleXMoveCollisionResult = obstacleCollisionResult
		collidedOnXWithObstacle = obstacleCollisionResult.CollidedWithObstacle
		if collidedOnXWithObstacle {

			finalMoveAmount := Vector2{0.0, 0.0}
			maxMoveAmountX := 0.0

			if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
				if inputDirection.X > 0 {
					maxMoveAmountX = obstacleCollisionResult.PositionOfCollidedWithObstacle.X - (entityPos.X + entityCollisionShapeRef.GetBoundingBoxDims().X)
				} else if inputDirection.X < 0 {
					maxMoveAmountX = entityPos.X - (obstacleCollisionResult.PositionOfCollidedWithObstacle.X + obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().X)
				}
				finalMoveAmount = Vector2{inputDirection.X * maxMoveAmountX, inputDirection.Y * entityMoveAmountPerFrame}
			} else if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

				penetrationDistance := obstacleCollisionResult.CollisionPenetrationAmount
				if penetrationDistance != 0 {

					circleColliderOriginPosition := obstacleCollisionResult.ColliderOffsetOriginPosition

					circleObstaclePenetrationEscapeDirectionX := Multiply_Float_Vector2(1.0, &obstacleCollisionResult.CollisionPenetrationNormal)
					circleObstaclePenetrationEscapeDirectionX = Multiply_Float_Vector2(penetrationDistance, &circleObstaclePenetrationEscapeDirectionX)

					if collideAndMoveParameters.SlideWhenCollide {
						moveAmountOnTangentForThisCircleObstacle := MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirectionX, circleColliderOriginPosition, entityPos, inputDirection, entityCollisionShapeRef, entityMoveAmountPerFrame)

						finalMoveAmount = Normalise_Vector2(&moveAmountOnTangentForThisCircleObstacle)
						finalMoveAmount = Multiply_Float_Vector2(entityMoveAmountPerFrame, &finalMoveAmount)
					} else {
						totalMoveAmountX := Vector2{totalMoveAmount.X, 0.0}
						circleObstaclePenetrationEscapeDirectionX = Multiply_Float_Vector2(-1.0, &circleObstaclePenetrationEscapeDirectionX)
						circleObstaclePenetrationEscapeDirectionX = Normalise_Vector2(&circleObstaclePenetrationEscapeDirectionX)
						finalMoveAmount = Add_Vector2(&circleObstaclePenetrationEscapeDirectionX, &totalMoveAmountX)
					}
				}
			}

			entityPosRef.X += finalMoveAmount.X
			entityPosRef.Y += finalMoveAmount.Y
			entityPos = entityComponentsRef.Positions[entityID]
		}
	}

	playerToMoveYPos := entityPos.Y + totalMoveAmount.Y

	collidingOnYWithTilemap := ShapeCollisionWithTilemapResult{
		Collided: false,
	}
	if collideAndMoveParameters.CollideWithTiles {
		collidingOnYWithTilemap = curScene.ShapeCollidesWithTilemapCollisionLayer(curRoomIndex, entityID, Vector2{0.0, inputDirection.Y}, Vector2{0.0, totalMoveAmount.Y}, entityMoveAmountPerFrame)
		collideAndMoveResultData.TilemapYMoveCollisionResult = collidingOnYWithTilemap
	}
	collidedOnYWithObstacle := false

	if !collidingOnYWithTilemap.Collided {
		obstacleCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, entityID, &Vector2{entityPos.X, playerToMoveYPos}, entityCollisionShapeRef)
		collideAndMoveResultData.ObstacleYMoveCollisionResult = obstacleCollisionResult
		collidedOnYWithObstacle = obstacleCollisionResult.CollidedWithObstacle
		if collidedOnYWithObstacle {

			finalMoveAmount := Vector2{0.0, 0.0}
			maxMoveAmountY := 0.0

			if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
				if inputDirection.Y > 0 {
					maxMoveAmountY = obstacleCollisionResult.PositionOfCollidedWithObstacle.Y - (entityPos.Y + entityCollisionShapeRef.GetBoundingBoxDims().Y)
				} else if inputDirection.Y < 0 {
					maxMoveAmountY = entityPos.Y - (obstacleCollisionResult.PositionOfCollidedWithObstacle.Y + obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().Y)
				}
				finalMoveAmount = Vector2{inputDirection.X * entityMoveAmountPerFrame, inputDirection.Y * maxMoveAmountY}

			} else if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

				penetrationDistance := obstacleCollisionResult.CollisionPenetrationAmount
				if penetrationDistance != 0 {
					circleColliderOriginPosition := obstacleCollisionResult.ColliderOffsetOriginPosition

					circleObstaclePenetrationEscapeDirectionY := Multiply_Float_Vector2(1.0, &obstacleCollisionResult.CollisionPenetrationNormal)
					circleObstaclePenetrationEscapeDirectionY = Multiply_Float_Vector2(penetrationDistance, &circleObstaclePenetrationEscapeDirectionY)

					if collideAndMoveParameters.SlideWhenCollide {
						moveAmountOnTangentForThisCircleObstacle := MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirectionY, circleColliderOriginPosition, entityPos, inputDirection, entityCollisionShapeRef, entityMoveAmountPerFrame)

						finalMoveAmount = Normalise_Vector2(&moveAmountOnTangentForThisCircleObstacle)
						finalMoveAmount = Multiply_Float_Vector2(entityMoveAmountPerFrame, &finalMoveAmount)
					} else {
						totalMoveAmountY := Vector2{0.0, totalMoveAmount.Y}
						circleObstaclePenetrationEscapeDirectionY = Multiply_Float_Vector2(-1.0, &circleObstaclePenetrationEscapeDirectionY)
						circleObstaclePenetrationEscapeDirectionY = Normalise_Vector2(&circleObstaclePenetrationEscapeDirectionY)
						finalMoveAmount = Add_Vector2(&circleObstaclePenetrationEscapeDirectionY, &totalMoveAmountY)
					}
				}
			}

			entityPosRef.X += finalMoveAmount.X
			entityPosRef.Y += finalMoveAmount.Y
		}
	}

	if !collidedOnXWithObstacle && !collidedOnYWithObstacle && !collidingOnXWithTilemap.Collided {
		entityPosRef.X += totalMoveAmount.X
	}
	if !collidedOnXWithObstacle && !collidedOnYWithObstacle && !collidingOnYWithTilemap.Collided {
		entityPosRef.Y += totalMoveAmount.Y
	}

	return collideAndMoveResultData
}
