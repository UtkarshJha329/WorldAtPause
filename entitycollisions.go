package main

import (
	"math"
)

type EntityObstacleCollisionData struct {
	collidedWithObstacle                         bool
	collidedWithObstacleEntityID                 int
	positionOfCollidedWithObstacle               Vector2
	colliderOffsetOriginPosition                 Vector2
	boxEntityCircleObstacleCirclePenetrationData CircleBoxOverlapCirclePenetrationData
	collidedWithObstacleCollisionShapeRef        CollisionShape
}

func (curScene *Scene) DoesEntityCollideWithObstacles(roomIndex Vector2Int, const_entityPositionRef *Vector2, const_entityCollisionShapeRef CollisionShape) EntityObstacleCollisionData {

	entityComponents := curScene.entityComponentsForScene
	curRoom, curRoomHasSomeData := curScene.roomsData[roomIndex]

	curCollisionData := EntityObstacleCollisionData{
		collidedWithObstacle: false,
	}

	if !curRoomHasSomeData {
		return curCollisionData
	}

	for _, obstacleEntityID := range curRoom.obstacleEntityIDs {

		obstaclePos := entityComponents.positions[obstacleEntityID]
		obstacleCollisionShapeRef := entityComponents.collisionShapes[obstacleEntityID]
		if CollisionShapeOverlapsWithCollisionShape(const_entityPositionRef, const_entityCollisionShapeRef, &obstaclePos, obstacleCollisionShapeRef) {

			curCollisionData.collidedWithObstacle = true
			curCollisionData.collidedWithObstacleEntityID = obstacleEntityID
			curCollisionData.positionOfCollidedWithObstacle = entityComponents.positions[obstacleEntityID]
			curCollisionData.collidedWithObstacleCollisionShapeRef = obstacleCollisionShapeRef

			if obstacleCollisionShapeRef.CollisionShapeType() == Circle {
				curCollisionData.boxEntityCircleObstacleCirclePenetrationData = GetCircleBoxOverlapPenetrationData(&obstaclePos, obstacleCollisionShapeRef.(*CircleCollider), const_entityPositionRef, const_entityCollisionShapeRef.(*BoxCollider))
				curCollisionData.colliderOffsetOriginPosition = *obstacleCollisionShapeRef.GetOffsetOrigin(&obstaclePos)
			}

			return curCollisionData
		}
	}

	return curCollisionData
}

func MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirectionX Vector2, circleColliderOriginPosition Vector2, entityPos Vector2, inputDirection Vector2, entityCollisionShapeRef CollisionShape, entityMoveAmountPerFrame float64) Vector2 {

	normal := circleObstaclePenetrationEscapeDirectionX
	entityBBDims := entityCollisionShapeRef.GetBoundingBoxDims()
	entityColliderCentrePos := Vector2{entityPos.x + entityBBDims.x*0.5, entityPos.y + entityBBDims.y*0.5}
	relPos := Subtract_Vector2(&entityColliderCentrePos, &circleColliderOriginPosition)
	relPos = Vector2{math.Copysign(1.0, relPos.x), math.Copysign(1.0, relPos.y)}

	absNormal := Vector2{math.Abs(normal.x), math.Abs(normal.y)}

	tangent := normal
	if inputDirection.x != 0 {
		tangent = Vector2{math.Copysign(1.0, inputDirection.x) * absNormal.y, relPos.y * absNormal.x}
	}
	if inputDirection.y != 0 {
		tangent = Vector2{relPos.x * absNormal.y, math.Copysign(1.0, inputDirection.y) * absNormal.x}
	}

	tangent = Normalise_Vector2(&tangent)
	return Multiply_Float_Vector2(entityMoveAmountPerFrame, &tangent)
}

