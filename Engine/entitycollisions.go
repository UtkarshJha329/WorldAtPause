package Engine

import (
	"fmt"
	"math"
)

type EntityObstacleCollisionData struct {
	CollidedWithObstacle                  bool
	CollidedWithObstacleEntityID          int
	PositionOfCollidedWithObstacle        Vector2
	ColliderOffsetOriginPosition          Vector2
	CollisionNormalFromObstacle           Vector2
	PenetrationAmount                     float64
	CollidedWithObstacleCollisionShapeRef CollisionShape
}

// boxEntityCircleObstacleCirclePenetrationData CircleBoxOverlapCirclePenetrationData

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

		// if obstacleEntityID == entityID {
		// 	continue
		// }

		obstaclePos := entityComponents.Positions[obstacleEntityID]
		obstacleCollisionShapeRef := entityComponents.CollisionShapes[obstacleEntityID]
		normal, penetrationAmount, collided := CollisionShapeOverlapsWithCollisionShape(const_entityPositionRef, const_entityCollisionShapeRef, &obstaclePos, obstacleCollisionShapeRef)
		if collided {

			curCollisionData.CollidedWithObstacle = true
			curCollisionData.CollidedWithObstacleEntityID = obstacleEntityID
			curCollisionData.PositionOfCollidedWithObstacle = entityComponents.Positions[obstacleEntityID]
			curCollisionData.CollidedWithObstacleCollisionShapeRef = obstacleCollisionShapeRef

			curCollisionData.ColliderOffsetOriginPosition = *obstacleCollisionShapeRef.GetOffsetOrigin(&obstaclePos)

			curCollisionData.CollisionNormalFromObstacle = normal
			curCollisionData.PenetrationAmount = penetrationAmount

			return curCollisionData
		}
	}

	// if !entityComponents.IsEntityPlayer(entityID) {

	// 	playerPos := entityComponents.Positions[entityComponents.PlayerEntityID]
	// 	playerCollisionShapeRef := entityComponents.CollisionShapes[entityComponents.PlayerEntityID]
	// 	normal, penetrationAmount, collided := CollisionShapeOverlapsWithCollisionShape(const_entityPositionRef, const_entityCollisionShapeRef, &playerPos, playerCollisionShapeRef)
	// 	if collided {

	// 		curCollisionData.CollidedWithObstacle = true
	// 		curCollisionData.CollidedWithObstacleEntityID = entityComponents.PlayerEntityID
	// 		curCollisionData.PositionOfCollidedWithObstacle = playerPos
	// 		curCollisionData.CollidedWithObstacleCollisionShapeRef = playerCollisionShapeRef

	// 		curCollisionData.ColliderOffsetOriginPosition = *playerCollisionShapeRef.GetOffsetOrigin(&playerPos)

	// 		curCollisionData.CollisionNormalFromObstacle = normal
	// 		curCollisionData.PenetrationAmount = penetrationAmount

	// 		return curCollisionData
	// 	}

	// }

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

			curTile := curScene.EntityComponentsForScene.Tilemap.GetTileOfPointInLevel(&roomIndex, &collisionPointWorldPos)
			curTilePos := curScene.EntityComponentsForScene.Tilemap.GetTilePosInWorld(&roomIndex, &curTile)
			collisionResult.CollidedTilePos = curTilePos

			collisionResult.CollisionSeparationNormal, collisionResult.PenetrationAmount, collisionResult.Collided = CollisionShapeOverlapsWithCollisionShape(&colliderEntityWorldPos, boxCollider, &curTilePos, &curTilemap.TileCollisionShape)
			// fmt.Println("Box tile col Pen amount :", collisionResult.PenetrationAmount)
			// collisionResult.CollisionNormal = Multiply_Float_Vector2(-1.0, &collisionResult.CollisionNormal)
			if collisionResult.Collided {
				boxColliderPos := boxCollider.GetOffsetOrigin(&colliderEntityWorldPos)
				boxColliderBounds := boxCollider.GetBoundingBoxDims()
				boxColliderHalfBounds := Multiply_Float_Vector2(0.5, boxColliderBounds)
				boxColliderCentre := Add_Vector2(boxColliderPos, &boxColliderHalfBounds)
				normalFromTileCollider := curTilemap.TileCollisionShape.GetNormalFromPoint(&curTilePos, &boxColliderCentre)
				collisionResult.CollisionTileNormal = *normalFromTileCollider

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

		collisionResult.CollisionSeparationNormal, collisionResult.PenetrationAmount, collisionResult.Collided = CollisionShapeOverlapsWithCollisionShape(&colliderEntityWorldPos, circleCollider, &collisionTilePosition, &tileBoundingBox)
		// fmt.Println("Circle tile col Pen amount :", collisionResult.PenetrationAmount)
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

func (curScene *Scene) CollideAndMoveEntityWithTilemapAndObstacles(curRoomIndex Vector2Int, entityID int, inputDirection Vector2, totalMoveAmount Vector2, entityMoveAmountPerFrame float64) CollideAndMoveCollisionResultSteps {

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
	// collidingOnXWithTilemap := curScene.ShapeCollidesWithTilemapCollisionLayer(curRoomIndex, entityID, Vector2{inputDirection.X, 0.0}, Vector2{totalMoveAmount.X, 0.0}, entityMoveAmountPerFrame)
	obstacleXCollisionResult := EntityObstacleCollisionData{}

	if !collidingOnXWithTilemap {
		// if !collidingOnXWithTilemap.Collided {
		obstacleXCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, entityID, &Vector2{playerToMoveXPos, entityPos.Y}, entityCollisionShapeRef)

		if obstacleXCollisionResult.CollidedWithObstacle {

			fmt.Println("Colliding on X with obs.")

			finalMoveAmount := Vector2{0.0, 0.0}
			maxMoveAmountX := 0.0

			if obstacleXCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
				if inputDirection.X > 0 {
					maxMoveAmountX = obstacleXCollisionResult.PositionOfCollidedWithObstacle.X - (entityPos.X + entityCollisionShapeRef.GetBoundingBoxDims().X)
				} else if inputDirection.X < 0 {
					maxMoveAmountX = entityPos.X - (obstacleXCollisionResult.PositionOfCollidedWithObstacle.X + obstacleXCollisionResult.CollidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().X)
				}
				finalMoveAmount = Vector2{inputDirection.X * maxMoveAmountX, inputDirection.Y * entityMoveAmountPerFrame}
			} else if obstacleXCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

				penetrationDistance := obstacleXCollisionResult.PenetrationAmount
				if penetrationDistance != 0 {

					circleColliderOriginPosition := obstacleXCollisionResult.ColliderOffsetOriginPosition

					circleObstaclePenetrationEscapeDirectionX := Multiply_Float_Vector2(1.0, &obstacleXCollisionResult.CollisionNormalFromObstacle)
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
	// collidingOnYWithTilemap := curScene.ShapeCollidesWithTilemapCollisionLayer(curRoomIndex, entityID, Vector2{0.0, inputDirection.Y}, Vector2{0.0, totalMoveAmount.Y}, entityMoveAmountPerFrame)
	obstacleYCollisionResult := EntityObstacleCollisionData{}

	if !collidingOnYWithTilemap {
		// if !collidingOnYWithTilemap.Collided {
		obstacleYCollisionResult := curScene.DoesEntityCollideWithObstacles(curRoomIndex, entityID, &Vector2{entityPos.X, playerToMoveYPos}, entityCollisionShapeRef)

		if obstacleYCollisionResult.CollidedWithObstacle {

			fmt.Println("Colliding on Y with obs.")

			finalMoveAmount := Vector2{0.0, 0.0}
			maxMoveAmountY := 0.0

			if obstacleYCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Box {
				if inputDirection.Y > 0 {
					maxMoveAmountY = obstacleYCollisionResult.PositionOfCollidedWithObstacle.Y - (entityPos.Y + entityCollisionShapeRef.GetBoundingBoxDims().Y)
				} else if inputDirection.Y < 0 {
					maxMoveAmountY = entityPos.Y - (obstacleYCollisionResult.PositionOfCollidedWithObstacle.Y + obstacleYCollisionResult.CollidedWithObstacleCollisionShapeRef.GetBoundingBoxDims().Y)
				}
				finalMoveAmount = Vector2{inputDirection.X * entityMoveAmountPerFrame, inputDirection.Y * maxMoveAmountY}

			} else if obstacleYCollisionResult.CollidedWithObstacleCollisionShapeRef.CollisionShapeType() == Circle {

				penetrationDistance := obstacleYCollisionResult.PenetrationAmount
				if penetrationDistance != 0 {
					circleColliderOriginPosition := obstacleYCollisionResult.ColliderOffsetOriginPosition

					circleObstaclePenetrationEscapeDirectionY := Multiply_Float_Vector2(1.0, &obstacleYCollisionResult.CollisionNormalFromObstacle)
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

	// if !obstacleXCollisionResult.CollidedWithObstacle && !obstacleYCollisionResult.CollidedWithObstacle && !collidingOnXWithTilemap.Collided {
	if !obstacleXCollisionResult.CollidedWithObstacle && !obstacleYCollisionResult.CollidedWithObstacle && !collidingOnXWithTilemap {
		entityPosRef.X += totalMoveAmount.X
	}
	// if !obstacleXCollisionResult.CollidedWithObstacle && !obstacleYCollisionResult.CollidedWithObstacle && !collidingOnYWithTilemap.Collided {
	if !obstacleXCollisionResult.CollidedWithObstacle && !obstacleYCollisionResult.CollidedWithObstacle && !collidingOnYWithTilemap {
		entityPosRef.Y += totalMoveAmount.Y
	}

	// return CollideAndMoveCollisionResultSteps{
	// 	TilemapXMoveCollisionResult:  collidingOnXWithTilemap,
	// 	TilemapYMoveCollisionResult:  collidingOnYWithTilemap,
	// 	ObstacleXMoveCollisionResult: obstacleXCollisionResult,
	// 	ObstacleYMoveCollisionResult: obstacleYCollisionResult,
	// }

	return CollideAndMoveCollisionResultSteps{
		TilemapXMoveCollisionResult:  ShapeCollisionWithTilemapResult{Collided: collidingOnXWithTilemap},
		TilemapYMoveCollisionResult:  ShapeCollisionWithTilemapResult{Collided: collidingOnYWithTilemap},
		ObstacleXMoveCollisionResult: obstacleXCollisionResult,
		ObstacleYMoveCollisionResult: obstacleYCollisionResult,
	}
}
