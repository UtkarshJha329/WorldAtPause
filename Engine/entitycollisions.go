package Engine

import "math"

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

func (curScene *Scene) DoesEntityCollideWithObstacles(roomIndex Vector2Int, entityID int, const_entityPositionRef *Vector2, const_entityCollisionShapeRef CollisionShape, CollideWithPlayer bool) EntityObstacleCollisionData {

	curRoom, curRoomHasSomeData := curScene.RoomsData[roomIndex]
	entityComponents := curScene.EntityComponentsForScene

	curCollisionData := EntityObstacleCollisionData{
		CollidedWithObstacle: false,
	}

	if !curRoomHasSomeData {
		return curCollisionData
	}

	for _, obstacleEntityID := range curRoom.ObstacleEntityIDs {

		if obstacleEntityID == entityID || entityComponents.IsEntityDead(obstacleEntityID) {
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

	if CollideWithPlayer && entityID != curScene.EntityComponentsForScene.PlayerEntityID && entityComponents.IsEntityAlive(curScene.EntityComponentsForScene.PlayerEntityID) {

		playerPos := curScene.EntityComponentsForScene.Positions[curScene.EntityComponentsForScene.PlayerEntityID]
		playerCollisionShapeRef := curScene.EntityComponentsForScene.CollisionShapes[curScene.EntityComponentsForScene.PlayerEntityID]

		curCollisionData.CollisionNormal, curCollisionData.CollisionPenetrationNormal, curCollisionData.CollisionPenetrationAmount, curCollisionData.CollidedWithObstacle = CollisionShapeOverlapsWithCollisionShape(const_entityPositionRef, const_entityCollisionShapeRef, &playerPos, playerCollisionShapeRef)
		if curCollisionData.CollidedWithObstacle {

			curCollisionData.CollidedWithObstacleEntityID = curScene.EntityComponentsForScene.PlayerEntityID
			curCollisionData.PositionOfCollidedWithObstacle = playerPos
			curCollisionData.CollidedWithObstacleCollisionShapeRef = playerCollisionShapeRef
			curCollisionData.ColliderOffsetOriginPosition = *playerCollisionShapeRef.GetOffsetOrigin(&playerPos)

			return curCollisionData
		}

	}

	return curCollisionData
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

	curTilemap := curScene.EntityComponentsForScene.Tilemap

	colliderCollisionPoints := circleBoundingBoxCol.GetCollisionPoints()
	for _, collisionPoint := range *colliderCollisionPoints {
		collisionPointWorldPos := Add_Vector2(&collisionPoint, &colliderEntityWorldPos)

		if curTilemap.PointCollidesWithTilemapCollisionLayerInRoom(&roomIndex, &collisionPointWorldPos) {
			collisionResult.Collided = true

			curTile := curScene.EntityComponentsForScene.Tilemap.GetTileOfPointInRoom(&roomIndex, &collisionPointWorldPos)
			curTilePos := curScene.EntityComponentsForScene.Tilemap.GetTilePosInWorld(&roomIndex, &curTile)

			_, _, _, collided := CollisionShapeOverlapsWithCollisionShape(&colliderEntityWorldPos, &circleBoundingBoxCol, &curTilePos, &curTilemap.TileCollisionShape)

			if collided {
				tileBoundingBox := curScene.EntityComponentsForScene.Tilemap.TileCollisionShape

				collisionResult.CollisionTileNormal, collisionResult.CollisionSeparationNormal, collisionResult.PenetrationAmount, collisionResult.Collided = CollisionShapeOverlapsWithCollisionShape(&colliderEntityWorldPos, circleCollider, &curTilePos, &tileBoundingBox)

				if collisionResult.Collided {
					circleColliderOrigin := circleCollider.GetOffsetOrigin(&colliderEntityWorldPos)
					normalFromTileCollider := tileBoundingBox.GetNormalFromPoint(&curTilePos, circleColliderOrigin)
					collisionResult.CollisionTileNormal = *normalFromTileCollider

					break
				}
			}
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
	CollideWithPlayer    bool
	MovementLock         Vector2
}

func MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirection Vector2, circleColliderOriginPosition Vector2, entityPos Vector2, inputDirection Vector2, entityCollisionShapeRef CollisionShape, entityMoveAmountPerFrame float64) Vector2 {

	normal := circleObstaclePenetrationEscapeDirection
	entityBBDims := entityCollisionShapeRef.GetBoundingBoxDims()
	entityColliderCentrePos := Vector2{entityPos.X + entityBBDims.X*0.5, entityPos.Y + entityBBDims.Y*0.5}
	relPos := Subtract_Vector2(&entityColliderCentrePos, &circleColliderOriginPosition)
	relPos = Vector2{math.Copysign(1.0, relPos.X), math.Copysign(1.0, relPos.Y)}

	absNormal := Vector2{math.Abs(normal.X), math.Abs(normal.Y)}

	tangentA := Vector2{math.Copysign(1.0, inputDirection.X) * absNormal.Y, relPos.Y * absNormal.X}
	tangentA = Normalise_Vector2(&tangentA)
	tangentB := Vector2{relPos.X * absNormal.Y, math.Copysign(1.0, inputDirection.Y) * absNormal.X}
	tangentB = Normalise_Vector2(&tangentB)

	normalisedInputDir := Normalise_Vector2(&inputDirection)
	inputDirDotTangentA := Dot_Vector2(&tangentA, &normalisedInputDir)
	inputDirDotTangentB := Dot_Vector2(&tangentB, &normalisedInputDir)

	tangent := tangentA
	if inputDirDotTangentB > inputDirDotTangentA {
		tangent = tangentB
	}

	tangent = Normalise_Vector2(&tangent)
	return Multiply_Float_Vector2(entityMoveAmountPerFrame, &tangent)
}

func FinalMoveAmountIfObstacleObjectHasCircleCollider(axis int, entityID int, entityPos Vector2, entityCollisionShapeRef CollisionShape, inputDirection Vector2, totalMoveAmount Vector2, entityMoveAmountPerFrame float64, obstacleCollisionResult *EntityObstacleCollisionData, collideAndMoveParameters *CollideAndMoveCollisionParameters) (finalMoveAmount Vector2) {

	boxCircleCollisionEpsilon := -0.1
	boxCircleNormalMultiplier := -1.0

	circleCircleCollisionEpsilon := 0.1
	circleCircleNormalMultiplier := 1.0

	normalPlusEpsilon := boxCircleNormalMultiplier + boxCircleCollisionEpsilon

	if entityCollisionShapeRef.CollisionShapeType() == Circle {
		normalPlusEpsilon = circleCircleNormalMultiplier + circleCircleCollisionEpsilon
	}

	penetrationDistance := obstacleCollisionResult.CollisionPenetrationAmount
	circleColliderOriginPosition := obstacleCollisionResult.ColliderOffsetOriginPosition

	circleObstaclePenetrationEscapeDirection := Multiply_Float_Vector2(penetrationDistance, &obstacleCollisionResult.CollisionPenetrationNormal)

	if collideAndMoveParameters.SlideWhenCollide {
		moveAmountOnTangentForThisCircleObstacle := MoveAmountOnTangentForCircleCollider(circleObstaclePenetrationEscapeDirection, circleColliderOriginPosition, entityPos, inputDirection, entityCollisionShapeRef, entityMoveAmountPerFrame)

		finalMoveAmount = moveAmountOnTangentForThisCircleObstacle
	} else {
		totalMoveAmount := Vector2{totalMoveAmount.X, 0.0}
		if axis == 1 {
			totalMoveAmount = Vector2{0.0, totalMoveAmount.Y}
		}

		circleObstaclePenetrationEscapeDirection = Multiply_Float_Vector2(normalPlusEpsilon, &circleObstaclePenetrationEscapeDirection)
		finalMoveAmount = Add_Vector2(&circleObstaclePenetrationEscapeDirection, &totalMoveAmount)
	}

	return finalMoveAmount
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
		obstacleCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, entityID, &Vector2{playerToMoveXPos, entityPos.Y}, entityCollisionShapeRef, collideAndMoveParameters.CollideWithPlayer)
		collideAndMoveResultData.ObstacleXMoveCollisionResult = obstacleCollisionResult
		collidedOnXWithObstacle = obstacleCollisionResult.CollidedWithObstacle
		if collidedOnXWithObstacle {

			finalMoveAmount := Vector2{0.0, 0.0}
			maxMoveAmountX := 0.0

			if entityCollisionShapeRef.CollisionShapeType() == Box {

				if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
					if inputDirection.X > 0 {
						entityCollisionPoints := entityCollisionShapeRef.GetCollisionPoints()
						topRight := (*entityCollisionPoints)[1]

						maxMoveAmountX = obstacleCollisionResult.PositionOfCollidedWithObstacle.X - (entityPos.X + topRight.X)
					} else if inputDirection.X < 0 {
						entityCollisionPoints := entityCollisionShapeRef.GetCollisionPoints()
						topLeft := (*entityCollisionPoints)[0]

						maxMoveAmountX = (entityPos.X + topLeft.X) - (obstacleCollisionResult.PositionOfCollidedWithObstacle.X + obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().X)
					}
					finalMoveAmount = Vector2{inputDirection.X * maxMoveAmountX, inputDirection.Y * entityMoveAmountPerFrame}
				} else if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

					finalMoveAmount = FinalMoveAmountIfObstacleObjectHasCircleCollider(0, entityID, entityPos, entityCollisionShapeRef, inputDirection, totalMoveAmount, entityMoveAmountPerFrame, &obstacleCollisionResult, &collideAndMoveParameters)
				}
			} else if entityCollisionShapeRef.CollisionShapeType() == Circle {

				if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {

					finalMoveAmount = FinalMoveAmountIfObstacleObjectHasCircleCollider(0, entityID, entityPos, entityCollisionShapeRef, inputDirection, totalMoveAmount, entityMoveAmountPerFrame, &obstacleCollisionResult, &collideAndMoveParameters)

				} else if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

					finalMoveAmount = FinalMoveAmountIfObstacleObjectHasCircleCollider(0, entityID, entityPos, entityCollisionShapeRef, inputDirection, totalMoveAmount, entityMoveAmountPerFrame, &obstacleCollisionResult, &collideAndMoveParameters)
				}
			}

			entityPosRef.X += finalMoveAmount.X * collideAndMoveParameters.MovementLock.X
			entityPosRef.Y += finalMoveAmount.Y * collideAndMoveParameters.MovementLock.Y
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

	if !collidingOnYWithTilemap.Collided && collideAndMoveParameters.CollideWithObstacles {
		obstacleCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, entityID, &Vector2{entityPos.X, playerToMoveYPos}, entityCollisionShapeRef, collideAndMoveParameters.CollideWithPlayer)
		collideAndMoveResultData.ObstacleYMoveCollisionResult = obstacleCollisionResult
		collidedOnYWithObstacle = obstacleCollisionResult.CollidedWithObstacle
		if collidedOnYWithObstacle {

			finalMoveAmount := Vector2{0.0, 0.0}
			maxMoveAmountY := 0.0

			if entityCollisionShapeRef.CollisionShapeType() == Box {

				if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
					if inputDirection.Y > 0 {
						entityCollisionPoints := entityCollisionShapeRef.GetCollisionPoints()
						bottomRight := (*entityCollisionPoints)[2]
						maxMoveAmountY = obstacleCollisionResult.PositionOfCollidedWithObstacle.Y - (entityPos.Y + bottomRight.Y)
					} else if inputDirection.Y < 0 {
						entityCollisionPoints := entityCollisionShapeRef.GetCollisionPoints()
						topRight := (*entityCollisionPoints)[1]
						maxMoveAmountY = (entityPos.Y + topRight.Y) - (obstacleCollisionResult.PositionOfCollidedWithObstacle.Y + obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().Y)
					}
					finalMoveAmount = Vector2{inputDirection.X * entityMoveAmountPerFrame, inputDirection.Y * maxMoveAmountY}

				} else if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

					finalMoveAmount = FinalMoveAmountIfObstacleObjectHasCircleCollider(1, entityID, entityPos, entityCollisionShapeRef, inputDirection, totalMoveAmount, entityMoveAmountPerFrame, &obstacleCollisionResult, &collideAndMoveParameters)

				}
			} else if entityCollisionShapeRef.CollisionShapeType() == Circle {

				if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {

					finalMoveAmount = FinalMoveAmountIfObstacleObjectHasCircleCollider(1, entityID, entityPos, entityCollisionShapeRef, inputDirection, totalMoveAmount, entityMoveAmountPerFrame, &obstacleCollisionResult, &collideAndMoveParameters)

				} else if obstacleCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

					finalMoveAmount = FinalMoveAmountIfObstacleObjectHasCircleCollider(1, entityID, entityPos, entityCollisionShapeRef, inputDirection, totalMoveAmount, entityMoveAmountPerFrame, &obstacleCollisionResult, &collideAndMoveParameters)
				}

			}
			entityPosRef.X += finalMoveAmount.X * collideAndMoveParameters.MovementLock.X
			entityPosRef.Y += finalMoveAmount.Y * collideAndMoveParameters.MovementLock.Y
		}
	}

	if !collidedOnXWithObstacle && !collidedOnYWithObstacle && !collidingOnXWithTilemap.Collided {
		entityPosRef.X += totalMoveAmount.X * collideAndMoveParameters.MovementLock.X
	}
	if !collidedOnXWithObstacle && !collidedOnYWithObstacle && !collidingOnYWithTilemap.Collided {
		entityPosRef.Y += totalMoveAmount.Y * collideAndMoveParameters.MovementLock.Y
	}

	return collideAndMoveResultData
}
