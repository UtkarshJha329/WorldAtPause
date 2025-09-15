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

func (entityComponents *EntityComponents) DoesEntityCollideWithObstacles(const_entityPositionRef *Vector2, const_entityCollisionShapeRef CollisionShape) EntityObstacleCollisionData {

	curCollisionData := EntityObstacleCollisionData{
		collidedWithObstacle: false,
	}
	for _, obstacleEntityID := range entityComponents.obstacleEntityIDs {

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

func Sign(x float64) float64 {
	if x < 0.0 {
		return -1.0
	}
	return 1.0
}

func (entityComponentsRef *EntityComponents) MoveAndCollideEntityWithTilemapAndObstacles(entityID int, inputDirection Vector2, totalMoveAmount Vector2, entityMoveAmountPerFrame float64) {

	entityPos := entityComponentsRef.positions[entityID]
	entityPosRef := &entityComponentsRef.positions[entityID]

	entityCollisionShapeRef := entityComponentsRef.collisionShapes[entityID]

	playerToMoveXPos := entityPos.x + totalMoveAmount.x
	moveDirectionCheckOffsetX := 0.0
	if inputDirection.x > 0 {
		moveDirectionCheckOffsetX = entityCollisionShapeRef.GetBoundingBoxDims().x
	}

	collidingOnXWithTilemap := entityComponentsRef.tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{playerToMoveXPos + moveDirectionCheckOffsetX, entityPos.y})
	collidingOnX := collidingOnXWithTilemap

	circleCollision := false
	circlePenetrationObjectEscapeDirection := Vector2{0.0, 0.0}
	circleObstaclePenetrationEscapeDirectionX := Vector2{0.0, 0.0}
	circleObstaclePenetrationEscapeDirectionY := Vector2{0.0, 0.0}
	circlePenetrationEscapeEpsilon := 0.0

	circleColliderOriginPosition := Vector2{0.0, 0.0}

	maxMoveAmtOnX := 0.0
	if !collidingOnXWithTilemap {
		obstacleCollisionResult := entityComponentsRef.DoesEntityCollideWithObstacles(&Vector2{playerToMoveXPos, entityPos.y}, entityCollisionShapeRef)
		if obstacleCollisionResult.collidedWithObstacle {
			collidingOnX = true
			if obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
				if inputDirection.x > 0 {
					maxMoveAmtOnX = obstacleCollisionResult.positionOfCollidedWithObstacle.x - (entityPos.x + entityCollisionShapeRef.GetBoundingBoxDims().x)
				} else if inputDirection.x < 0 {
					maxMoveAmtOnX = entityPos.x - (obstacleCollisionResult.positionOfCollidedWithObstacle.x + obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().x)
				}
			} else if obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

				penetrationDistance := obstacleCollisionResult.boxEntityCircleObstacleCirclePenetrationData.penetrationAmount
				if penetrationDistance != 0 {

					circleCollision = true
					circleColliderOriginPosition = obstacleCollisionResult.colliderOffsetOriginPosition

					penetrationDistance += circlePenetrationEscapeEpsilon
					circleObstaclePenetrationEscapeDirectionX = Multiply_Float_Vector2(1.0, &obstacleCollisionResult.boxEntityCircleObstacleCirclePenetrationData.normal)
					circleObstaclePenetrationEscapeDirectionX = Multiply_Float_Vector2(penetrationDistance, &circleObstaclePenetrationEscapeDirectionX)

				}
			}
		}
	}

	playerToMoveYPos := entityPos.y + totalMoveAmount.y
	moveDirectionCheckOffsetY := 0.0
	if inputDirection.y > 0 {
		moveDirectionCheckOffsetY = entityCollisionShapeRef.GetBoundingBoxDims().y
	}

	collidingOnYWithTilemap := entityComponentsRef.tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{entityPos.x, playerToMoveYPos + moveDirectionCheckOffsetY})
	collidingOnY := collidingOnYWithTilemap

	maxMoveAmtOnY := 0.0
	if !collidingOnYWithTilemap {
		obstacleCollisionResult := entityComponentsRef.DoesEntityCollideWithObstacles(&Vector2{entityPos.x, playerToMoveYPos}, entityCollisionShapeRef)
		if obstacleCollisionResult.collidedWithObstacle {
			collidingOnY = true
			if obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
				if inputDirection.y > 0 {
					maxMoveAmtOnY = obstacleCollisionResult.positionOfCollidedWithObstacle.y - (entityPos.y + entityCollisionShapeRef.GetBoundingBoxDims().y)
				} else if inputDirection.y < 0 {
					maxMoveAmtOnY = entityPos.y - (obstacleCollisionResult.positionOfCollidedWithObstacle.y + obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().y)
				}
			} else if obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

				penetrationDistance := obstacleCollisionResult.boxEntityCircleObstacleCirclePenetrationData.penetrationAmount
				if penetrationDistance != 0 {

					circleCollision = true
					circleColliderOriginPosition = obstacleCollisionResult.colliderOffsetOriginPosition

					penetrationDistance += circlePenetrationEscapeEpsilon
					circleObstaclePenetrationEscapeDirectionY = Multiply_Float_Vector2(1.0, &obstacleCollisionResult.boxEntityCircleObstacleCirclePenetrationData.normal)
					circleObstaclePenetrationEscapeDirectionY = Multiply_Float_Vector2(penetrationDistance, &circleObstaclePenetrationEscapeDirectionY)
				}
			}
		}
	}

	circlePenetrationObjectEscapeDirection.x += circleObstaclePenetrationEscapeDirectionX.x + circleObstaclePenetrationEscapeDirectionY.x
	circlePenetrationObjectEscapeDirection.y += circleObstaclePenetrationEscapeDirectionX.y + circleObstaclePenetrationEscapeDirectionY.y

	finalInputDirection := inputDirection
	if circleCollision {
		normal := circlePenetrationObjectEscapeDirection

		entityBBDims := entityCollisionShapeRef.GetBoundingBoxDims()
		entityColliderCentrePos := Vector2{entityPos.x + entityBBDims.x*0.5, entityPos.y + entityBBDims.y*0.5}
		relPos := Subtract_Vector2(&entityColliderCentrePos, &circleColliderOriginPosition)
		relPos = Vector2{Sign(relPos.x), Sign(relPos.y)}

		absNormal := Vector2{math.Abs(normal.x), math.Abs(normal.y)}

		tangent := normal
		if inputDirection.x != 0 {
			tangent = Vector2{Sign(inputDirection.x) * absNormal.y, relPos.y * absNormal.x}
		}
		if inputDirection.y != 0 {
			tangent = Vector2{relPos.x * absNormal.y, Sign(inputDirection.y) * absNormal.x}
		}

		tangent = Normalise_Vector2(&tangent)
		magnitudeInput := Magnitude_Vector2(&inputDirection)
		finalInputDirection = Multiply_Float_Vector2(magnitudeInput, &tangent)
	}

	if !collidingOnX && !collidingOnY {
		entityPosRef.x += totalMoveAmount.x
		entityPosRef.y += totalMoveAmount.y
	} else if !circleCollision && !collidingOnX && collidingOnY {
		entityPosRef.x += finalInputDirection.x * entityMoveAmountPerFrame
		entityPosRef.y += finalInputDirection.y * maxMoveAmtOnY
	} else if !circleCollision && collidingOnX && !collidingOnY {
		entityPosRef.y += finalInputDirection.y * entityMoveAmountPerFrame
		entityPosRef.x += finalInputDirection.x * maxMoveAmtOnX
	} else if circleCollision {
		finalMoveAmt := Multiply_Float_Vector2(entityMoveAmountPerFrame, &finalInputDirection)
		entityPosRef.x += finalMoveAmt.x
		entityPosRef.y += finalMoveAmt.y
	}

}