func (curScene *Scene) CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex Vector2Int, entityID int, inputDirection Vector2, totalMoveAmount Vector2, entityMoveAmountPerFrame float64) {

	entityComponentsRef := curScene.entityComponentsForScene

	entityPos := entityComponentsRef.positions[entityID]
	entityPosRef := &entityComponentsRef.positions[entityID]

	entityCollisionShapeRef := entityComponentsRef.collisionShapes[entityID]

	playerToMoveXPos := entityPos.x + totalMoveAmount.x
	moveDirectionCheckOffsetX := 0.0
	if inputDirection.x > 0 {
		moveDirectionCheckOffsetX = entityCollisionShapeRef.GetBoundingBoxDims().x
	}

	collidingOnXWithTilemap := entityComponentsRef.tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{playerToMoveXPos + moveDirectionCheckOffsetX, entityPos.y})
	collidedOnXWithObstacle := false

	if !collidingOnXWithTilemap {
		obstacleCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, &Vector2{playerToMoveXPos, entityPos.y}, entityCollisionShapeRef)
		collidedOnXWithObstacle = obstacleCollisionResult.collidedWithObstacle
		if collidedOnXWithObstacle {

			finalMoveAmount := Vector2{0.0, 0.0}
			maxMoveAmountX := 0.0

			if obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
				if inputDirection.x > 0 {
					maxMoveAmountX = obstacleCollisionResult.positionOfCollidedWithObstacle.x - (entityPos.x + entityCollisionShapeRef.GetBoundingBoxDims().x)
				} else if inputDirection.x < 0 {
					maxMoveAmountX = entityPos.x - (obstacleCollisionResult.positionOfCollidedWithObstacle.x + obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().x)
				}
				finalMoveAmount = Vector2{inputDirection.x * maxMoveAmountX, inputDirection.y * entityMoveAmountPerFrame}
			} else if obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

				penetrationDistance := obstacleCollisionResult.boxEntityCircleObstacleCirclePenetrationData.penetrationAmount
				if penetrationDistance != 0 {

					circleColliderOriginPosition := obstacleCollisionResult.colliderOffsetOriginPosition

					circleObstaclePenetrationEscapeDirectionX := Multiply_Float_Vector2(1.0, &obstacleCollisionResult.boxEntityCircleObstacleCirclePenetrationData.normal)
					circleObstaclePenetrationEscapeDirectionX = Multiply_Float_Vector2(penetrationDistance, &circleObstaclePenetrationEscapeDirectionX)

					moveAmountOnTangentForThisCircleObstacle := MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirectionX, circleColliderOriginPosition, entityPos, inputDirection, entityCollisionShapeRef, entityMoveAmountPerFrame)

					finalMoveAmount = Normalise_Vector2(&moveAmountOnTangentForThisCircleObstacle)
					finalMoveAmount = Multiply_Float_Vector2(entityMoveAmountPerFrame, &finalMoveAmount)
				}
			}

			entityPosRef.x += finalMoveAmount.x
			entityPosRef.y += finalMoveAmount.y
			entityPos = entityComponentsRef.positions[entityID]
		}
	}

	playerToMoveYPos := entityPos.y + totalMoveAmount.y
	moveDirectionCheckOffsetY := 0.0
	if inputDirection.y > 0 {
		moveDirectionCheckOffsetY = entityCollisionShapeRef.GetBoundingBoxDims().y
	}

	collidingOnYWithTilemap := entityComponentsRef.tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{entityPos.x, playerToMoveYPos + moveDirectionCheckOffsetY})
	collidedOnYWithObstacle := false

	if !collidingOnYWithTilemap {
		obstacleCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, &Vector2{entityPos.x, playerToMoveYPos}, entityCollisionShapeRef)
		collidedOnYWithObstacle = obstacleCollisionResult.collidedWithObstacle
		if collidedOnYWithObstacle {

			finalMoveAmount := Vector2{0.0, 0.0}
			maxMoveAmountY := 0.0

			if obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
				if inputDirection.y > 0 {
					maxMoveAmountY = obstacleCollisionResult.positionOfCollidedWithObstacle.y - (entityPos.y + entityCollisionShapeRef.GetBoundingBoxDims().y)
				} else if inputDirection.y < 0 {
					maxMoveAmountY = entityPos.y - (obstacleCollisionResult.positionOfCollidedWithObstacle.y + obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().y)
				}
				finalMoveAmount = Vector2{inputDirection.x * entityMoveAmountPerFrame, inputDirection.y * maxMoveAmountY}

			} else if obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

				penetrationDistance := obstacleCollisionResult.boxEntityCircleObstacleCirclePenetrationData.penetrationAmount
				if penetrationDistance != 0 {
					circleColliderOriginPosition := obstacleCollisionResult.colliderOffsetOriginPosition

					circleObstaclePenetrationEscapeDirectionY := Multiply_Float_Vector2(1.0, &obstacleCollisionResult.boxEntityCircleObstacleCirclePenetrationData.normal)
					circleObstaclePenetrationEscapeDirectionY = Multiply_Float_Vector2(penetrationDistance, &circleObstaclePenetrationEscapeDirectionY)

					moveAmountOnTangentForThisCircleObstacle := MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirectionY, circleColliderOriginPosition, entityPos, inputDirection, entityCollisionShapeRef, entityMoveAmountPerFrame)

					finalMoveAmount = Normalise_Vector2(&moveAmountOnTangentForThisCircleObstacle)
					finalMoveAmount = Multiply_Float_Vector2(entityMoveAmountPerFrame, &finalMoveAmount)
				}
			}

			entityPosRef.x += finalMoveAmount.x
			entityPosRef.y += finalMoveAmount.y
		}
	}

	if !collidedOnXWithObstacle && !collidedOnYWithObstacle && !collidingOnXWithTilemap {
		entityPosRef.x += totalMoveAmount.x
	}
	if !collidedOnXWithObstacle && !collidedOnYWithObstacle && !collidingOnYWithTilemap {
		entityPosRef.y += totalMoveAmount.y
	}
}
