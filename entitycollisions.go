package main

type CollisionData struct {
	collidedWithObstacle                  bool
	collidedWithObstacleEntityID          int
	collidedWithObstaclePosition          Vector2
	collidedWithObstacleCollisionShapeRef CollisionShape
}

func (entityComponents *EntityComponents) DoesEntityCollideWithObstacles(const_entityPositionRef *Vector2, const_entityCollisionShapeRef CollisionShape) CollisionData {

	curCollisionData := CollisionData{
		collidedWithObstacle: false,
	}
	for _, obstacleEntityID := range entityComponents.obstacleEntityIDs {

		obstaclePosRef := &entityComponents.positions[obstacleEntityID]
		obstacleCollisionShapeRef := entityComponents.collisionShapes[obstacleEntityID]
		if CollisionShapeOverlapsWithCollisionShape(const_entityPositionRef, const_entityCollisionShapeRef, obstaclePosRef, obstacleCollisionShapeRef) {

			curCollisionData.collidedWithObstacle = true
			curCollisionData.collidedWithObstacleEntityID = obstacleEntityID
			curCollisionData.collidedWithObstaclePosition = *obstaclePosRef
			curCollisionData.collidedWithObstacleCollisionShapeRef = obstacleCollisionShapeRef

			return curCollisionData
		}
	}

	return curCollisionData
}

func (entityComponentsRef *EntityComponents) MoveAndCollideEntityWithTilemapAndObstacles(entityID int, inputDirection Vector2, totalMoveAmount Vector2, entityMoveAmountPerFrame float64) {

	entityPos := entityComponentsRef.positions[entityID]
	entityPosRef := &entityComponentsRef.positions[entityID]

	entityCollisionShapeRef := entityComponentsRef.collisionShapes[entityID]

	playerToMoveXPos := entityPos.x + totalMoveAmount.x
	moveDirectionCheckOffsetX := 0.0
	if inputDirection.x > 0 {
		moveDirectionCheckOffsetX = entityCollisionShapeRef.(*BoxCollider).size.x
	}

	collidingOnXWithTilemap := entityComponentsRef.tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{playerToMoveXPos + moveDirectionCheckOffsetX, entityPos.y})
	collidingOnX := collidingOnXWithTilemap

	maxMoveAmtOnX := 0.0
	if !collidingOnXWithTilemap {
		obstacleCollisionResult := entityComponentsRef.DoesEntityCollideWithObstacles(&Vector2{playerToMoveXPos, entityPos.y}, entityCollisionShapeRef)
		if obstacleCollisionResult.collidedWithObstacle {
			collidingOnX = true
			if inputDirection.x > 0 {
				maxMoveAmtOnX = obstacleCollisionResult.collidedWithObstaclePosition.x - (entityPos.x + entityCollisionShapeRef.(*BoxCollider).size.x)
			} else if inputDirection.x < 0 {
				maxMoveAmtOnX = entityPos.x - (obstacleCollisionResult.collidedWithObstaclePosition.x + obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.(*BoxCollider).size.x)
			}
		}
	}

	playerToMoveYPos := entityPos.y + totalMoveAmount.y
	moveDirectionCheckOffsetY := 0.0
	if inputDirection.y > 0 {
		moveDirectionCheckOffsetY = entityCollisionShapeRef.(*BoxCollider).size.y
	}

	collidingOnYWithTilemap := entityComponentsRef.tilemap.PointCollidesWithTilemapCollisionLayer(&Vector2{entityPos.x, playerToMoveYPos + moveDirectionCheckOffsetY})
	collidingOnY := collidingOnYWithTilemap

	maxMoveAmtOnY := 0.0
	if !collidingOnYWithTilemap {
		obstacleCollisionResult := entityComponentsRef.DoesEntityCollideWithObstacles(&Vector2{entityPos.x, playerToMoveYPos}, entityCollisionShapeRef)
		if obstacleCollisionResult.collidedWithObstacle {
			collidingOnY = true
			if inputDirection.y > 0 {
				maxMoveAmtOnY = obstacleCollisionResult.collidedWithObstaclePosition.y - (entityPos.y + entityCollisionShapeRef.(*BoxCollider).size.y)
			} else if inputDirection.y < 0 {
				maxMoveAmtOnY = entityPos.y - (obstacleCollisionResult.collidedWithObstaclePosition.y + obstacleCollisionResult.collidedWithObstacleCollisionShapeRef.(*BoxCollider).size.y)
			}
		}
	}

	if !collidingOnX && !collidingOnY {
		entityPosRef.x += totalMoveAmount.x
		entityPosRef.y += totalMoveAmount.y
	} else if !collidingOnX && collidingOnY {
		entityPosRef.x += inputDirection.x * entityMoveAmountPerFrame
		entityPosRef.y += inputDirection.y * maxMoveAmtOnY
	} else if collidingOnX && !collidingOnY {
		entityPosRef.y += inputDirection.y * entityMoveAmountPerFrame
		entityPosRef.x += inputDirection.x * maxMoveAmtOnX
	}
}
