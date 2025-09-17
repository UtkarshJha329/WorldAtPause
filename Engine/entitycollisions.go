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
	CollidedWithObstacleCollisionShapeRef CollisionShape
}

func (curScene *Scene) DoesEntityCollideWithObstacles(roomIndex Vector2Int, const_entityPositionRef *Vector2, const_entityCollisionShapeRef CollisionShape) EntityObstacleCollisionData {

	entityComponents := curScene.EntityComponentsForScene
	curRoom, curRoomHasSomeData := curScene.RoomsData[roomIndex]

	curCollisionData := EntityObstacleCollisionData{
		CollidedWithObstacle: false,
	}

	if !curRoomHasSomeData {
		return curCollisionData
	}

	for _, obstacleEntityID := range curRoom.ObstacleEntityIDs {

		obstaclePos := entityComponents.Positions[obstacleEntityID]
		obstacleCollisionShapeRef := entityComponents.CollisionShapes[obstacleEntityID]
		curCollisionData.CollisionPenetrationNormal, curCollisionData.CollisionPenetrationAmount, curCollisionData.CollidedWithObstacle = CollisionShapeOverlapsWithCollisionShape(const_entityPositionRef, const_entityCollisionShapeRef, &obstaclePos, obstacleCollisionShapeRef)
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

func (curScene *Scene) CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex Vector2Int, entityID int, inputDirection Vector2, totalMoveAmount Vector2, entityMoveAmountPerFrame float64) {

	entityComponentsRef := curScene.EntityComponentsForScene

	entityPos := entityComponentsRef.Positions[entityID]
	entityPosRef := &entityComponentsRef.Positions[entityID]

	entityCollisionShapeRef := entityComponentsRef.CollisionShapes[entityID]

	playerToMoveXPos := entityPos.X + totalMoveAmount.X
	moveDirectionCheckOffsetX := 0.0
	if inputDirection.X > 0 {
		moveDirectionCheckOffsetX = entityCollisionShapeRef.GetBoundingBoxDims().X
	}

	collidingOnXWithTilemap := entityComponentsRef.Tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{playerToMoveXPos + moveDirectionCheckOffsetX, entityPos.Y})
	collidedOnXWithObstacle := false

	if !collidingOnXWithTilemap {
		obstacleCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, &Vector2{playerToMoveXPos, entityPos.Y}, entityCollisionShapeRef)
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

					moveAmountOnTangentForThisCircleObstacle := MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirectionX, circleColliderOriginPosition, entityPos, inputDirection, entityCollisionShapeRef, entityMoveAmountPerFrame)

					finalMoveAmount = Normalise_Vector2(&moveAmountOnTangentForThisCircleObstacle)
					finalMoveAmount = Multiply_Float_Vector2(entityMoveAmountPerFrame, &finalMoveAmount)
				}
			}

			entityPosRef.X += finalMoveAmount.X
			entityPosRef.Y += finalMoveAmount.Y
			entityPos = entityComponentsRef.Positions[entityID]
		}
	}

	playerToMoveYPos := entityPos.Y + totalMoveAmount.Y
	moveDirectionCheckOffsetY := 0.0
	if inputDirection.Y > 0 {
		moveDirectionCheckOffsetY = entityCollisionShapeRef.GetBoundingBoxDims().Y
	}

	collidingOnYWithTilemap := entityComponentsRef.Tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{entityPos.X, playerToMoveYPos + moveDirectionCheckOffsetY})
	collidedOnYWithObstacle := false

	if !collidingOnYWithTilemap {
		obstacleCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, &Vector2{entityPos.X, playerToMoveYPos}, entityCollisionShapeRef)
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

					moveAmountOnTangentForThisCircleObstacle := MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirectionY, circleColliderOriginPosition, entityPos, inputDirection, entityCollisionShapeRef, entityMoveAmountPerFrame)

					finalMoveAmount = Normalise_Vector2(&moveAmountOnTangentForThisCircleObstacle)
					finalMoveAmount = Multiply_Float_Vector2(entityMoveAmountPerFrame, &finalMoveAmount)
				}
			}

			entityPosRef.X += finalMoveAmount.X
			entityPosRef.Y += finalMoveAmount.Y
		}
	}

	if !collidedOnXWithObstacle && !collidedOnYWithObstacle && !collidingOnXWithTilemap {
		entityPosRef.X += totalMoveAmount.X
	}
	if !collidedOnXWithObstacle && !collidedOnYWithObstacle && !collidingOnYWithTilemap {
		entityPosRef.Y += totalMoveAmount.Y
	}
}
