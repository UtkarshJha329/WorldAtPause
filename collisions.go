package main

import (
	"math"
)

type BoxCollider struct {
	size Vector2
}

func (bc *BoxCollider) CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool {
	bcTopLeft := Vector2{const_origin.x, const_origin.y}
	bcBottomRight := Vector2{const_origin.x + bc.size.x, const_origin.y + bc.size.y}

	if const_point.x >= bcTopLeft.x && const_point.x <= bcBottomRight.x {
		if const_point.y >= bcTopLeft.y && const_point.y <= bcBottomRight.y {
			return true
		}
	}

	return false
}

func (bc *BoxCollider) CollisionPoints(const_origin *Vector2) *[]Vector2 {

	topLeft := *const_origin
	bottomRight := Add_Vector2(const_origin, &bc.size)
	topRight := Vector2{bottomRight.x, topLeft.y}
	bottomLeft := Vector2{topLeft.x, bottomRight.y}

	return &[]Vector2{
		topLeft,
		topRight,
		bottomRight,
		bottomLeft,
	}
}

type CircleCollider struct {
	radius float64
}

func (cc *CircleCollider) CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool {
	return DistanceSquare_Vector2(const_point, const_origin) <= math.Pow(cc.radius, 2.0)
}

func (cc *CircleCollider) CollisionPoints(const_origin *Vector2) *[]Vector2 {
	return &[]Vector2{*const_origin}
}

type CollisionShape interface {
	CollisionPoints(const_origin *Vector2) *[]Vector2
	CollisionShapeIncludesPoint(const_origin *Vector2, const_point *Vector2) bool
}

func CollisionShapeOverlapsWithCollisionShape(const_originA *Vector2, csA CollisionShape, const_originB *Vector2, csB CollisionShape) bool {

	csAPoints := csA.CollisionPoints(const_originA)
	csBPoints := csB.CollisionPoints(const_originB)

	numPointsA := len(*csAPoints)
	numPointsB := len(*csBPoints)

	if numPointsA == 1 && numPointsB == 1 {
		return DistanceSquare_Vector2(const_originA, const_originB) <= math.Pow(csA.(*CircleCollider).radius+csB.(*CircleCollider).radius, 2.0)
	}
	if numPointsA == 1 {
		for _, collisionPointInB := range *csBPoints {
			if csA.CollisionShapeIncludesPoint(const_originA, &collisionPointInB) {
				return true
			}
		}
		return false
	} else if numPointsB == 1 {
		for _, collisionPointInA := range *csAPoints {
			if csB.CollisionShapeIncludesPoint(const_originB, &collisionPointInA) {
				return true
			}
		}
		return false
	} else {
		for _, collisionPointInB := range *csBPoints {
			if csA.CollisionShapeIncludesPoint(const_originA, &collisionPointInB) {
				return true
			}
		}
		return false
	}
}

type CollisionShapeAssetData struct {
	CollisionShapeType string   `json:"CollisionShapeType"`
	Size               *Vector2 `json:"Size"`
	Radius             *float64 `json:"Radius"`
}
